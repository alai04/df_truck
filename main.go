package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

const defaultPort = "8080"

var Db *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbHost := os.Getenv("DBHOST")
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := os.Getenv("DBNAME")
	initDB(fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName))

	r := gin.Default()
	r.GET("/query", queryApproved)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	r.Run(":" + port)
}

func initDB(dsn string) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	Db = db
	log.Printf("Init DB [%s] success", dsn)
}

func queryApproved(c *gin.Context) {
	plate := c.Query("p")
	c.JSON(http.StatusOK, gin.H{
		"result": Approved(plate),
	})
}

func Approved(plate string) bool {
	stmt, err := Db.Prepare(
		`select count(1) 
		from t_traffic_car C inner join t_traffic_main M on C.main_id = M.id 
		where M.status_code = 3 and C.plate_no = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(plate).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}
