package db

import (
	"database/sql"
)

type AitDbProp struct {
	AitNo             string
	ID                string
	Profile           string
	DbType            string
	MachineName       string
	DbName            string
	SchemaName        string
	UserID            string
	Password          string
	TopicName         string
	JdbcUrl           string
	NoOfConnection    *int
	EmailID           string
	IsActive          bool
	LastUpdatedUser   string
	SddActive         bool
	AimlIsActive      bool
	TableList         string
	ExecStatus        string
	AitCadence        string
	IdwSdd            bool
	IdwUdd            bool
	IedpsSdd          bool
	ReportTopicName   string
	FunnelUdd         bool
	FunnelSdd         bool
	AimlSdd           bool
	AimlUdd           bool
	FunnelDestination string
	FunnelDiscovery   bool
	AimlDiscovery     bool
	IdwDiscovery      bool
	IedpsDiscovery    bool
	AimlValidation    bool
	FftDestination    string
	AitNum            *int
	JdbcConStr        string
	Lob               string
	Environment       string
}

func SelectAitDbProps(db *sql.DB, conditions map[string]interface{}, orderBy string) ([]AitDbProp, error) {
	query := "SELECT * FROM AIT_DBPROP1"

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

	var props []AitDbProp
	for rows.Next() {
		var prop AitDbProp
		err := rows.Scan(
			&prop.AitNo, &prop.ID, &prop.Profile, &prop.DbType, &prop.MachineName,
			&prop.DbName, &prop.SchemaName, &prop.UserID, &prop.Password, &prop.TopicName,
			&prop.JdbcUrl, &prop.NoOfConnection, &prop.EmailID, &prop.IsActive,
			&prop.LastUpdatedUser, &prop.SddActive, &prop.AimlIsActive, &prop.TableList,
			&prop.ExecStatus, &prop.AitCadence, &prop.IdwSdd, &prop.IdwUdd,
			&prop.IedpsSdd, &prop.ReportTopicName, &prop.FunnelUdd, &prop.FunnelSdd,
			&prop.AimlSdd, &prop.AimlUdd, &prop.FunnelDestination, &prop.FunnelDiscovery,
			&prop.AimlDiscovery, &prop.IdwDiscovery, &prop.IedpsDiscovery,
			&prop.AimlValidation, &prop.FftDestination, &prop.AitNum,
			&prop.JdbcConStr, &prop.Lob, &prop.Environment,
		)
		if err != nil {
			return nil, err
		}
		props = append(props, prop)
	}
	return props, nil
}

func InsertAitDbProp(db *sql.DB, prop *AitDbProp) error {
	query := `
        INSERT INTO AIT_DBPROP1 (
            AIT_NO, ID, PROFILE, DB_TYPE, MACHINE_NAME, DB_NAME, SCHEMA_NAME,
            USER_ID, PASS_WORD, TOPIC_NAME, JDBC_URL, NO_OF_CONNECTION,
            EMAIL_ID, IS_ACTIVE, LAST_UPDATED_USER, SDD_Active, AIML_IS_ACTIVE,
            TABLE_LIST, EXEC_STATUS, AIT_CADENCE, IDW_SDD, IDW_UDD,
            IEDPS_SDD, REPORT_TOPIC_NAME, FUNNEL_UDD, FUNNEL_SDD,
            AIML_SDD, AIML_UDD, FUNNEL_DESTINATION, FUNNEL_DISCOVERY,
            AIML_DISCOVERY, IDW_DISCOVERY, IEDPS_DISCOVERY, AIML_VALIDATION,
            FFT_DESTINATION, AIT_NUM, JDBC_CON_STR, LOB, ENVIRONMENT
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
                 ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	_, err := db.Exec(query,
		prop.AitNo, prop.ID, prop.Profile, prop.DbType, prop.MachineName,
		prop.DbName, prop.SchemaName, prop.UserID, prop.Password, prop.TopicName,
		prop.JdbcUrl, prop.NoOfConnection, prop.EmailID, prop.IsActive,
		prop.LastUpdatedUser, prop.SddActive, prop.AimlIsActive, prop.TableList,
		prop.ExecStatus, prop.AitCadence, prop.IdwSdd, prop.IdwUdd,
		prop.IedpsSdd, prop.ReportTopicName, prop.FunnelUdd, prop.FunnelSdd,
		prop.AimlSdd, prop.AimlUdd, prop.FunnelDestination, prop.FunnelDiscovery,
		prop.AimlDiscovery, prop.IdwDiscovery, prop.IedpsDiscovery,
		prop.AimlValidation, prop.FftDestination, prop.AitNum,
		prop.JdbcConStr, prop.Lob, prop.Environment,
	)
	return err
}

func UpdateAitDbProp(db *sql.DB, updates map[string]interface{}, conditions map[string]interface{}) error {
	query := "UPDATE AIT_DBPROP1"

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

func DeleteAitDbProp(db *sql.DB, conditions map[string]interface{}) error {
	query := "DELETE FROM AIT_DBPROP1"

	whereClause, values := buildWhereClause(conditions)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	_, err := db.Exec(query, values...)
	return err
}
