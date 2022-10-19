package main

import (
	"gwlkm-resend-transaction/delivery"
	"gwlkm-resend-transaction/delivery/router"
	"gwlkm-resend-transaction/repository/databasefactory"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kpango/glg"
)

func main() {
	go delivery.PrintoutObserver()
	router.Start()
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())

	LoadConfiguration(false)
	if os.Getenv("app.database_driver") != "" {
		PrepareDatabase()
	}

	go ReloadObserver()
}

func LoadConfiguration(isReload bool) {
	var er error
	if isReload {
		_ = glg.Log("e-Channel Recycle Transaction Service")
		_ = glg.Log("Reloading configuration file...")
		er = godotenv.Overload(".env")
	} else {
		_ = glg.Log("e-Channel Recycle Transaction Service")
		_ = glg.Log("Loading configuration file...")
		er = godotenv.Load(".env")
	}

	if er != nil {
		_ = glg.Error("Configuration file not found...")
		os.Exit(1)
	}

	//Opsi agar log utk level LOG, DEBUG, INFO dicatat atau tidak
	//Jika menggunakan docker atau dibuatkan service, log sudah dibuatkan, sehingga direkomendasikan
	//app log di set false
	appLog := os.Getenv("app.log")
	if appLog == "true" {
		log := glg.FileWriter("log/application.log", 0666)
		glg.Get().
			SetMode(glg.BOTH).
			AddLevelWriter(glg.LOG, log).
			AddLevelWriter(glg.DEBG, log).
			AddLevelWriter(glg.INFO, log).
			AddLevelWriter(glg.ERR, log)
	}

	//Untuk error, akan selalu dicatat dalam file
	logEr := glg.FileWriter("log/application.err", 0666)
	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, logEr).
		AddLevelWriter(glg.WARN, logEr)
}

func PrepareDatabase() {
	var er error

	// # DB 1
	databasefactory.AppDb1, er = databasefactory.GetDatabase()
	databasefactory.AppDb1.SetEnvironmentVariablePrefix("ech.")
	if er != nil {
		glg.Fatal(er.Error())
	}

	_ = glg.Log("Connecting to echannelv3...")
	if er = databasefactory.AppDb1.Connect(); er != nil {
		_ = glg.Error("Connection to echannelv3 failed : ", er.Error())
		os.Exit(1)
	}

	if er = databasefactory.AppDb1.Ping(); er != nil {
		_ = glg.Error("Cannot ping echannelv3 : ", er.Error())
		os.Exit(1)
	}

	// # DB 2
	databasefactory.AppDb2, er = databasefactory.GetDatabase()
	databasefactory.AppDb2.SetEnvironmentVariablePrefix("apx.")
	if er != nil {
		glg.Fatal(er.Error())
	}

	_ = glg.Log("Connecting to apex..")
	if er = databasefactory.AppDb2.Connect(); er != nil {
		_ = glg.Error("Connection to apex failed: ", er.Error())
		os.Exit(1)
	}

	if er = databasefactory.AppDb2.Ping(); er != nil {
		_ = glg.Error("Cannot ping apex: ", er.Error())
		os.Exit(1)
	}

	_ = glg.Log("Database Connected")
	_ = glg.Log("Service Started")
}

func ReloadObserver() {
	sign := make(chan os.Signal, 1)     // bikin channel yang isinya dari signal
	signal.Notify(sign, syscall.SIGHUP) // kalo ada signal HUP simpan ke channel sign

	func() {
		for {
			<-sign
			LoadConfiguration(true)
		}
	}()
}
