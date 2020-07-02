package main

import (
	"database/sql"
  "github.com/joho/godotenv"
  "github.com/gorilla/mux"
  "github.com/gorilla/sessions"
  "encoding/json"
	"fmt"
	"log"
  "os"
  "net/http"

	_ "github.com/lib/pq"
)

var (
  // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
  key = []byte("super-secret-key")
  store = sessions.NewCookieStore(key)
)

type userInfo struct {
	ID       int
	Email    string
	Password string
}

type users struct {
  Users []userInfo
}

func login(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, "cookie-name")
  // auth here
  session.Values["authenticated"] = true
  session.Save(r, w)
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

func queryUsers(userData *users) error {
  DB_DSN := dbConnection()
	// Create DB pool
	db, err := sql.Open("postgres", DB_DSN)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}

  rows, err := db.Query(`SELECT id, email, password FROM users`)
  if err != nil {
    return err
  }
  defer db.Close()

  for rows.Next() {
    user := userInfo{}
    err := rows.Scan(
      &user.ID,
      &user.Email,
      &user.Password,
    )
    if err != nil {
      return err
    }
    userData.Users = append(userData.Users, user)
  }
  err = rows.Err()
  if err != nil {
    return err
  }
  return nil

}

func getUsers(w http.ResponseWriter, req *http.Request) {

  session, _ := store.Get(req, "cookie-name")

  if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
    http.Error(w, "Forbidden", http.StatusForbidden)
    return
  }

  userData := users{}

  err := queryUsers(&userData)
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  out, err := json.Marshal(userData)
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  fmt.Fprintf(w, string(out))
  //w.Header().Set("Content-Type", "application/json")
  //json.NewEncoder(w).Encode(string(out))
}

func main() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", homeLink)
  router.HandleFunc("/users", getUsers).Methods("GET")
  router.HandleFunc("/login", login)
  log.Fatal(http.ListenAndServe(":8090", router))
}
