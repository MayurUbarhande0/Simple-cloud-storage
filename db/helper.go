package helper

import (
	"database/sql"
)

type Db struct {
	db *sql.DB
}

func NewDb(db *sql.DB) *Db {
	return &Db{db: db}
}

// Check returns true if the file_id exists
func (db *Db) Check(file_id string) (bool, error) {
	var exists int
	// Using SELECT 1 is faster for existence checks
	err := db.db.QueryRow("SELECT 1 FROM files WHERE file_id = ? LIMIT 1", file_id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Insert adds a new file record
func (db *Db) Insert(user_id string, file_id string, path string) error {
	_, err := db.db.Exec("INSERT INTO files (user_id, file_id, path) VALUES (?, ?, ?)", user_id, file_id, path)
	return err
}

// Delete removes a file record
func (db *Db) Delete(file_id string) error {
	_, err := db.db.Exec("DELETE FROM files WHERE file_id = ?", file_id)
	return err
}

// GetPath retrieves the storage path for a specific file and user
func (db *Db) GetPath(user_id string, file_id string) (string, error) {
	var path string
	// Using AND ensures standard SQL compliance
	err := db.db.QueryRow("SELECT path FROM files WHERE file_id = ? AND user_id = ?", file_id, user_id).Scan(&path)
	if err != nil {
		return "", err
	}
	return path, nil
}
