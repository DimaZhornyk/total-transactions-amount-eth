package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"transaction-amount-eth/server/db"
	"transaction-amount-eth/server/gracefully"
	"transaction-amount-eth/server/handlers"
	totalTransactionAmount "transaction-amount-eth/server/handlers/TotalTransactionAmount"
	"transaction-amount-eth/server/models"
)

func main() {
	// .env file with token was left in the project by purpose to simplify the review
	_ = godotenv.Load(".env")
	log.Println(os.Getenv("ETHERSCAN_TOKEN"))
	log.SetOutput(os.Stdout)
	var config models.Config
	models.ReadConfig(&config)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.User, config.PostgreSQL.Password, config.PostgreSQL.Dbname)
	log.Println(psqlInfo)
	conn := db.Connect(psqlInfo)
	db.CreateTables(conn)

	h := handlers.BaseHandler{
		DbManager:        &db.Manager{DB: conn},
		EtherScanManager: models.EtherScanManager{Url: "https://api.etherscan.io/api", Module: "proxy", ApiKey: os.Getenv("ETHERSCAN_TOKEN")},
	}

	router := mux.NewRouter()
	ttaHandler := totalTransactionAmount.NewHandler(&h)
	router.HandleFunc("/api/block/{block_number:[0-9]+}/total", ttaHandler.TotalTransactionAmount).Methods("GET")
	server := http.Server{
		Addr:    config.Server.Port,
		Handler: router,
	}

	log.Printf("Starting golang server on port %s\n", config.Server.Port)
	if err := gracefully.Serve(server.ListenAndServe, func(ctx context.Context) error {
		if err := server.Shutdown(ctx); err != nil {
			return err
		}
		conn.Close()
		return nil
	}); err != nil {
		log.Fatalln("ERR:", err)
	}
}
