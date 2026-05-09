// Seed tool — creates a confirmed-and-active account, optionally bound to a
// company as admin so the React app's launcher resolves to a real company
// UUID instead of "null". Invite-only product means no public signup; an
// admin always seeds users directly.
//
// Usage:
//   go run ./tools/seed -email manager@acme.local -password acmedev -name "Acme Manager"
//   go run ./tools/seed -email manager@acme.local -password acmedev -company "Acme Diner" -team "Front of house"
//
// Connects to the docker-compose mysql container by default (localhost:3306).
// Idempotent: existing rows with matching email/company name are deleted first.
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"v2.staffjoy.com/crypto"
)

const (
	defaultMySQLHost = "127.0.0.1"
	defaultMySQLPort = "3306"
	defaultMySQLUser = "root"
	defaultMySQLPass = "devroot"
	defaultTimezone  = "America/Los_Angeles"
	defaultDayStart  = "monday"
)

func main() {
	email := flag.String("email", "", "email (required)")
	password := flag.String("password", "", "plaintext password (required, ≥6 chars)")
	name := flag.String("name", "", "display name (defaults to email local-part)")
	support := flag.Bool("support", false, "set the support flag (sudo across the app)")
	company := flag.String("company", "", "company name; if set, create the company and make this user admin")
	team := flag.String("team", "Default Team", "team name to create inside the company (only if -company set)")
	tz := flag.String("timezone", defaultTimezone, "company default timezone")
	mysqlHost := flag.String("mysql-host", defaultMySQLHost, "mysql host")
	mysqlPort := flag.String("mysql-port", defaultMySQLPort, "mysql port")
	mysqlUser := flag.String("mysql-user", defaultMySQLUser, "mysql user")
	mysqlPass := flag.String("mysql-pass", defaultMySQLPass, "mysql password")
	flag.Parse()

	if *email == "" || *password == "" {
		fatal("-email and -password are required")
	}
	if len(*password) < 6 {
		fatal("password must be at least 6 chars")
	}
	if *name == "" {
		if i := strings.Index(*email, "@"); i > 0 {
			*name = (*email)[:i]
		} else {
			*name = *email
		}
	}

	accountDSN := mysqlDSN(*mysqlUser, *mysqlPass, *mysqlHost, *mysqlPort, "account")
	companyDSN := mysqlDSN(*mysqlUser, *mysqlPass, *mysqlHost, *mysqlPort, "company")

	userUUID := seedAccount(accountDSN, *email, *password, *name, *support)
	fmt.Printf("✓ account seeded: %s (%s, support=%v)\n", *email, userUUID, *support)

	if *company != "" {
		companyUUID, teamUUID := seedCompany(companyDSN, userUUID, *company, *team, *tz)
		fmt.Printf("✓ company seeded: %s (%s)\n", *company, companyUUID)
		fmt.Printf("✓ team seeded:    %s (%s)\n", *team, teamUUID)
		fmt.Printf("✓ admin link:     %s ↔ %s\n", *email, *company)
		fmt.Printf("✓ directory entry added\n")
	}

	fmt.Println()
	fmt.Println("Login at http://localhost:8080/login/ with:")
	fmt.Printf("  email:    %s\n", *email)
	fmt.Printf("  password: %s\n", *password)
}

func seedAccount(dsn, email, password, name string, support bool) string {
	db := mustOpen(dsn)
	defer db.Close()

	uuid, err := crypto.NewUUID()
	if err != nil {
		fatal("uuid: " + err.Error())
	}
	salt, err := crypto.NewSalt()
	if err != nil {
		fatal("salt: " + err.Error())
	}
	hash, err := crypto.HashPassword(salt, []byte(password))
	if err != nil {
		fatal("hash: " + err.Error())
	}

	must(db.Exec("DELETE FROM account WHERE email=?", email))
	supportInt := 0
	if support {
		supportInt = 1
	}
	must(db.Exec(`
INSERT INTO account
  (uuid, email, name, confirmed_and_active, member_since, password_hash, password_salt, support, phonenumber, photo_url)
VALUES
  (?, ?, ?, 1, ?, ?, ?, ?, '', '')`,
		uuid.String(), email, name, time.Now().UTC().Format("2006-01-02 15:04:05"),
		hash, salt, supportInt))

	return uuid.String()
}

func seedCompany(dsn, userUUID, companyName, teamName, tz string) (string, string) {
	db := mustOpen(dsn)
	defer db.Close()

	companyUUID, err := crypto.NewUUID()
	if err != nil {
		fatal("uuid: " + err.Error())
	}
	teamUUID, err := crypto.NewUUID()
	if err != nil {
		fatal("uuid: " + err.Error())
	}

	must(db.Exec("DELETE FROM company WHERE name=?", companyName))
	must(db.Exec("DELETE FROM admin WHERE user_uuid=?", userUUID))
	must(db.Exec("DELETE FROM directory WHERE user_uuid=?", userUUID))

	must(db.Exec(`
INSERT INTO company (uuid, name, archived, default_timezone, default_day_week_starts)
VALUES (?, ?, 0, ?, ?)`,
		companyUUID.String(), companyName, tz, defaultDayStart))

	must(db.Exec(`
INSERT INTO team (uuid, company_uuid, name, archived, timezone, day_week_starts, color)
VALUES (?, ?, ?, 0, ?, ?, '48B7AB')`,
		teamUUID.String(), companyUUID.String(), teamName, tz, defaultDayStart))

	must(db.Exec(`
INSERT INTO admin (company_uuid, user_uuid) VALUES (?, ?)`,
		companyUUID.String(), userUUID))

	must(db.Exec(`
INSERT INTO directory (company_uuid, user_uuid, internal_id) VALUES (?, ?, '')`,
		companyUUID.String(), userUUID))

	return companyUUID.String(), teamUUID.String()
}

func mysqlDSN(user, pass, host, port, db string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4", user, pass, host, port, db)
}

func mustOpen(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fatal("open: " + err.Error())
	}
	if err := db.Ping(); err != nil {
		fatal("ping: " + err.Error())
	}
	return db
}

func must(_ sql.Result, err error) {
	if err != nil {
		fatal("sql: " + err.Error())
	}
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, "ERROR:", msg)
	os.Exit(1)
}
