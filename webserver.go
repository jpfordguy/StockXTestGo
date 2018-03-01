package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	uuid "github.com/nu7hatch/gouuid"
)

const (
	DB_HOST     = "localhost"
	DB_PORT     = 5432
	DB_DBTYPE   = "postgres"
	DB_USER     = "postgres"
	DB_PASSWORD = "password"
	DB_NAME     = "stockxtest"
)

var mainDb *sql.DB

func CreateGuid() string {
	id, err := uuid.NewV4()
	if err == nil {
		return fmt.Sprintf("%s", id)
	}
	return ""
}

func DoPanicError(funcname string, err error) {
	text := fmt.Sprintf("Error occured in %s - %s", funcname, err.Error())
	panic(text)
}

func DoPanicString(funcname string, err string) {
	text := fmt.Sprintf("Error occured in %s - %s", funcname, err)
	panic(text)
}

func CreateName(name string, db *sql.DB) {
	id := CreateGuid()
	qry := fmt.Sprintf("INSERT INTO ShoeNames (Name, Id) VALUES ('%s', '%s')", name, id)
	res, err := db.Exec(qry)
	if err != nil {
		DoPanicError("CreateName", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		DoPanicError("CreateName", err)
	}
	if rows < 1 {
		DoPanicString("CreateName", "no rows were updated")
	}
}

func SetId(name string, id string, size int, db *sql.DB) {
	qry := fmt.Sprintf("INSERT INTO ShoeSizes(SizeId,Id,Size) VALUES ('%s', '%s', %d)", name, id, size)
	res, err := db.Exec(qry)
	if err != nil {
		DoPanicError("SetId", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		DoPanicError("SetId", err)
	}
	if rows < 1 {
		DoPanicString("SetId", "no rows were updated")
	}
}

func GetId(name string, db *sql.DB) string {
	qry := fmt.Sprintf("SELECT Id FROM ShoeNames WHERE Name='%s'", name)
	rows, err := db.Query(qry)
	if err != nil {
		return ""
	}
	defer rows.Close()
	var ret string
	if rows.Next() {
		err := rows.Scan(&ret)
		if err != nil {
			DoPanicError("GetId", err)
		}
	}
	if err == nil {
		return ret
	}
	return ""
}

func AppendSizeToDatabase(name string, size int64, db *sql.DB) {
	id := CreateGuid()
	nameid := GetId(name, db)
	if nameid == "" {
		CreateName(name, db)
	}
	qry := fmt.Sprintf("INSERT INTO ShoeSizes(SizeId,Id,Size) VALUES ('%s', '%s', %d)", id, nameid, size)
	res, err := db.Exec(qry)
	if err != nil {
		panic(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	if rows < 1 {
		panic("Unable to insert size into database")
	}
}

func GetTrueToSize(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["name"]
	//if name == nil {
	//	http.Error(w, "Invalid manufacturer", 400)
	//	return
	//}
	if name == "" {
		http.Error(w, "Invalid manufacturer", 400)
		return
	}
	id := GetId(name, mainDb)
	if id == "" {
		http.Error(w, "Invalid manufacturer", 400)
		return
	}
	qry := fmt.Sprintf("SELECT AVG(Size) FROM ShoeSizes WHERE Id = '%s'", id)
	rows, err := mainDb.Query(qry)
	if err != nil {
		http.Error(w, "Error querying database", 400)
		return
	}
	defer rows.Close()
	var ret float64
	for rows.Next() {
		err := rows.Scan(&ret)
		if err != nil {
			DoPanicError("GetTrueToSize", err)
		}
	}
	if err != nil {
		http.Error(w, "Error querying database", 400)
		return
	}
	outstr := fmt.Sprintf("%g", ret)
	http.Error(w, outstr, 200)
}

func AppendSize(w http.ResponseWriter, r *http.Request) {
	qp := r.URL.Query()
	name := qp.Get("name")
	size := qp.Get("size")

	if name == "" {
		http.Error(w, "Invalid manufacturer", 400)
		return
	}
	if size == "" {
		http.Error(w, "Invalid size", 400)
		return
	}

	var isize int64
	isize, err := strconv.ParseInt(size, 10, 32)
	if err != nil {
		http.Error(w, "Invalid size", 400)
		return
	}
	if isize < 1 || isize > 5 {
		http.Error(w, "Invalid size", 400)
		return
	}

	AppendSizeToDatabase(name, isize, mainDb)
	http.Error(w, "Value saved", 201)
}

func main() {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	mainDb = db
	if err != nil {
		DoPanicError("main", err)
	}

	router := mux.NewRouter()
	//router.Host("localhost:8000")
	router.HandleFunc("/truetosize/{name}", GetTrueToSize).Methods("GET")
	router.HandleFunc("/append", AppendSize).Methods("GET")
	http.ListenAndServe(":8000", router)
	defer mainDb.Close()
}
