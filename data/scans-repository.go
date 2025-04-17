package data

import (
	"database/sql"
	"time"
)

// Scan represents a record to be saved
type Scan struct {
	ID        int
	UserId    int
	Target    string
	ResultXML string
	CreatedAt time.Time
}

func StoreNmapScan(db *sql.DB, userId int, target string, resultXML string) error {
	query := "INSERT INTO scans (user_id, target, result_xml) VALUES (?, ?, ?)"
	_, err := db.Exec(query, userId, target, resultXML)
	return err
}

func GetNampXMLScanByID(db *sql.DB, ID string) (string, error) {
	var xmlData string
	query := "SELECT result_xml FROM scans WHERE id = ?"
	err := db.QueryRow(query, ID).Scan(&xmlData)
	return xmlData, err
}
