// Seed tool — creates dev fixtures so the manager UI has something useful
// to display. Two modes:
//
//   bare     -email + -password : creates one confirmed account, optionally
//            -company / -team   bound to a company as admin (no employees,
//                               no shifts).
//
//   demo     -demo              : wipes + rebuilds the full Acme Diner
//                               fixture: 3 teams, 5 jobs, 7 employees
//                               (manager + 6 workers), ~30 shifts spread
//                               across the next 14 days. Skips other flags.
//
// Connects to the docker-compose mysql container by default (localhost:3306).
// Idempotent: existing rows for the email / company name are deleted first.
//
// Usage:
//   go run ./tools/seed -demo
//   go run ./tools/seed -email manager@acme.local -password acmedev -support
//   go run ./tools/seed -email m@x.com -password p -company "Acme" -team "FOH"
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"math/rand"
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
	demo := flag.Bool("demo", false, "wipe and rebuild the full Acme Diner demo fixture")
	email := flag.String("email", "", "email (required unless -demo)")
	password := flag.String("password", "", "plaintext password (required unless -demo, ≥6 chars)")
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

	accountDSN := mysqlDSN(*mysqlUser, *mysqlPass, *mysqlHost, *mysqlPort, "account")
	companyDSN := mysqlDSN(*mysqlUser, *mysqlPass, *mysqlHost, *mysqlPort, "company")

	if *demo {
		seedDemo(accountDSN, companyDSN, *tz)
		return
	}

	if *email == "" || *password == "" {
		fatal("-email and -password are required (or use -demo)")
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

	userUUID := seedAccount(accountDSN, *email, *password, *name, *support)
	fmt.Printf("✓ account seeded: %s (%s, support=%v)\n", *email, userUUID, *support)

	if *company != "" {
		companyUUID, teamUUID := seedCompanyAndTeam(companyDSN, userUUID, *company, *team, *tz)
		fmt.Printf("✓ company seeded: %s (%s)\n", *company, companyUUID)
		fmt.Printf("✓ team seeded:    %s (%s)\n", *team, teamUUID)
	}

	fmt.Println()
	fmt.Println("Login at http://localhost:8080/login/ with:")
	fmt.Printf("  email:    %s\n", *email)
	fmt.Printf("  password: %s\n", *password)
}

// ----- Acme Diner demo fixture --------------------------------------------------

type demoEmployee struct {
	Name    string
	Email   string
	Phone   string
	Teams   []string // team names this employee works on
}

type demoTeam struct {
	Name   string
	Color  string
	Jobs   []demoJob
}

type demoJob struct {
	Name  string
	Color string
}

var demoTeams = []demoTeam{
	{Name: "Front of House", Color: "48B7AB", Jobs: []demoJob{
		{Name: "Server", Color: "FF8C42"},
		{Name: "Host", Color: "4A90E2"},
	}},
	{Name: "Back of House", Color: "E94B3C", Jobs: []demoJob{
		{Name: "Line Cook", Color: "C0392B"},
		{Name: "Dishwasher", Color: "27AE60"},
	}},
	{Name: "Bar", Color: "9B59B6", Jobs: []demoJob{
		{Name: "Bartender", Color: "8E44AD"},
	}},
}

var demoWorkers = []demoEmployee{
	{Name: "Alice Chen", Email: "alice@acme.local", Phone: "+15551234001", Teams: []string{"Front of House"}},
	{Name: "Bob Diaz", Email: "bob@acme.local", Phone: "+15551234002", Teams: []string{"Front of House", "Bar"}},
	{Name: "Carol Patel", Email: "carol@acme.local", Phone: "+15551234003", Teams: []string{"Front of House"}},
	{Name: "Dave Kim", Email: "dave@acme.local", Phone: "+15551234004", Teams: []string{"Back of House"}},
	{Name: "Eve Mendoza", Email: "eve@acme.local", Phone: "+15551234005", Teams: []string{"Back of House"}},
	{Name: "Frank Owens", Email: "frank@acme.local", Phone: "+15551234006", Teams: []string{"Bar"}},
}

