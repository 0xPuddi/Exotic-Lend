package main

import (
	"fmt"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Price struct {
	Price     int `json:"price"`
	Timestamp int `json:"timestamp"`
}

type Asset struct {
	Id     string `josn:"id"`
	Ticker string `json:"ticker"`
	Price  Price  `json:"Price"`
}

func main() {
	fmt.Println("Hello World")
	utils.HandleFatalError(godotenv.Load(".env"))

	// db.

	// URL := "postgresql://postgres:" + os.Getenv(DB_PASSWORD) + "@" + os.Getenv(DB_IP) + ":" + os.Getenv(DB_PORT) + "/postgres?sslmode=" + os.Getenv(DB_SSL_MODE)
	// db, err := sql.Open("postgres", URL)
	// utils.HandleFatalError(err)
	// defer db.Close()

	// stats := db.Stats()
	// fmt.Println(stats)

	// insertEntry(db)
	// addEntry(db)
}
