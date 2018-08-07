package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"os"

	fdk "github.com/fnproject/fdk-go"
	_ "github.com/mattn/go-oci8"
)

//global scope for hot function support
var db *sql.DB

//gets called before main() and sets up the Db connection
func init() {
	log.Println("Inside init()")

	//can also use app config and fetch it from context - fn config app fn-oradb-go-app DB_USER sys ..etc. etc.
	/* fnctx := fdk.Context(ctx)
	initDB(fnctx.Config["DB_USER"], fnctx.Config["DB_PASSWORD"], fnctx.Config["DB_HOST"], fnctx.Config["DB_PORT"], fnctx.Config["DB_SERVICE_NAME"])*/

	initDB(getEnvVar("DB_USER", "scott"), getEnvVar("DB_PASSWORD", "tiger"), getEnvVar("DB_HOST", "localhost"), getEnvVar("DB_PORT", "1521"), getEnvVar("DB_SERVICE_NAME", "orcl"), getEnvVar("IS_SYSDBA", "false"))
}

func initDB(username, password, host, port, serviceName, isSysDBA string) {
	log.Println("Init DB connection.....")
	sysdbaPart := ""
	if isSysDBA == "true" {
		sysdbaPart = "?as=sysdba"
	}
	connString := username + "/" + password + "@" + host + ":" + port + "/" + serviceName + sysdbaPart
	//for debugging only - DO NOT log your password!
	log.Println("DB Connection string " + connString)
	db, _ = sql.Open("oci8", connString)
	err := db.Ping()

	if err != nil {
		log.Println("failed to connect due to err " + err.Error())
	} else {
		log.Println("Connected to DB...")
	}
}

func main() {

	fdk.Handle(fdk.HandlerFunc(readHandler))
}

//Employee - struct for employee details
type Employee struct {
	EmpEmail string `json:"emp_email"`
	EmpName  string `json:"emp_name"`
	EmpDept  string `json:"emp_dept"`
}

func readHandler(ctx context.Context, in io.Reader, out io.Writer) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(in)
	emp := buf.String()

	var query string
	if emp == "" {
		query = "select * from EMPLOYEES"
		log.Println("searching ALL employees..")
		rows, _ := db.Query(query)
		defer rows.Close()
		var emps []Employee
		for rows.Next() {
			var empl Employee
			rows.Scan(&empl.EmpEmail, &empl.EmpName, &empl.EmpDept)
			emps = append(emps, empl)
		}
		fdk.SetHeader(out, "Content-Type", "application/json")
		json.NewEncoder(out).Encode(emps)

	} else {
		log.Println("searching for employee " + emp)

		query = "SELECT * from EMPLOYEES WHERE EMP_EMAIL=:1"
		empl := &Employee{}
		err := db.QueryRow(query, emp).Scan(&empl.EmpEmail, &empl.EmpName, &empl.EmpDept)
		if err != nil {
			out.Write([]byte("Could not find employee " + emp))
		} else {
			//log.Println("sending JSON employee response")
			fdk.SetHeader(out, "Content-Type", "application/json")
			json.NewEncoder(out).Encode(&empl)
		}

	}

}
func getEnvVar(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val

}
