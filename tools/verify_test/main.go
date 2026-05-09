// One-off: connect to mysql, fetch hash+salt, attempt to verify a known password.
package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"v2.staffjoy.com/crypto"
)

func main() {
	dsn := os.Getenv("MYSQL_CONFIG")
	if dsn == "" {
		fmt.Fprintln(os.Stderr, "ERROR: set MYSQL_CONFIG env (no scheme prefix; e.g. root:devroot@tcp(127.0.0.1:3306)/account)")
		os.Exit(1)
	}
	email := os.Args[1]
	password := os.Args[2]

	db, err := sql.Open("mysql", dsn+"?parseTime=true")
	if err != nil {
		fmt.Fprintln(os.Stderr, "open:", err)
		os.Exit(1)
	}
	defer db.Close()

	var uuid, dbHash, salt sql.NullString
	var confirmedAndActive sql.NullBool
	err = db.QueryRow("SELECT uuid, password_hash, password_salt, confirmed_and_active FROM account WHERE email=?", email).Scan(&uuid, &dbHash, &salt, &confirmedAndActive)
	if err != nil {
		fmt.Fprintln(os.Stderr, "query:", err)
		os.Exit(1)
	}

	fmt.Printf("uuid=%q\n", uuid.String)
	fmt.Printf("confirmed=%v\n", confirmedAndActive.Bool)
	fmt.Printf("hash_len=%d hash=%q\n", len(dbHash.String), dbHash.String)
	fmt.Printf("salt_len=%d\n", len(salt.String))

	err = crypto.CheckPasswordHash([]byte(dbHash.String), []byte(salt.String), []byte(password))
	if err != nil {
		fmt.Println("VERIFY FAIL:", err)
		os.Exit(2)
	}
	fmt.Println("VERIFY OK")
}
