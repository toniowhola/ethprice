package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

var config Config

// Config - config yaml
type Config struct {
	Db struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"db"`
}

// Price - for price we get from the db
type Price struct {
	ID        int     `json:"id"`
	Coin      string  `json:"coin"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
}

func getPrices(w http.ResponseWriter, r *http.Request) {
	db := dbConnection()
	defer db.Close()
	results, err := db.Query(`
    SELECT id, coin, price, created_at
    FROM prices
    ORDER BY created_at DESC
    LIMIT 10
  `)
	if err != nil {
		fmt.Println(err)
	}
	var prices []Price

	for results.Next() {
		var price Price
		err = results.Scan(&price.ID, &price.Coin, &price.Price, &price.CreatedAt)
		if err != nil {
			fmt.Println(err)
		}
		prices = append(prices, price)
	}

	json.NewEncoder(w).Encode(prices)
}

func dbConnection() *sql.DB {
	db, err := sql.Open(
		"mysql",
		config.Db.Username+":"+config.Db.Password+"@tcp(127.0.0.1:3306)/",
	)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE price_retrieval")
	if err != nil {
		fmt.Println(err)
	}

	return db
}

func readConfig() {
	file, err := os.Open("./config.yml")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		fmt.Println(err)
	}
}

func main() {
	readConfig()
	http.HandleFunc("/prices", getPrices)
	fmt.Println("Server started at :10000")
	log.Fatal(http.ListenAndServe(":10000", nil))
}
