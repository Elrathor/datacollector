package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	Timestamp time.Time
	Key       string
	Value     string
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS data (timestamp DATETIME, key TEXT, value TEXT)")
	if err != nil {
		log.Fatal(err)
	}
}

func saveData(data Data) {
	_, err := db.Exec("INSERT INTO data (timestamp, key, value) VALUES (?, ?, ?)", data.Timestamp, data.Key, data.Value)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(c *gin.Context) {
	for key, values := range c.Request.URL.Query() {
		for _, value := range values {
			saveData(Data{
				Timestamp: time.Now(),
				Key:       key,
				Value:     value,
			})
		}
	}
}

func main() {
	r := gin.Default()

	// Get username and password from environment variables
	username := os.Getenv("BASIC_AUTH_USERNAME")
	password := os.Getenv("BASIC_AUTH_PASSWORD")

	r.Use(gin.BasicAuth(gin.Accounts{
		username: password,
	}))

	// Create a rate limiter
	limiter := tollbooth.NewLimiter(1, nil) // 1 request per second

	// Add the rate limiter as a middleware
	r.Use(tollbooth_gin.LimitHandler(limiter))

	r.GET("/", handler)
	log.Fatal(r.Run(":80"))
}
