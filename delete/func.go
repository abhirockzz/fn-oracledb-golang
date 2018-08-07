package main

import (
	"bytes"
	"context"
	"database/sql"
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

	fdk.Handle(fdk.HandlerFunc(deleteHandler))
}


func deleteHandler(ctx context.Context, in io.Reader, out io.Writer) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(in)
	emp := buf.String()

	result := "Failed to delete employee " + emp

	log.Println("Deleting employee " + emp)

	execResult, execErr := db.Exec("delete from EMPLOYEES where EMP_EMAIL=:EMP_EMAIL", emp)
	log.Println("executed query...")

	if execErr != nil {
		log.Println("Unable to delete employee due to " + execErr.Error())
	} else {
		numRowsAffected, _ := execResult.RowsAffected()
		if numRowsAffected == 1 {
			result = "Deleted employee " + emp
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