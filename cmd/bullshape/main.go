package main

import (
	"bullshape/db"
	"bullshape/router"
	"bullshape/utils/kafkala"
	"bullshape/utils/logger"
)

func main() {
	// log := l.NewLogger("Initialise")
	// log.Printf("Lets start")

	// Logger
	//lg := logger.NewZapSugaredLogger(os.Stdout)
	lg := logger.GetLogger()
	lg.SetExtra("Req:", "Main things")
	lg.Slogger.Info("Initialise")

	//DB
	db := db.NewDatabaseConnection()
	// Kafka
	productsProducer := kafkala.NewCompanyProducer()

	server := router.NewServer(lg.Slogger, db, productsProducer)

	server.Run()
}
