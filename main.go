package main

import (
	"gwlkm-resend-transaction/delivery"
	"gwlkm-resend-transaction/delivery/handler"
	"gwlkm-resend-transaction/delivery/router"
	"gwlkm-resend-transaction/helper"
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
	// Temporary Functions
	go handler.RepostingSchedulerRepoObserver()

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
		_ = glg.Log("Reloading configuration file...")
		er = godotenv.Overload(".env")
	} else {
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

	//Untuk error jika terdapat gagal reposting menggunakan flag FAIL
	repostinglogEr := glg.FileWriter("log/repostings.err", 0666)
	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.FAIL, repostinglogEr)

	_ = glg.Log("=================Service Info===================")
	_ = glg.Log("Application Name:", helper.AppName)
	_ = glg.Log("Application Version:", helper.AppVersion)
	_ = glg.Log("Last Build:", helper.LastBuild)
	_ = glg.Log("================================================")
}

func PrepareDatabase() {
	var er error

	// # App DB: Echannel V3
	databasefactory.Echannel, er = databasefactory.GetDatabase()
	databasefactory.Echannel.SetEnvironmentVariablePrefix("ech.")
	if er != nil {
		glg.Fatal(er.Error())
	}

	_ = glg.Log("Connecting to echannel v3...")
	if er = databasefactory.Echannel.Connect(); er != nil {
		_ = glg.Error("Connection to echannelv3 failed : ", er.Error())
		os.Exit(1)
	}

	if er = databasefactory.Echannel.Ping(); er != nil {
		_ = glg.Error("Cannot ping echannel v3 : ", er.Error())
		os.Exit(1)
	}

	// # App DB: Apex
	databasefactory.Apex, er = databasefactory.GetDatabase()
	databasefactory.Apex.SetEnvironmentVariablePrefix("apx.")
	if er != nil {
		glg.Fatal(er.Error())
	}

	_ = glg.Log("Connecting to apex..")
	if er = databasefactory.Apex.Connect(); er != nil {
		_ = glg.Error("Connection to apex failed: ", er.Error())
		os.Exit(1)
	}

	if er = databasefactory.Apex.Ping(); er != nil {
		_ = glg.Error("Cannot ping apex: ", er.Error())
		os.Exit(1)
	}

	// # App DB: Sys Apex
	databasefactory.SysApex, er = databasefactory.GetDatabase()
	databasefactory.SysApex.SetEnvironmentVariablePrefix("sysapx.")
	if er != nil {
		glg.Fatal(er.Error())
	}

	_ = glg.Log("Connecting to sys apex..")
	if er = databasefactory.SysApex.Connect(); er != nil {
		_ = glg.Error("Connection to sys apex failed: ", er.Error())
		os.Exit(1)
	}

	if er = databasefactory.SysApex.Ping(); er != nil {
		_ = glg.Error("Cannot ping sys apex: ", er.Error())
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
