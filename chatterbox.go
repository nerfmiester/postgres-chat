package main
import (
	"fmt"
	"database/sql"
	"github.com/nu7hatch/gouuid"
	"net/http"
	"github.com/jackc/pgx"
	"os"
	"time"
	"gopkg.in/gorp.v1"
	_ "github.com/lib/pq"
)

const (
	DB_HOST = "localhost"
	DB_USER = "adrianjackson"
	DB_PASSWORD = "0hn0sh3cr13d"
	DB_NAME = "chitchat"
)

type Users struct {
	Id      int64            `db:"id"`
	Uuid    string           `db:"uuid"`
	Name    string           `db:"name"`
	Email   string           `db:"email"`
	Pwd     string           `db:"password"`
	Created time.Time        `db:"created_at"`
	Updated time.Time        `db:"updated_at"`
}

/*
Id       int64                        `id`
    Created  int64                        `created`
    Updated  int64                        `modified`
    FName    string                       `firstName`
    LName    string                       `lastName`
    Comments *SomeNonPersistentStructure  `db:"-"`


*/

var conn *pgx.Conn


type GorpController struct {
	Txn *gorp.Transaction
}


var dbmap gorp.DbMap

func init() {

	var err error
	conn, err = pgx.Connect(extractConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
}

func main() {

	var err error
	dbinfo := fmt.Sprintf("user=%s password=%s host=%s  dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_HOST, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	dialect := gorp.PostgresDialect{}

	// construct a gorp DbMap

	dbmap = gorp.DbMap{Db: db, Dialect: dialect}
	dbmap.AddTable(Users{}).SetKeys(true, "Id")


	http.HandleFunc("/users/create", usersCreate)
	http.HandleFunc("/users/amend", usersAmend)
	http.ListenAndServe(":3000", nil)

}

func usersCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")


	u4, err := uuid.NewV4()
	u4string := u4.String()

	usersX := Users{
		Uuid: u4string,
		Name: name,
		Email: email,
		Pwd: password,
		Created: time.Now(),
		Updated: time.Now(),
	}

	err = addTask(usersX)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to add task: %v\n", err)
		os.Exit(1)
	}

}

func usersAmend(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	if r.FormValue("uuid") == "" {
		fmt.Printf("No key to update with, process will exit . . . . ")
		os.Exit(1)
	}

	Uuid := r.FormValue("uuid")

	fmt.Printf("r.FormValue(uuid) . . . %s\n", r.FormValue("uuid"))
	fmt.Printf("Uuid . . . %s\n", Uuid)

	var userZ Users
	err := dbmap.SelectOne(&userZ, "select * from users where uuid=$1", Uuid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get record %s: %v\n", Uuid, err)
		os.Exit(1)
	}

	fmt.Printf("\nUser Name . . . %s\n", userZ.Name)
	fmt.Printf("User ID. . . %s\n", userZ.Id)

	var name string
	var email string
	var password string

	name = r.FormValue("name")
	email = r.FormValue("email")
	password = r.FormValue("password")
	fmt.Printf("User Name 1. . . %s\n", name)
	if r.FormValue("name") != "" {
		fmt.Printf("User Name 1.5 . . %s\n", userZ.Name)
		userZ.Name = name
		fmt.Printf("User Name 1.6 . . %s\n", userZ.Name)
	}
	if r.FormValue("email") != "" {
		userZ.Email= email
	}
	if r.FormValue("password") != "" {
		userZ.Pwd = password
	}

	fmt.Printf("r.FormValue(name) . . . %s\n", r.FormValue("name"))
	fmt.Printf("User Name 2. . . %s\n", name)

	usersX := Users{
		Id: userZ.Id,
		Uuid: userZ.Uuid,
		Name: userZ.Name,
		Email: userZ.Email,
		Pwd: userZ.Pwd,
		Created: userZ.Created,
		Updated: time.Now(),
	}

	// count is the # of rows updated, which should be 1 in this example
	count, err := dbmap.Update(&usersX)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to update record: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("The number of records updated is . . . . %d\n", count)

}

func addTask(user Users) error {
	_, err := conn.Exec("insert into users(uuid,name,email,password,created_at,updated_at) values($1,$2,$3,$4,$5,$6)", user.Uuid, user.Name, user.Email, user.Pwd, user.Created,user.Updated)
	return err
}

func extractConfig() pgx.ConnConfig {
	var config pgx.ConnConfig

	config.Host = DB_HOST
	if config.Host == "" {
		config.Host = "localhost"
	}

	config.User = DB_USER
	if config.User == "" {
		config.User = "postgres"
	}

	config.Password = DB_PASSWORD

	config.Database = DB_NAME
	if config.Database == "" {
		config.Database = "todo"
	}

	return config
}