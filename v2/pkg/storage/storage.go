package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/google/uuid"
)

const (
	containerName = "azureproject"
	accountName   = "projectstorage69"
)

var containerURL azblob.ContainerURL

// InitializeStorage creates azure storage instances
func InitializeStorage() {
	credential, err := azblob.NewSharedKeyCredential("projectstorage69", "+kpPRjIysUxKy2QhxRsJRrFGmOoJY/3o6eD4ZTOqTqC7wABPS0CUyvn3dzbOy4MYtBQ8sS+VO1SbxPmqupwgNA==")
	if err != nil {
		log.Println(err)
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))
	containerURL = azblob.NewContainerURL(*URL, p)

	ctx := context.Background()
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		// Ignore if already created
	}

	// fmt.Printf("Creating a dummy file to test the upload and download\n")
	// data := []byte("hello world this is a blob\n")
	// fileName := "randomfile.txt"
	// err = ioutil.WriteFile(fileName, data, 0700)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// blobURL := containerURL.NewBlockBlobURL(fileName)
	// file, err := os.Open(fileName)

	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// fmt.Printf("Uploading the file with blob name: %s\n", fileName)
	// _, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
	// 	BlockSize:   4 * 1024 * 1024,
	// 	Parallelism: 16})

	// if err != nil {
	// 	log.Println(err)
	// }
}

func writeToLocalStorage(readerCloser *io.ReadCloser) (string, error) {
	fileName := filepath.Join("tmp", uuid.New().String())
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		os.Mkdir("tmp", os.ModePerm)
	}

	outFile, err := os.Create(fileName)
	defer outFile.Close()

	if err != nil {
		log.Println(err)
		return "", err
	}

	_, err = io.Copy(outFile, *readerCloser)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return fileName, nil
}

// UploadToStorage will upload blob from reader to azure storage
func UploadToStorage(readCloser *io.ReadCloser, destination string) error {
	fileName, err := writeToLocalStorage(readCloser)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		// Ignore if already created
	}

	blobURL := containerURL.NewBlockBlobURL(destination)
	file, err := os.Open(fileName)
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})

	return err
}
