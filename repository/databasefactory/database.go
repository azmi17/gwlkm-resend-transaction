package databasefactory

var Echannel, Apex, SysApex Database

type Database interface {
	Connect() error
	Ping() error
	GetConnection() interface{}
	GetDriverName() string
	SetEnvironmentVariablePrefix(string)
	Close()
}
