package data

import (
	"database/sql"

	"github.com/KarolinaLop/dp/models"
)

// DeleteScan deletes a scan.
func DeleteScan(db *sql.DB, ID string) error {
	query := "DELETE FROM scans WHERE id = ?"
	_, err := db.Exec(query, ID)
	return err
}

func CreateScan(db *sql.DB, scanStatus string, userID int) (int64, error) {
	query := "INSERT INTO scans (scan_status, user_id) VALUES (?,?)"
	res, err := db.Exec(query, scanStatus, userID)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func UpdateScan(db *sql.DB, ID int64, status string, PID int, resultXML string) error {
	query := "UPDATE scans SET scan_status = ?, pid = ?, result_xml = ?  WHERE id = ?"
	_, err := db.Exec(query, status, PID, resultXML, ID)
	return err
}

func GetNampXMLScanByID(db *sql.DB, ID string) (string, error) {
	var xmlData string
	query := "SELECT result_xml FROM scans WHERE id = ?"
	err := db.QueryRow(query, ID).Scan(&xmlData) // fetches a row of data by ID, and maps the result_xml field from the database to the xmlData variable
	return xmlData, err
}

func GetAllNmapScans(db *sql.DB, userID int) ([]models.Scan, error) {
	results := []models.Scan{}

	query := "SELECT id, created_at, scan_status FROM scans WHERE user_id = ?"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := models.Scan{}
		if err := rows.Scan(&s.ID, &s.Timestamp, &s.Status); err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func GetScanStatus(db *sql.DB, ID string) (status string, err error) {
	query := "SELECT scan_status FROM scans WHERE id = ?"
	if err = db.QueryRow(query, ID).Scan(&status); err != nil {
		return "", err
	}

	return status, nil
}
