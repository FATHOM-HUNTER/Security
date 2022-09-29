package storage

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Ovenoboyo/basic_webserver/pkg/database"
	db "github.com/Ovenoboyo/basic_webserver/pkg/database"
	"github.com/google/uuid"
)

var containerURL azblob.ContainerURL

// InitializeStorage creates azure storage instances
func InitializeStorage() {
	var (
		containerName = os.Getenv("STORAGE_CONTAINER")
		accountName   = os.Getenv("STORAGE_ACCOUNT")
		accountKey    = os.Getenv("STORAGE_KEY")
	)

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Println(err)
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))
	containerURL = azblob.NewContainerURL(*URL, p)
}

func encryptLocalFile(fileName string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plaintext, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("cipher GCM err: %v", err.Error())
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("nonce  err: %v", err.Error())
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, plaintext, nil)

	encrypted := filepath.Join(filepath.Dir(fileName), uuid.New().String())
	err = ioutil.WriteFile(encrypted, cipherText, 0777)
	if err != nil {
		log.Fatalf("write file err: %v", err.Error())
		return "", err
	}

	return encrypted, nil
}

func writeToLocalStorage(readerCloser *io.ReadCloser) (string, error) {
	fileName := filepath.Join("tmp", uuid.New().String())
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		os.Mkdir("tmp", os.ModePerm)
	}

	outFile, err := os.Create(fileName)
	defer outFile.Close()

	if err != nil {
		return "", err
	}

	_, err = io.Copy(outFile, *readerCloser)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func getMD5(file *os.File) string {
	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	sum := h.Sum(nil)
	return hex.EncodeToString(sum[:])
}

// UploadToStorage will upload blob from reader to azure storage
func UploadToStorage(readCloser *io.ReadCloser, destination string, uid string, key string) error {
	plainFile, err := writeToLocalStorage(readCloser)
	defer os.Remove(plainFile)

	if err != nil {
		return err
	}

	fileName, err := encryptLocalFile(plainFile, key)
	defer os.Remove(fileName)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		// Ignore if already created
	}

	blobURL := containerURL.NewBlockBlobURL(destination + "-" + uid)
	file, err := os.Open(fileName)
	defer file.Close()

	if err != nil {
		return err
	}

	md5 := getMD5(file)
	stat, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	exists, version, err := db.GetExistingFile(destination, md5, uid)
	if err != nil {
		return err
	}

	if exists && len(version) != 0 {
		return errors.New("Exact file already exists")
	}

	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})

	if err != nil {
		return err
	}

	props, err := blobURL.GetProperties(ctx, azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return err
	}

	err = db.AddFileMetaToDB(destination, md5, uid, int(stat.Size()), props.VersionID())

	return err
}

func decryptFile(closer *io.ReadCloser, k string) (io.ReadCloser, error) {
	cipherText, err := io.ReadAll(*closer)
	if err != nil {
		return nil, err
	}

	log.Println("decrypt key: " + k)
	key := []byte(k)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("cipher err: %v", err.Error())
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("cipher GCM err: %v", err.Error())
		return nil, err
	}

	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		log.Fatalf("decrypt file err: %v", err.Error())
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(plainText)), nil
}

func DownloadBlob(fileName string, uid string, version string, key string) (io.ReadCloser, error) {
	ctx := context.Background()
	blobURL := containerURL.NewBlockBlobURL(fileName + "-" + uid).WithVersionID(version)
	resp, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, err
	}

	bodyStream := resp.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	closer, err := decryptFile(&bodyStream, key)
	if err != nil {
		return nil, err
	}
	return closer, nil
}

func DeleteBlob(fileName string, uid string, version string) error {
	ctx := context.Background()
	_, err := containerURL.NewBlockBlobURL(fileName+"-"+uid).WithVersionID(version).Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		_, err := containerURL.NewBlockBlobURL(fileName+"-"+uid).Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
		if err != nil {
			return err
		}
	}
	return database.RemoveBlob(uid, fileName, version)
}
