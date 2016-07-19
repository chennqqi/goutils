package mysqldump

import (
	"bufio"
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	csv "github.com/JensRantil/go-csv"
	"github.com/JensRantil/go-csv/dialect"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pquerna/ffjson/ffjson"
)

const (
	TypeCSV  TypeFlag = 1
	TypeJSON TypeFlag = 0
)

// Queryable interface that matches sql.DB and sql.Tx.
type queryable interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func dump2csv(tables []string, db queryable, outputDir string, compressOut, skipHeader bool, csvDialect *csv.Dialect) error {
	for _, table := range tables {
		err := dumpTable(table, db, outputDir, compressOut, skipHeader, csvDialect)
		if err != nil {
			fmt.Printf("Error dumping %s: %s\n", table, err)
		}
	}
	return nil
}

func dump2JSON(tables []string, db queryable, outputDir string) error {
	for _, table := range tables {
		err := dumpTable2JSON(table, db, outputDir)
		if err != nil {
			fmt.Printf("Error dumping %s: %s\n", table, err)
		}
	}
	return nil
}

func dumpTable(table string, db queryable, outputDir string,
	compressOut, skipHeader bool, csvDialect *csv.Dialect) error {

	fname := outputDir + "/" + table + ".csv"
	if compressOut {
		fname = fname + ".gz"
	}

	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	var out io.Writer
	if compressOut {
		gzout := gzip.NewWriter(f)
		defer gzout.Close()
		out = gzout
	} else {
		out = f
	}

	w := csv.NewDialectWriter(out, *csvDialect)

	rows, err := db.Query("SELECT * FROM " + table) // Couldn't get placeholder expansion to work here
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	if !skipHeader {
		err = w.Write(columns) // Header
		if err != nil {
			return err
		}
	}

	for rows.Next() {
		// Shamelessly ripped (and modified) from http://play.golang.org/p/jxza3pbqq9

		// Create interface set
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// Scan for arbitrary values
		err = rows.Scan(scanArgs...)
		if err != nil {
			return err
		}

		// Print data
		csvData := make([]string, 0, len(values))
		for _, value := range values {
			switch value.(type) {
			default:
				s := fmt.Sprintf("%s", value)
				csvData = append(csvData, string(s))
			}
		}
		err = w.Write(csvData)
		if err != nil {
			return err
		}
	}

	w.Flush()
	err = w.Error()
	if err != nil {
		return err
	}

	return nil
}

func dumpTable2JSON(table string, db queryable, outputDir string) error {
	fname := outputDir + "/" + table + ".tjson"

	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	rows, err := db.Query("SELECT * FROM " + table) // Couldn't get placeholder expansion to work here
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	mjson := make(map[string]interface{}, len(columns))

	for rows.Next() {
		// Shamelessly ripped (and modified) from http://play.golang.org/p/jxza3pbqq9

		// Create interface set
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// Scan for arbitrary values
		err = rows.Scan(scanArgs...)
		if err != nil {
			return err
		}

		for idx, value := range values {
			switch value.(type) {
			default:
				s := fmt.Sprintf("%s", value)
				mjson[columns[idx]] = s
			}
		}
		txt, err := ffjson.Marshal(mjson)
		if err != nil {
			log.Panic(err)
		}
		_, err = w.Write(txt)
		w.WriteByte('\n')
		if err != nil {
			log.Panic(err)
		}
	}

	w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func getTables(db queryable) ([]string, error) {
	tables := make([]string, 0, 10)
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var table string
		rows.Scan(&table)
		tables = append(tables, table)
	}
	return tables, nil
}

type TypeFlag int

type DumpOpt struct {
	DumpDir        string   `json:"dumpdir"`
	TableNames     []string `json:"table_names"`
	Host           string   `json:"db_host"`
	User           string   `json:"db_user"`
	DBName         string   `json:"db_name"`
	Password       string   `json:"db_password"`
	Type           TypeFlag `json:"export_type"`
	UseTrasnAction bool     `json:"transaction_enable"`
}

func MysqlDump(opt DumpOpt) error {
	compressFiles := false
	skipHeader := false

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s", opt.User, opt.Password,
		opt.Host, opt.DBName)
	db, err := sql.Open("mysql", dbUrl)

	if err != nil {
		return fmt.Errorf("Could not connect to server", err)
	}
	defer db.Close()

	var q queryable
	if opt.UseTrasnAction {
		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}
		defer tx.Rollback()
		q = tx
	} else {
		q = db
	}

	var tables []string
	if len(opt.TableNames) < 1 {
		tables, err = getTables(q)
		if err != nil {
			return err
		}
	} else {
		tables = opt.TableNames
	}
	os.Mkdir(opt.DumpDir, 0777)

	switch opt.Type {
	case TypeJSON:
		err = dump2JSON(tables, q, opt.DumpDir)

	case TypeCSV:
		dialectBuilder := dialect.FromCommandLine()
		csvDialect, err := dialectBuilder.Dialect()
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		err = dump2csv(tables, q, opt.DumpDir, compressFiles, skipHeader, csvDialect)
	}

	if err != nil {
		return err
	}
	return nil
}
