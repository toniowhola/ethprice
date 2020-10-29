package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

var config Config

// EthPrice - contains the price from GeckoCoin
type EthPrice struct {
	Ethereum struct {
		Usd float64 `json:"usd"`
	} `json:"ethereum"`
}

// Config - config yaml
type Config struct {
	Db struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"db"`
}

func readPrice(db *sql.DB) {
	coinName := "ethereum"
	rootURL := "https://api.coingecko.com/api/v3"

	params := url.Values{}
	params.Add("ids", coinName)
	params.Add("vs_currencies", "USD")

	requestURL := rootURL + "/simple/price?" + params.Encode()
	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var ethPrice EthPrice
	json.Unmarshal(bodyBytes, &ethPrice)
	fmt.Println(ethPrice.Ethereum.Usd)

	insert, err := db.Query(
		"INSERT INTO prices (coin, price, created_at) VALUES (?,?,?)",
		coinName,
		ethPrice.Ethereum.Usd,
		time.Now().Local(),
	)
	if err != nil {
		fmt.Println(err)
	}
	defer insert.Close()
}

func dbConnection() *sql.DB {
	db, err := sql.Open(
		"mysql",
		config.Db.Username+":"+config.Db.Password+"@tcp(127.0.0.1:3306)/",
	)
	if err != nil {
		panic(err)
	}
	return db
}

func setupDB(db *sql.DB) {

	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS price_retrieval")
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("USE price_retrieval")
	if err != nil {
		fmt.Println(err)
	}

	stmt, err := db.Prepare(
		`
		CREATE TABLE IF NOT EXISTS prices(
			id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			coin VARCHAR(10),
			price FLOAT,
			created_at DATETIME
		)
		`,
	)
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
	}
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
	db := dbConnection()
	defer db.Close()
	setupDB(db)
	for {
		readPrice(db)
		time.Sleep(5 * time.Second)
	}

}
