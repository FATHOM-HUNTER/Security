package database

import (
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type FileMetadata struct {
	ID           string `json:"id"`
	FileName     string `json:"file_name"`
	UID          string `json:"uid"`
	LastModified int64  `json:"last_modified"`
	MD5Hash      string `json:"md5"`
	Size         string `json:"size"`
	Version      string `json:"version"`
}

// AddFileMetaToDB adds uploaded file to database
func AddFileMetaToDB(fileName string, md5 string, uid string, contents int) error {
	lastModified := strconv.FormatInt(time.Now().UnixMilli(), 10)
	id := uuid.New().String()
	_, err := dbConnection.Exec(`INSERT INTO file_meta (id, file_name, uid, last_modified, md5_hash, file_contents, version) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)`, id, fileName, uid, lastModified, md5, contents, 1)

	return err
}

// ListFilesForUser returns all files uploaded by the user with given UID
func ListFilesForUser(uid string) (ret []FileMetadata, err error) {
	rows, err := dbConnection.Query(`SELECT id, file_name, uid, last_modified, md5_hash, file_contents, version FROM file_meta WHERE uid = @p1`, uid)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		data := FileMetadata{}
		rows.Scan(&data.ID, &data.FileName, &data.UID, &data.LastModified, &data.MD5Hash, &data.Size, &data.Version)
		ret = append(ret, data)
	}

	return
}
