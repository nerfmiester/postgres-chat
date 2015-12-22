package main
import (
	"fmt"
	"database/sql"
	"github.com/nu7hatch/gouuid"
	"net/http"
	"github.com/jackc/pgx"
	"os"
	"time"
)

const (
	DB_HOST		=	"localhost"
	DB_USER     = 	"adrianjackson"
	DB_PASSWORD = 	"0hn0sh3cr13d"
	DB_NAME     = 	"chitchat"
)

type Users struct {
	Id		int64			`id`
	Uuid 	string		`uuid`
	Name 	string			`name`
	Email 	string			`email`
	Pwd		string			`password`
	Created time.Time			`created_at`
}

/*
Id       int64                        `id`
    Created  int64                        `created`
    Updated  int64                        `modified`
    FName    string                       `firstName`
    LName    string                       `lastName`
    Comments *SomeNonPersistentStructure  `db:"-"`


*/
var db *sql.DB
var pool *pgx.ConnPool
var conn *pgx.Conn

func init() {

	var err error
	conn, err = pgx.Connect(extractConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
}


func main() {

	http.HandleFunc("/users", usersCreate)
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
	if name == "" || email == "" || password == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	u4, err := uuid.NewV4()
	u4string := u4.String()

	usersX := Users{
		Uuid: u4string,
		Name: name,
		Email: email,
		Pwd: password,
		Created: time.Now(),
	}

   err = addTask(usersX)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to add task: %v\n", err)
		os.Exit(1)
	}

}

func addTask(user Users) error {
	_, err := conn.Exec("insert into users(uuid,name,email,password,created_at) values($1,$2,$3,$4,$5)", user.Uuid,user.Name,user.Email,user.Pwd,user.Created)
	return err
}

func extractConfig() pgx.ConnConfig {
	var config pgx.ConnConfig

	config.Host = DB_HOST
	if config.Host == "" {
		config.Host = "localhost"
	}

	config.User =DB_USER
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