package database

import (
	"log"
	"strconv"
	"strings"
	"time"
)

type FileMetadata struct {
	FileName     string `json:"file_name"`
	UID          string `json:"uid"`
	LastModified int64  `json:"last_modified"`
	MD5Hash      string `json:"md5"`
	Size         string `json:"size"`
	Version      string `json:"version"`
}

// AddFileMetaToDB adds uploaded file to database
func AddFileMetaToDB(fileName string, md5 string, uid string, contents int, version string) error {
	lastModified := strconv.FormatInt(time.Now().UnixMilli(), 10)
	_, err := dbConnection.Exec(`INSERT INTO file_meta (file_name, uid, last_modified, md5_hash, file_contents, version) VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`, fileName, uid, lastModified, md5, contents, version)

	return err
}

// GetExistingFile check if file name already exists in database and if contents match too
func GetExistingFile(fileName string, hash string, uid string) (exists bool, lastVersion string, err error) {
	rows, err := dbConnection.Query(`SELECT file_name, md5_hash, version FROM file_meta WHERE file_name = @p1 AND uid = @p2`, fileName, uid)
	if err != nil {
		return false, "", err
	}

	defer rows.Close()
	existing := make([]FileMetadata, 0)

	for rows.Next() {
		data := FileMetadata{}
		rows.Scan(&data.FileName, &data.MD5Hash, &data.Version)
		existing = append(existing, data)
	}

	if len(existing) > 0 {
		for _, file := range existing {
			if file.Version > lastVersion {
				lastVersion = file.Version
			}

			if strings.TrimSpace(file.MD5Hash) == strings.TrimSpace(hash) {
				return true, lastVersion, nil
			}
		}
		return true, "", nil
	}

	return false, "", nil
}

// ListFilesForUser returns all files uploaded by the user with given UID
func ListFilesForUser(uid string) (ret []FileMetadata, err error) {
	rows, err := dbConnection.Query(`SELECT file_name, uid, last_modified, md5_hash, file_contents, version FROM file_meta WHERE uid = @p1`, uid)
	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		data := FileMetadata{}
		rows.Scan(&data.FileName, &data.UID, &data.LastModified, &data.MD5Hash, &data.Size, &data.Version)
		ret = append(ret, data)
	}

	return
}

func RemoveBlob(uid string, fileName string, version string) error {
	_, err := dbConnection.Exec(`DELETE FROM file_meta WHERE file_name = @p1 AND version = @p2 AND uid = @p3`, fileName, version, uid)
	return err
}
