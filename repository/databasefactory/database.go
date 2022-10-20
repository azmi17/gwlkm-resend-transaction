package databasefactory

var Echannel, Apex Database

type Database interface {
	Connect() error
	Ping() error
	GetConnection() interface{}
	GetDriverName() string
	SetEnvironmentVariablePrefix(string)
	Close()
}
