package data

import (
	"database/sql"
)

// DeleteScan deletes a scan.
func DeleteScan(db *sql.DB, ID string) error {
	query := "DELETE FROM scans WHERE id = ?"
	_, err := db.Exec(query, ID)
	return err
}

func StoreNmapScan(db *sql.DB, userID int, resultXML string) error {
	query := "INSERT INTO scans (user_id, result_xml) VALUES (?, ?)"
	_, err := db.Exec(query, userID, resultXML)
	return err
}

func GetNampXMLScanByID(db *sql.DB, ID string) (string, error) {
	var xmlData string
	query := "SELECT result_xml FROM scans WHERE id = ?"
	err := db.QueryRow(query, ID).Scan(&xmlData) // fetches a row of data by ID, and maps the result_xml field from the database to the xmlData variable
	return xmlData, err
}

func GetAllNmapScans(db *sql.DB, userID int) (map[int]string, error) {
	results := make(map[int]string)

	query := "SELECT id, created_at FROM scans WHERE user_id = ?"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		id := 0
		createdAt := ""
		if err := rows.Scan(&id, &createdAt); err != nil {
			return nil, err
		}
		results[id] = createdAt
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
