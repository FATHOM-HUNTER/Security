package database

import (
	"github.com/google/uuid"
)

func AddSharedEmail(fileName string, ownerUID string, sharedWithEmail string) error {
	uid := uuid.New()
	_, err := dbConnection.Exec(`INSERT INTO shares (uid, file_name, owner_uid, shared_with_email) VALUES (@p1, @p2, @p3, @p4)`, uid, fileName, ownerUID, sharedWithEmail)
	return err
}

func GetSharedEmails(fileName string, ownerUID string) (ret []string, err error) {
	rows, err := dbConnection.Query(`SELECT shared_with_email FROM shares WHERE file_name = @p1 AND owner_uid = @p2`, fileName, ownerUID)
	if err != nil {
		return
	}

	defer rows.Close()

	var data string
	for rows.Next() {
		rows.Scan(&data)
		ret = append(ret, data)
	}

	return
}

func RemoveSharedEmail(fileName string, ownerUID string, sharedWithEmail string) error {
	_, err := dbConnection.Exec(`DELETE FROM shares WHERE file_name = @p1 AND owner_uid = @p2 AND shared_with_email = @p3`, fileName, ownerUID, sharedWithEmail)
	return err
}

func CanAccessFile(UID string, fileName string, ownerUID string) bool {
	// if UID == ownerUID {
		return true
	// }

	shared, err := GetSharedEmails(fileName, ownerUID)
	if err != nil {
		return false
	}

	email, err := GetUserEmail(UID)
	if err != nil {
		return false
	}

	return contains(shared, email)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
