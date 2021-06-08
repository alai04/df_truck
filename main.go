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

type approvedRec struct {
	plate string
	bDate string
	eDate string
	bTime string
	eTime string
}

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
	result := Approved(plate)
	c.JSON(http.StatusOK, gin.H{
		"result": len(result) > 0,
		"num":    len(result),
		"desc":   Rec2String(result),
	})
}

func Approved(plate string) []approvedRec {
	stmt, err := Db.Prepare(
		`select C.plate_no, D.begin_date, D.end_date, T.begin_time, T.end_time
		from t_traffic_car C inner join t_traffic_main M on C.main_id = M.id 
		inner join t_traffic_date D on D.main_id = M.id 
		inner join t_traffic_time T on T.main_id = M.id
		where M.status_code = 3 and D.data_type = 1 and T.data_type = 1 and C.plate_no = ?
		order by D.begin_date`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var recs []approvedRec
	rows, err := stmt.Query(plate)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var rec approvedRec
		if err = rows.Scan(&rec.plate, &rec.bDate, &rec.eDate, &rec.bTime, &rec.eTime); err != nil {
			log.Fatal(err)
		}
		recs = append(recs, rec)
	}
	return recs
}

func Rec2String(recs []approvedRec) string {
	if len(recs) == 0 {
		return "Empty record"
	}
	str := "审批通行时间：\n"
	for _, rec := range recs {
		str += fmt.Sprintf("%s - %s, %s - %s\n", rec.bDate, rec.eDate, rec.bTime, rec.eTime)
	}
	return str
}
