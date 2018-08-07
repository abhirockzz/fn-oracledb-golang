package main

import (
	"os"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"

	fdk "github.com/fnproject/fdk-go"
	_ "github.com/mattn/go-oci8"
)

//CreateEmployeeInfo - struct for employee details
type CreateEmployeeInfo struct {
	EmpEmail string `json:"emp_email"`
	EmpName  string `json:"emp_name"`
	EmpDept  string `json:"emp_dept"`
}

//global scope for hot function support
var db *sql.DB

//gets called before main() and sets up the Db connection
func init(){
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
	fdk.Handle(fdk.HandlerFunc(createHandler))
}

func createHandler(ctx context.Context, in io.Reader, out io.Writer) {

	empInfo := &CreateEmployeeInfo{}
	json.NewDecoder(in).Decode(empInfo)
	result := "Failed to create employee " + empInfo.EmpEmail

	log.Println("Creating employee ---- " + empInfo.EmpEmail)
	
	execResult, execErr := db.Exec("INSERT INTO EMPLOYEES VALUES (:1,:2,:3)", empInfo.EmpEmail, empInfo.EmpName, empInfo.EmpDept)
	log.Println("executed query...")

	if execErr != nil {
		log.Println("Unable to create employee due to - " + execErr.Error())
	} else {
		numRows, _ := execResult.RowsAffected()
		if numRows == 1 {
			result = "Created employee " + empInfo.EmpEmail
		}
	}

	out.Write([]byte(result))
}

func getEnvVar(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	} 
	return val
	
}
