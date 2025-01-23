package db

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"testcaseautomation/constants"
	"time"
)

func ConnectToSQLServer() (*sql.DB, error) {
	connString := "server=wva60bddpadb11v.corp.bankofamerica.com;port=15001;trusted_connection=yes;user id=CORP\\zsppuwid;password=c5prrP4s;database=idwdb"
	conn, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func UpdateTable(db *sql.DB, tableName string, whereCondition string, updateValues map[string]interface{}) error {
	query := "UPDATE " + tableName + " SET "
	values := make([]interface{}, 0)
	for key, value := range updateValues {
		query += key + " = ?, "
		values = append(values, value)
	}
	query = query[:len(query)-2] // to remove the last comma and space
	query += " WHERE " + whereCondition

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFromTable(db *sql.DB, tableName string, whereCondition string, args ...interface{}) error {
	query := "DELETE FROM " + tableName + " WHERE " + whereCondition

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}
	return nil
}

func InsertScanWindowRecords(db *sql.DB, aitNo, dbType, startTime, endTime string) error {
	// Delete existing records for the specified AIT_NO
	_, err := db.Exec("DELETE FROM AIT_SCAN_WINDOW WHERE AIT_NO = ?", aitNo)
	if err != nil {
		return err
	}

	windowStartTime := "00:00:00"
	if startTime != "" {
		windowStartTime = startTime
	}

	windowEndTime := "23:59:00"
	if endTime != "" {
		windowEndTime = endTime
	}

	lastUpdatedUser := "testframework"
	profile := "public"

	// Insert new records
	for _, day := range constants.Days {
		currentTime := time.Now().Format("2006-01-02 15:04:05") // Get current timestamp
		_, err := db.Exec("INSERT INTO AIT_SCAN_WINDOW (AIT_NO, DB_TYPE, SCAN_DAY, START_TIME_EST, END_TIME_EST, LAST_UPDATED, LAST_UPDATED_USER, PROFILE) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			aitNo, dbType, day, windowStartTime, windowEndTime, currentTime, lastUpdatedUser, profile)
		if err != nil {
			return err
		}
	}
	return nil
}
