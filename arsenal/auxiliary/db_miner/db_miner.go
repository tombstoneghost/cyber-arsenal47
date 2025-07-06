// DB Miner
package auxiliary

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb" // Microsoft SQL Server driver
	_ "github.com/go-sql-driver/mysql"   // MySQL driver
	_ "github.com/lib/pq"                // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3"      // SQLite driver
)

// DatabaseMiner stores database connection details
type DatabaseMiner struct {
	DBType   string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	DB       *sql.DB
}

// Connect establishes a connection to the database
func (miner *DatabaseMiner) Connect() error {
	var err error
	var dsn string

	switch miner.DBType {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", miner.User, miner.Password, miner.Host, miner.Port, miner.DBName)
		miner.DB, err = sql.Open("mysql", dsn)
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", miner.Host, miner.Port, miner.User, miner.Password, miner.DBName)
		miner.DB, err = sql.Open("postgres", dsn)
	case "sqlite":
		dsn = fmt.Sprintf("%s", miner.DBName) // For SQLite, the database name is the file path
		miner.DB, err = sql.Open("sqlite3", dsn)
	case "mssql":
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", miner.User, miner.Password, miner.Host, miner.Port, miner.DBName)
		miner.DB, err = sql.Open("sqlserver", dsn)
	default:
		return fmt.Errorf("[X] Unsupported database type: %s", miner.DBType)
	}

	if err != nil {
		return fmt.Errorf("[X] Error creating connection: %v", err)
	}

	// Check the connection
	err = miner.DB.Ping()
	if err != nil {
		return fmt.Errorf("[X] Error pinging database: %v", err)
	}
	fmt.Println("[!] Connected to database successfully!")
	return nil
}

// FetchAllTables retrieves the list of all tables in the database
func (miner *DatabaseMiner) FetchAllTables() ([]string, error) {
	var query string
	var tableName string
	var tables []string

	// Query to get all tables for different DB types
	switch miner.DBType {
	case "mysql":
		query = "SHOW TABLES"
	case "postgres":
		query = "SELECT table_name FROM information_schema.tables WHERE table_schema='public'"
	case "sqlite":
		query = "SELECT name FROM sqlite_master WHERE type='table'"
	case "mssql":
		query = "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE'"
	default:
		return nil, fmt.Errorf("[X] Unsupported database type: %s", miner.DBType)
	}

	rows, err := miner.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("[X] Error fetching tables: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("[X] Error scanning table name: %v", err)
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

// FetchTableData fetches and prints all data from a given table
func (miner *DatabaseMiner) FetchTableData(table string) error {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	rows, err := miner.DB.Query(query)
	if err != nil {
		return fmt.Errorf("[X] Error fetching data from table %s: %v", table, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("[X] Error fetching columns for table %s: %v", table, err)
	}

	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	fmt.Printf("\nTable: %s\n", table)
	fmt.Println(columns)

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return fmt.Errorf("[X] Error scanning row data: %v", err)
		}

		for _, col := range values {
			fmt.Printf("%v\t", col)
		}
		fmt.Println()
	}

	return nil
}

// Run executes the database mining process
func (miner *DatabaseMiner) Mine() {
	// Fetch all tables
	tables, err := miner.FetchAllTables()
	if err != nil {
		log.Fatalf("[X] Failed to retrieve tables: %v", err)
	}

	fmt.Println("[!] Tables found:", tables)

	// Fetch data from each table
	for _, table := range tables {
		err := miner.FetchTableData(table)
		if err != nil {
			log.Printf("[X] Error fetching data from table %s: %v", table, err)
		}
	}
}

func DBMiner_Init(dbType string, host string, port int, username string, password string, dbName string) {
	fmt.Println("[>] DBType: ", dbType)
	fmt.Println("[>] Host: ", host)
	fmt.Println("[>] Port: ", port)
	fmt.Println("[>] Username: ", username)
	fmt.Println("[>] Password: ", password)
	fmt.Println("[>] DBName: ", dbName)

	dbMiner := &DatabaseMiner{
		DBType:   dbType, // Choose: "mysql", "postgres", "sqlite", "mssql"
		Host:     host,
		Port:     port, // MySQL: 3306, PostgreSQL: 5432, MSSQL: 1433, SQLite: No port needed
		User:     username,
		Password: password,
		DBName:   dbName, // For SQLite, this would be the file path (e.g., "./mydb.db")
	}

	// Connect to the database
	err := dbMiner.Connect()
	if err != nil {
		log.Fatalf("[X] Could not connect to the database: %v", err)
	}

	// Run the miner to fetch all tables and their data
	dbMiner.Mine()
}
