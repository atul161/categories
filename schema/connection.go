package connection

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Connection interface {
	MakeConnection() *sql.DB
}

//Info contains the details of new connection
type Info struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

//New Connection will open the connection with the database information
// that is passed as an argument.
func NewConnection(info *Info) Connection {
	conn := Info{
		Host:     info.Host,
		Port:     info.Port,
		User:     info.User,
		Password: info.Password,
		Dbname:   info.Dbname,
	}
	return &conn
}

//Make Connection will open the connection with the DB
// and return the instance of that connection.
//Note: Caller must have to close the connection after the use
func (info *Info) MakeConnection() *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		info.Host, info.Port, info.User, info.Password, info.Dbname)

	DB, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		panic(err)
	}
	return DB
}