func seedDemo(accountDSN, companyDSN, tz string) {
	rand.Seed(42) // deterministic so reseeds produce the same fixture

	accountDB := mustOpen(accountDSN)
	defer accountDB.Close()
	companyDB := mustOpen(companyDSN)
	defer companyDB.Close()

	companyName := "Acme Diner"
	managerEmail := "manager@acme.local"

	// Wipe the prior demo fixture so re-running this is idempotent.
	wipeDemo(accountDB, companyDB, companyName, managerEmail)

	// Manager
	managerUUID := seedAccountTx(accountDB, managerEmail, "acmedev", "Acme Manager", true)
	fmt.Printf("✓ manager       %s  uuid=%s\n", managerEmail, managerUUID)

	// Company
	companyUUID := mustNewUUIDStr()
	must(companyDB.Exec(`
INSERT INTO company (uuid, name, archived, default_timezone, default_day_week_starts)
VALUES (?, ?, 0, ?, ?)`, companyUUID, companyName, tz, defaultDayStart))
	must(companyDB.Exec(`
INSERT INTO admin (company_uuid, user_uuid) VALUES (?, ?)`, companyUUID, managerUUID))
	must(companyDB.Exec(`
INSERT INTO directory (company_uuid, user_uuid, internal_id) VALUES (?, ?, '')`, companyUUID, managerUUID))
	fmt.Printf("✓ company       %s  uuid=%s\n", companyName, companyUUID)

	// Teams + jobs
	teamUUIDs := make(map[string]string, len(demoTeams))
	type jobRef struct {
		teamUUID string
		jobUUID  string
	}
	jobsByTeam := make(map[string][]jobRef, len(demoTeams))
	for _, t := range demoTeams {
		teamUUID := mustNewUUIDStr()
		teamUUIDs[t.Name] = teamUUID
		must(companyDB.Exec(`
INSERT INTO team (uuid, company_uuid, name, archived, timezone, day_week_starts, color)
VALUES (?, ?, ?, 0, ?, ?, ?)`, teamUUID, companyUUID, t.Name, tz, defaultDayStart, t.Color))
		fmt.Printf("✓ team          %-16s uuid=%s\n", t.Name, teamUUID)

		for _, j := range t.Jobs {
			jobUUID := mustNewUUIDStr()
			must(companyDB.Exec(`
INSERT INTO job (uuid, team_uuid, name, archived, color)
VALUES (?, ?, ?, 0, ?)`, jobUUID, teamUUID, j.Name, j.Color))
			jobsByTeam[t.Name] = append(jobsByTeam[t.Name], jobRef{teamUUID: teamUUID, jobUUID: jobUUID})
			fmt.Printf("  job           %-16s uuid=%s\n", j.Name, jobUUID)
		}
	}

	// Workers — accounts + directory entries + worker links
	type workerRef struct {
		uuid  string
		teams []string
	}
	workers := make([]workerRef, 0, len(demoWorkers))
	for _, w := range demoWorkers {
		workerUUID := seedAccountTx(accountDB, w.Email, "worker", w.Name, false)
		must(accountDB.Exec(`UPDATE account SET phonenumber=? WHERE uuid=?`, w.Phone, workerUUID))
		must(companyDB.Exec(`
INSERT INTO directory (company_uuid, user_uuid, internal_id) VALUES (?, ?, '')`, companyUUID, workerUUID))
		for _, teamName := range w.Teams {
			teamUUID := teamUUIDs[teamName]
			must(companyDB.Exec(`
INSERT INTO worker (team_uuid, user_uuid) VALUES (?, ?)`, teamUUID, workerUUID))
		}
		fmt.Printf("✓ employee      %-15s %s  teams=%v\n", w.Name, w.Email, w.Teams)
		workers = append(workers, workerRef{uuid: workerUUID, teams: w.Teams})
	}

	// Shifts — 14 days, ~30 shifts. Mix of:
	//   - assigned vs unassigned
	//   - published (for current week) vs draft (for next week)
	loc, err := time.LoadLocation(tz)
	if err != nil {
		fatal("timezone: " + err.Error())
	}
	now := time.Now().In(loc)
	weekStart := startOfWeek(now)
	shiftCount := 0

	shiftPatterns := []struct {
		startHour, durationHours int
		teamName, jobName        string
	}{
		// Front of House — daytime + evening
		{startHour: 7, durationHours: 6, teamName: "Front of House", jobName: "Server"},
		{startHour: 11, durationHours: 6, teamName: "Front of House", jobName: "Server"},
		{startHour: 17, durationHours: 7, teamName: "Front of House", jobName: "Server"},
		{startHour: 7, durationHours: 8, teamName: "Front of House", jobName: "Host"},
		// Back of House
		{startHour: 6, durationHours: 7, teamName: "Back of House", jobName: "Line Cook"},
		{startHour: 12, durationHours: 8, teamName: "Back of House", jobName: "Line Cook"},
		{startHour: 16, durationHours: 6, teamName: "Back of House", jobName: "Dishwasher"},
		// Bar — evenings only
		{startHour: 17, durationHours: 7, teamName: "Bar", jobName: "Bartender"},
	}

	for dayOffset := 0; dayOffset < 14; dayOffset++ {
		dayStart := weekStart.AddDate(0, 0, dayOffset)
		// First week (Mon-Sun) = published. Second week = draft.
		published := dayOffset < 7
		// Saturdays + Sundays add an extra brunch shift; weekdays skip the early line cook on Mondays.
		patterns := shiftPatterns
		if dayStart.Weekday() == time.Sunday || dayStart.Weekday() == time.Saturday {
			patterns = append(patterns, struct {
				startHour, durationHours int
				teamName, jobName        string
			}{startHour: 9, durationHours: 5, teamName: "Front of House", jobName: "Server"})
		}

		for _, p := range patterns {
			start := time.Date(dayStart.Year(), dayStart.Month(), dayStart.Day(),
				p.startHour, 0, 0, 0, loc).UTC()
			stop := start.Add(time.Duration(p.durationHours) * time.Hour)
			jobs := jobsByTeam[p.teamName]
			if len(jobs) == 0 {
				continue
			}

			// Pick the matching job by name (we know the demo data).
			var jobRef jobRef
			for _, j := range jobs {
				// look up name by re-querying — small, only runs at seed time
				var jobName string
				_ = companyDB.QueryRow(`SELECT name FROM job WHERE uuid=?`, j.jobUUID).Scan(&jobName)
				if jobName == p.jobName {
					jobRef = j
					break
				}
			}
			if jobRef.jobUUID == "" {
				continue
			}

			// Assign 70% of shifts to a random worker on that team; leave 30% open.
			userUUID := ""
			if rand.Float32() < 0.7 {
				eligible := []string{}
				for _, w := range workers {
					for _, t := range w.teams {
						if t == p.teamName {
							eligible = append(eligible, w.uuid)
							break
						}
					}
				}
				if len(eligible) > 0 {
					userUUID = eligible[rand.Intn(len(eligible))]
				}
			}

			pubInt := 0
			if published {
				pubInt = 1
			}
			must(companyDB.Exec(`
INSERT INTO shift (uuid, team_uuid, job_uuid, user_uuid, published, start, stop)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
				mustNewUUIDStr(), jobRef.teamUUID, jobRef.jobUUID, userUUID,
				pubInt, start.Format("2006-01-02 15:04:05"), stop.Format("2006-01-02 15:04:05")))
			shiftCount++
		}
	}
	fmt.Printf("✓ shifts        %d total (published next 7d, draft after)\n", shiftCount)

	fmt.Println()
	fmt.Println("================================")
	fmt.Println("ACME DINER demo fixture ready.")
	fmt.Println("================================")
	fmt.Println()
	fmt.Println("Manager login: http://localhost:8080/login/")
	fmt.Printf("  email:    %s\n", managerEmail)
	fmt.Println("  password: acmedev")
	fmt.Println()
	fmt.Println("Workers (login currently not available — Staffjoy v2 had no worker login;")
	fmt.Println("they appear in the directory + on shifts):")
	for _, w := range demoWorkers {
		fmt.Printf("  %-15s %s\n", w.Name, w.Email)
	}
}

func wipeDemo(accountDB, companyDB *sql.DB, companyName, managerEmail string) {
	// Delete in dependency order. The company DB has no FKs (raw schema)
	// so we wipe children first to keep state consistent.
	emails := []string{managerEmail}
	for _, w := range demoWorkers {
		emails = append(emails, w.Email)
	}

	// Find the existing company UUID before we wipe it.
	var oldCompanyUUID sql.NullString
	_ = companyDB.QueryRow(`SELECT uuid FROM company WHERE name=?`, companyName).Scan(&oldCompanyUUID)

	// Find old user UUIDs (so we can clean their team/admin/directory rows).
	oldUserUUIDs := []string{}
	for _, e := range emails {
		var u sql.NullString
		_ = accountDB.QueryRow(`SELECT uuid FROM account WHERE email=?`, e).Scan(&u)
		if u.Valid {
			oldUserUUIDs = append(oldUserUUIDs, u.String)
		}
	}

	if oldCompanyUUID.Valid {
		// Find team UUIDs for the company so we can wipe their jobs/shifts.
		teamRows, err := companyDB.Query(`SELECT uuid FROM team WHERE company_uuid=?`, oldCompanyUUID.String)
		if err == nil {
			var teamUUIDs []string
			for teamRows.Next() {
				var u string
				if err := teamRows.Scan(&u); err == nil {
					teamUUIDs = append(teamUUIDs, u)
				}
			}
			teamRows.Close()
			for _, t := range teamUUIDs {
				_, _ = companyDB.Exec(`DELETE FROM shift WHERE team_uuid=?`, t)
				_, _ = companyDB.Exec(`DELETE FROM job WHERE team_uuid=?`, t)
				_, _ = companyDB.Exec(`DELETE FROM worker WHERE team_uuid=?`, t)
				_, _ = companyDB.Exec(`DELETE FROM manager WHERE team_uuid=?`, t)
			}
			_, _ = companyDB.Exec(`DELETE FROM team WHERE company_uuid=?`, oldCompanyUUID.String)
		}
		_, _ = companyDB.Exec(`DELETE FROM admin WHERE company_uuid=?`, oldCompanyUUID.String)
		_, _ = companyDB.Exec(`DELETE FROM directory WHERE company_uuid=?`, oldCompanyUUID.String)
		_, _ = companyDB.Exec(`DELETE FROM company WHERE uuid=?`, oldCompanyUUID.String)
	}

	// Belt-and-suspenders cleanup: anything pointing to the old user UUIDs.
	for _, u := range oldUserUUIDs {
		_, _ = companyDB.Exec(`DELETE FROM admin WHERE user_uuid=?`, u)
		_, _ = companyDB.Exec(`DELETE FROM directory WHERE user_uuid=?`, u)
		_, _ = companyDB.Exec(`DELETE FROM worker WHERE user_uuid=?`, u)
		_, _ = companyDB.Exec(`DELETE FROM manager WHERE user_uuid=?`, u)
		_, _ = companyDB.Exec(`UPDATE shift SET user_uuid='' WHERE user_uuid=?`, u)
	}

	// Account rows.
	for _, e := range emails {
		_, _ = accountDB.Exec(`DELETE FROM account WHERE email=?`, e)
	}
}

// ----- shared helpers -----------------------------------------------------------

func startOfWeek(t time.Time) time.Time {
	d := t
	for d.Weekday() != time.Monday {
		d = d.AddDate(0, 0, -1)
	}
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func mustNewUUIDStr() string {
	u, err := crypto.NewUUID()
	if err != nil {
		fatal("uuid: " + err.Error())
	}
	return u.String()
}

func seedAccount(dsn, email, password, name string, support bool) string {
	db := mustOpen(dsn)
	defer db.Close()
	return seedAccountTx(db, email, password, name, support)
}

func seedAccountTx(db *sql.DB, email, password, name string, support bool) string {
	uuid := mustNewUUIDStr()
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
		uuid, email, name, time.Now().UTC().Format("2006-01-02 15:04:05"),
		hash, salt, supportInt))
	return uuid
}

func seedCompanyAndTeam(dsn, userUUID, companyName, teamName, tz string) (string, string) {
	db := mustOpen(dsn)
	defer db.Close()

	companyUUID := mustNewUUIDStr()
	teamUUID := mustNewUUIDStr()

	must(db.Exec("DELETE FROM company WHERE name=?", companyName))
	must(db.Exec("DELETE FROM admin WHERE user_uuid=?", userUUID))
	must(db.Exec("DELETE FROM directory WHERE user_uuid=?", userUUID))

	must(db.Exec(`
INSERT INTO company (uuid, name, archived, default_timezone, default_day_week_starts)
VALUES (?, ?, 0, ?, ?)`, companyUUID, companyName, tz, defaultDayStart))
	must(db.Exec(`
INSERT INTO team (uuid, company_uuid, name, archived, timezone, day_week_starts, color)
VALUES (?, ?, ?, 0, ?, ?, '48B7AB')`, teamUUID, companyUUID, teamName, tz, defaultDayStart))
	must(db.Exec(`
INSERT INTO admin (company_uuid, user_uuid) VALUES (?, ?)`, companyUUID, userUUID))
	must(db.Exec(`
INSERT INTO directory (company_uuid, user_uuid, internal_id) VALUES (?, ?, '')`, companyUUID, userUUID))

	return companyUUID, teamUUID
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
