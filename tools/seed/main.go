// Seed tool — emits SQL inserts for a confirmed-and-active account so we can
// log in without going through the signup → email-activation flow. Useful
// for local dev (invite-only product means no public signup; an admin will
// always create accounts directly).
//
// Usage:
//   go run ./tools/seed -email manager@acme.local -password acmedev -name "Acme Manager"
// then pipe through mysql:
//   go run ./tools/seed ... | docker compose exec -T mysql mysql -uroot -p$(grep MYSQL_ROOT_PASSWORD .env | cut -d= -f2) account
//
// Idempotent: existing rows with the same email are deleted first.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"v2.staffjoy.com/crypto"
)

func main() {
	email := flag.String("email", "", "email (required)")
	password := flag.String("password", "", "plaintext password (required, ≥6 chars)")
	name := flag.String("name", "", "display name (defaults to email local-part)")
	support := flag.Bool("support", false, "set the support flag (sudo across the app)")
	flag.Parse()

	if *email == "" || *password == "" {
		fmt.Fprintln(os.Stderr, "ERROR: -email and -password are required")
		flag.Usage()
		os.Exit(1)
	}
	if len(*password) < 6 {
		fmt.Fprintln(os.Stderr, "ERROR: password must be at least 6 chars")
		os.Exit(1)
	}
	if *name == "" {
		if i := strings.Index(*email, "@"); i > 0 {
			*name = (*email)[:i]
		} else {
			*name = *email
		}
	}

	uuid, err := crypto.NewUUID()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR generating uuid:", err)
		os.Exit(1)
	}
	salt, err := crypto.NewSalt()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR generating salt:", err)
		os.Exit(1)
	}
	hash, err := crypto.HashPassword(salt, []byte(*password))
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR hashing password:", err)
		os.Exit(1)
	}

	supportInt := 0
	if *support {
		supportInt = 1
	}

	// MySQL accepts BINARY columns via the X'…' hex literal form — that
	// avoids any escaping pitfalls with raw bytes.
	memberSince := time.Now().UTC().Format("2006-01-02 15:04:05")
	fmt.Printf(`-- seed account: %s
DELETE FROM account WHERE email=%s;
INSERT INTO account
  (uuid, email, name, confirmed_and_active, member_since, password_hash, password_salt, support, phonenumber, photo_url)
VALUES
  (%s, %s, %s, 1, %s, X'%x', X'%x', %d, '', '');
SELECT uuid, email, name, support, confirmed_and_active FROM account WHERE email=%s;
`,
		*email,
		quote(*email),
		quote(uuid.String()), quote(*email), quote(*name), quote(memberSince),
		hash, salt, supportInt,
		quote(*email),
	)
}

// quote returns a single-quoted SQL literal with embedded single quotes
// doubled.  Sufficient for the seed-tool inputs we control.
func quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
