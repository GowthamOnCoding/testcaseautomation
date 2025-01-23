package dbutils

import (
	"database/sql"
	"time"
)

type AitScanWindow struct {
	AitNo           string
	DbType          string
	ScanDay         string
	StartTimeEst    time.Time
	EndTimeEst      time.Time
	LastUpdated     *time.Time
	LastUpdatedUser string
	Profile         string
}

func SelectAitScanWindows(db *sql.DB, conditions map[string]interface{}, orderBy string) ([]AitScanWindow, error) {
	query := "SELECT * FROM AIT_SCAN_WINDOW1"

	whereClause, values := buildWhereClause(conditions)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	if orderBy != "" {
		query += " ORDER BY " + orderBy
	}

	rows, err := db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var windows []AitScanWindow
	for rows.Next() {
		var window AitScanWindow
		err := rows.Scan(
			&window.AitNo,
			&window.DbType,
			&window.ScanDay,
			&window.StartTimeEst,
			&window.EndTimeEst,
			&window.LastUpdated,
			&window.LastUpdatedUser,
			&window.Profile,
		)
		if err != nil {
			return nil, err
		}
		windows = append(windows, window)
	}
	return windows, nil
}

func InsertAitScanWindow(db *sql.DB, window *AitScanWindow) error {
	query := `
        INSERT INTO AIT_SCAN_WINDOW1 (
            AIT_NO, DB_TYPE, SCAN_DAY, START_TIME_EST, END_TIME_EST,
            LAST_UPDATED, LAST_UPDATED_USER, PROFILE
        ) VALUES (?, ?, ?, ?, ?, GETDATE(), ?, ?)
    `

	_, err := db.Exec(query,
		window.AitNo,
		window.DbType,
		window.ScanDay,
		window.StartTimeEst,
		window.EndTimeEst,
		window.LastUpdatedUser,
		window.Profile,
	)
	return err
}

func UpdateAitScanWindow(db *sql.DB, updates map[string]interface{}, conditions map[string]interface{}) error {
	query := "UPDATE AIT_SCAN_WINDOW1"

	setClause, setValues := buildSetClause(updates)
	whereClause, whereValues := buildWhereClause(conditions)

	query += " SET " + setClause + ", LAST_UPDATED = GETDATE()"
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	values := append(setValues, whereValues...)
	_, err := db.Exec(query, values...)
	return err
}

func DeleteAitScanWindow(db *sql.DB, conditions map[string]interface{}) error {
	query := "DELETE FROM AIT_SCAN_WINDOW1"

	whereClause, values := buildWhereClause(conditions)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	_, err := db.Exec(query, values...)
	return err
}

func BulkInsertAitScanWindows(db *sql.DB, windows []AitScanWindow) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := `
        INSERT INTO AIT_SCAN_WINDOW1 (
            AIT_NO, DB_TYPE, SCAN_DAY, START_TIME_EST, END_TIME_EST,
            LAST_UPDATED, LAST_UPDATED_USER, PROFILE
        ) VALUES (?, ?, ?, ?, ?, GETDATE(), ?, ?)
    `

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, window := range windows {
		_, err := stmt.Exec(
			window.AitNo,
			window.DbType,
			window.ScanDay,
			window.StartTimeEst,
			window.EndTimeEst,
			window.LastUpdatedUser,
			window.Profile,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
