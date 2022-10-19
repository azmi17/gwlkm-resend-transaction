package databasefactory

var AppDb1, AppDb2 Database

type Database interface {
	Connect() error
	Ping() error
	GetConnection() interface{}
	GetDriverName() string
	SetEnvironmentVariablePrefix(string)
	Close()
}
