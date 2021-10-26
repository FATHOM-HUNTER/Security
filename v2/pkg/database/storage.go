package database

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

// AddFileMetaToDB adds uploaded file to database
func AddFileMetaToDB(fileName string, md5 string, uid string, contents int) error {
	lastModified := strconv.FormatInt(time.Now().UnixMilli(), 10)
	id := uuid.New().String()
	_, err := dbConnection.Exec(`INSERT INTO file_meta (id, file_name, uid, last_modified, md5_hash, file_contents, version) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)`, id, fileName, uid, lastModified, md5, contents, 1)

	return err
}
