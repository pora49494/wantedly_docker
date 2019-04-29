package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	user     = "docker"
	dbname   = "docker"
	password = "docker"
)

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type Users []User

func errorHandler(w http.ResponseWriter, status_code int) {
	message := make(map[string]int)
	message["status_code"] = status_code
	w.WriteHeader(status_code)
	json.NewEncoder(w).Encode(message)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	message := make(map[string]string)
	message["message"] = "Hello World!!"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/users", getUsersEndpoint).Methods("GET")
	myRouter.HandleFunc("/users/{id}", getUserEndpoint).Methods("GET")
	myRouter.HandleFunc("/users", createUserEndpoint).Methods("POST")
	myRouter.HandleFunc("/users/{id}", updateUserEndpoint).Methods("PUT")
	myRouter.HandleFunc("/users/{id}", deleteUserEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func createUserEndpoint(w http.ResponseWriter, r *http.Request) {
	var user User

	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &user)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError)
		return
	}

	sqlStatement := `
	INSERT INTO users (name, email)
	VALUES ($1, $2)
	RETURNING id, name, email, created_at, updated_at;`

	err = db.QueryRow(sqlStatement, user.Name, user.Email).Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &user.Updated_at)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			errorHandler(w, http.StatusBadRequest)
		} else {
			errorHandler(w, http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

	// fmt.Println("Hit createUserEndpoint()")
}

func updateUserEndpoint(w http.ResponseWriter, r *http.Request) {
	var user User
	params := mux.Vars(r)
	id := params["id"]

	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &user)

	sqlStatement := `
	UPDATE users
	SET name = $2, email = $3
	WHERE id = $1
	RETURNING id, name, email, created_at, updated_at;`
	err = db.QueryRow(sqlStatement, id, user.Name, user.Email).Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &user.Updated_at)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			errorHandler(w, http.StatusBadRequest)
		} else {
			errorHandler(w, http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

	// fmt.Println("Hit updateUserEndpoint()")
}

func getUserEndpoint(w http.ResponseWriter, r *http.Request) {
	var user User
	params := mux.Vars(r)
	id := params["id"]

	sqlStatement := `
	SELECT id, name, email, created_at, updated_at FROM users
	WHERE id = $1;`

	err := db.QueryRow(sqlStatement, id).Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &user.Updated_at)
	if err != nil {
		errorHandler(w, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

	// fmt.Println("Hit getUserEndpoint()")
}

func deleteUserEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	sqlStatement := `
	DELETE FROM users
	WHERE id = $1;`

	result, err := db.Exec(sqlStatement, id)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError)
		return
	}
	row_num, err := result.RowsAffected()
	if row_num == 0 {
		errorHandler(w, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)

	// fmt.Println("Hit deleteUserEndpoint()")
}

func getUsersEndpoint(w http.ResponseWriter, r *http.Request) {
	users := Users{}

	sqlStatement := `
	SELECT id, name, email, created_at, updated_at FROM users`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Updated_at, &user.Created_at)
		if err != nil {
			errorHandler(w, http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

	fmt.Println("Hit getUsersEndpoint()")
}

func databaseSetUp() {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)
	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}

var db *sql.DB

func main() {
	databaseSetUp()
	defer db.Close()
	handleRequest()
}
