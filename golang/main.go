package main

import (
	"database/sql"
  "github.com/joho/godotenv"
  "github.com/gorilla/mux"
  "encoding/json"
	"fmt"
	"log"
  "os"
  "net/http"

	_ "github.com/lib/pq"
)

type User struct {
	ID       int
	Email    string
	Password string
}

func dbConnection() string {
  // TODO fill this in directly or through environment variable
  // Build a DSN e.g. postgres://username:password@url.com:5432/dbName
  //set connection string:
  err := godotenv.Load(".env")
  if err != nil {
    log.Fatalf("Error loading .env file")
  }
  user := os.Getenv("PSQL_USER")
  pass := os.Getenv("PSQL_PASS")
  host := os.Getenv("PSQL_HOST")
  db := os.Getenv("PSQL_DB")

  psqlConn := fmt.Sprintf(
    "postgres://%v:%v@%v:5432/%v?sslmode=disable",
    user,
    pass,
    host,
    db,
  )

  return psqlConn
}

func homeLink(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Welcome to our homepage")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
  DB_DSN := dbConnection()
	// Create DB pool
	db, err := sql.Open("postgres", DB_DSN)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()
	// Create an empty user and make the sql query (using $1 for the parameter)
	var myUser User
	userSql := "SELECT id, email, password FROM users WHERE id = $1"

  err = db.QueryRow(userSql, 1).Scan(&myUser.ID, &myUser.Email, &myUser.Password)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(myUser)
  //fmt.Fprint(w, "Hi " + myUser.Email + " welcome back!\n", )
}

func main() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", homeLink)
  router.HandleFunc("/users", getUsers).Methods("GET")
  log.Fatal(http.ListenAndServe(":8090", router))
}
