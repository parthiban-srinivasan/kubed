package county

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	_ "os"
	"strings"
	_ "time"
)

var db *sql.DB
var err error

type county struct {
	name  string
	state string
}

// Init module...
func main() {

	//      databaseConfig := os.Getenv("MYSQL_CONNECTION")

	//      db, err = sql.Open("mysql", databaseConfig)
	//      db, err = sql.Open("mysql", "root:root@tcp(104.196.22.179:3306)/testdb")

	//  Open validates the database arguments without creating connections
	db, err = sql.Open("mysql", "root@cloudsql(mygo-1217:us-central1:locdb)/testdb")

	if err != nil {
		log.Printf("Database not found, not good you know")
		log.Fatal("database not found")
	}

	//  Root request is handled here
	http.HandleFunc("/", rootHandler)

	//  Health check is handled by "healthz" handler
	http.HandleFunc("/healthz", healthyHandler)

	//  Create table via createhandler
	http.HandleFunc("/create", createHandler)

	//  Warmup of instance (code load during instance creation) is handled by here
	http.HandleFunc("/warmup", warmupHandler)

	log.Print("Listening on port 8080")
  log.Fatal(http.ListenAndServe(":8080", nil))

}

// end of Init function

// Root request will be handled.
func rootHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		fmt.Fprint(w, "only GET method allow \n")
		http.NotFound(w, r)
		return
	}

	if r.URL.Path != "/county" {
		fmt.Fprint(w, "only / PATH allow \n")
		http.NotFound(w, r)
		return
	}

	countyName := r.URL.RawQuery
	count := strings.Count(countyName, "")

	if count == 1 {
		fmt.Fprint(w, "Query parm absent\n")
		return
	}

	county, err := queryByCountyName(countyName)
	//rA := r.RemoteAddr

	switch {
	case err == sql.ErrNoRows:
		fmt.Fprint(w, "No county with that name\n")
	case err != nil:
		//log.Fatal(err)
		fmt.Fprint(w, "Query by county failed %v \n", err)
	default:
		fmt.Fprint(w, "Welcome to County Service %s %s", county.name, county.state)
	}

}

// Warmup request will be handled here.
func healthyHandler(w http.ResponseWriter, r *http.Request) {

	//  Test database connection
	err = db.Ping()

	if err != nil {
		//log.Fatal(err)
		fmt.Fprint(w, "Database failure %s", err)
	} else {
		fmt.Fprint(w, "Healthy - county service")
	}
}

// Warmup request will be handled here.
func createHandler(w http.ResponseWriter, r *http.Request) {

	create_stmt := `CREATE TABLE IF NOT EXISTS county (
				name       VARCHAR(15),
				state      CHAR(2)
			)`
	_, err := db.Exec(create_stmt)

	if err != nil {
		fmt.Fprint(w, "Failed - County table not created %s", err)
	} else {
		fmt.Fprint(w, "County table created")
	}
}

// Warmup request will be handled here.
func warmupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Service warmed up \n")
}

func queryByCountyName(n string) (county, error) {

	var c county

	err := db.QueryRow("SELECT name, state FROM county WHERE name=?", n).Scan(&c.name, &c.state)

	return c, err
}
