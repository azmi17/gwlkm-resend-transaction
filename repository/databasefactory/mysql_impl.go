package databasefactory

import (
	"database/sql"
	"fmt"
	"gwlkm-resend-transaction/repository/databasefactory/drivers"
	"os"
	"time"

	"github.com/randyardiansyah25/libpkg/util/env"

	_ "github.com/go-sql-driver/mysql"
	aes "github.com/randyardiansyah25/libpkg/security/aes"
)

type mysqlImpl struct {
	conn   *sql.DB
	prefix string
}

func newMysqlImpl() Database {
	return &mysqlImpl{}
}

func (m *mysqlImpl) Connect() error {

	//Mysql Connection String Format : [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	//Example : root:123456@tcp(127.0.0.1:3306)/employees?charset=utf8
	usr := os.Getenv(m.prefix + "mysql.username")
	usr, _ = aes.Decrypt([]byte("ECHRESENDTXT00LS"), []byte("ECHRESENDTXT00LS"), usr)

	pwd := os.Getenv(m.prefix + "mysql.password")
	pwd, _ = aes.Decrypt([]byte("ECHRESENDTXT00LS"), []byte("ECHRESENDTXT00LS"), pwd)

	addr := os.Getenv(m.prefix + "mysql.address")
	addr, _ = aes.Decrypt([]byte("ECHRESENDTXT00LS"), []byte("ECHRESENDTXT00LS"), addr)

	port := os.Getenv(m.prefix + "mysql.port")
	port, _ = aes.Decrypt([]byte("ECHRESENDTXT00LS"), []byte("ECHRESENDTXT00LS"), port)

	dbName := os.Getenv(m.prefix + "mysql.name")
	dbName, _ = aes.Decrypt([]byte("ECHRESENDTXT00LS"), []byte("ECHRESENDTXT00LS"), dbName)

	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		usr,
		pwd,
		addr,
		port,
		dbName,
	)
	maxPool := env.GetInt(m.prefix + "mysql.maxpoolsize")
	maxIdleConn := env.GetInt(m.prefix + "mysql.maxidleconn")
	maxLifeTime := env.GetInt(m.prefix + "mysql.maxconnlifetime")
	var err error
	if m.conn, err = sql.Open(drivers.MYSQL, connectionString); err != nil {
		return err
	}

	m.conn.SetMaxOpenConns(maxPool)
	m.conn.SetMaxIdleConns(maxIdleConn)
	m.conn.SetConnMaxLifetime(time.Minute * time.Duration(maxLifeTime))

	return nil
}

func (m *mysqlImpl) Ping() error {
	return m.conn.Ping()
}

func (m *mysqlImpl) GetConnection() interface{} {
	return m.conn
}

func (m *mysqlImpl) GetDriverName() string {
	return drivers.MYSQL
}

func (m *mysqlImpl) SetEnvironmentVariablePrefix(prefix string) {
	m.prefix = prefix
}

func (m *mysqlImpl) Close() {
	_ = m.conn.Close()
}
