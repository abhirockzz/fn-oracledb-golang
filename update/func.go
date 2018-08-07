package main

import (
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
	fdk.Handle(fdk.HandlerFunc(updateHandler))
}

//UpdateEmployee - struct for update details
type UpdateEmployee struct {
	EmpEmail string `json:"emp_email"`
	EmpDept  string `json:"emp_dept"`
}

func updateHandler(ctx context.Context, in io.Reader, out io.Writer) {
	
	empInfo := &UpdateEmployee{}
	json.NewDecoder(in).Decode(empInfo)

	result := "Failed to update employee " + empInfo.EmpEmail

	log.Println("Updating department for employee " + empInfo.EmpEmail + " to " + empInfo.EmpDept)

	execResult, execErr := db.Exec("update EMPLOYEES set EMP_DEPT=:1 where EMP_EMAIL=:2", empInfo.EmpDept, empInfo.EmpEmail)
	log.Println("executed query...")

	if execErr != nil {
		log.Println("Unable to update employee info due to " + execErr.Error())
	} else {
		numRows, _ := execResult.RowsAffected()
		if numRows == 1 {
			result = "Updated employee " + empInfo.EmpEmail
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