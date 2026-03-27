package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	database "go-script-moph-appoint/src/database"
	"go-script-moph-appoint/src/loadenv"
	"go-script-moph-appoint/src/moph"
	"go-script-moph-appoint/src/schedule"
)

func main() {
	db := database.ConnectDB()
	defer db.Close()
	log.Println("server started")

	url, ck, sk := loadenv.LoadMOPH()
	moph.Init(url, ck, sk)

	schedule.Start(db, loadenv.LoadScheduleTime())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("server stopped")
}
