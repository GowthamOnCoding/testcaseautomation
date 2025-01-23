package dbutils

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type KafkaStatRecord struct {
	AitNo             string
	StartTime         *time.Time
	Event             string
	ProcessNo         string
	TableName         string
	TotalRows         *int64
	TotalFftInstances *int64
	TotalMessages     *int
	EndTime           *time.Time
	Duration          *float64
	SeqNo             int
	Remarks           string
	LastUpdated       time.Time
	Status            string
	TopicName         string
	GroupID           string
	DbType            string
	SchemaName        string
	MachineName       string
	Profile           string
	ConfigID          string
	LastUpdatedUser   string
	SchedulerID       string
}

func SelectKafkaStats(db *sql.DB, conditions map[string]interface{}, orderBy string) ([]KafkaStatRecord, error) {
	query := "SELECT * FROM KAFKA_STAT_TMP"

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

	var records []KafkaStatRecord
	for rows.Next() {
		var record KafkaStatRecord
		err := rows.Scan(
			&record.AitNo, &record.StartTime, &record.Event, &record.ProcessNo,
			&record.TableName, &record.TotalRows, &record.TotalFftInstances,
			&record.TotalMessages, &record.EndTime, &record.Duration, &record.SeqNo,
			&record.Remarks, &record.LastUpdated, &record.Status, &record.TopicName,
			&record.GroupID, &record.DbType, &record.SchemaName, &record.MachineName,
			&record.Profile, &record.ConfigID, &record.LastUpdatedUser, &record.SchedulerID,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func InsertKafkaStat(db *sql.DB, record *KafkaStatRecord) error {
	query := `
		INSERT INTO KAFKA_STAT (
			AIT_NO, START_TIME, EVENT, PROCESS_NO, TABLE_NAME,
			TOTAL_ROWS, TOTAL_FFT_INSTANCES, TOTAL_MESSAGES, END_TIME,
			DURATION, SEQ_NO, REMARKS, LAST_UPDATED, STATUS,
			TOPIC_NAME, GROUP_ID, DB_TYPE, SCHEMA_NAME, MACHINE_NAME,
			PROFILE, CONFIG_ID, LAST_UPDATED_USER, SCHEDULERID
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, GETDATE(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query,
		record.AitNo, record.StartTime, record.Event, record.ProcessNo,
		record.TableName, record.TotalRows, record.TotalFftInstances,
		record.TotalMessages, record.EndTime, record.Duration, record.SeqNo,
		record.Remarks, record.Status, record.TopicName,
		record.GroupID, record.DbType, record.SchemaName, record.MachineName,
		record.Profile, record.ConfigID, record.LastUpdatedUser, record.SchedulerID,
	)
	return err
}

func UpdateKafkaStat(db *sql.DB, updates map[string]interface{}, conditions map[string]interface{}) error {
	query := "UPDATE KAFKA_STAT"

	setClause, setValues := buildSetClause(updates)
	whereClause, whereValues := buildWhereClause(conditions)

	query += " SET " + setClause
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	values := append(setValues, whereValues...)
	_, err := db.Exec(query, values...)
	return err
}

func DeleteKafkaStat(db *sql.DB, conditions map[string]interface{}) error {
	query := "DELETE FROM KAFKA_STAT"

	whereClause, values := buildWhereClause(conditions)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	_, err := db.Exec(query, values...)
	return err
}

func buildWhereClause(conditions map[string]interface{}) (string, []interface{}) {
	var clauses []string
	var values []interface{}

	for key, value := range conditions {
		clauses = append(clauses, fmt.Sprintf("%s = ?", key))
		values = append(values, value)
	}

	return strings.Join(clauses, " AND "), values
}

func buildSetClause(updates map[string]interface{}) (string, []interface{}) {
	var clauses []string
	var values []interface{}

	for key, value := range updates {
		clauses = append(clauses, fmt.Sprintf("%s = ?", key))
		values = append(values, value)
	}

	return strings.Join(clauses, ", "), values
}
