package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var db *sql.DB
var server = "sonichood.database.windows.net"
var port = 1433
var user = os.Getenv("DBNAME")
var password = os.Getenv("DBWORD")
var database = "dbo.sonichood"

func isPrime(value int) bool {
	for i := 2; i <= value/2; i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func main() {

	db := &Database{
		Name:     "SonicHood",
		Server:   "sonichood.database.windows.net",
		Port:     1433,
		User:     os.Getenv("DBNAME"),
		Password: os.Getenv("DBWORD"),
		Database: "dbo.sonichood"}

	err := db.getBands()
	if err != nil {

	}

	m := martini.Classic()
	m.Use(render.Renderer(
		render.Options{
			Directory: "templates",
		},
	))

	m.Get("/", func(r render.Render, req *http.Request) {
		if req.URL.Query().Get("wait") != "" {
			sleep, _ := strconv.Atoi(req.URL.Query().Get("wait"))
			log.Printf("Sleep for %d seconds\n", sleep)
			time.Sleep(time.Duration(sleep) * time.Second)
		}
		if req.URL.Query().Get("prime") != "" {
			val, _ := strconv.Atoi(req.URL.Query().Get("prime"))
			log.Printf("Is %d prime: %t", val, isPrime(val))
		}
		r.HTML(200, "index", nil)
	})

	// db = &Database{
	// 	Name:     "SonicHood",
	// 	Server:   "sonichood.database.windows.net",
	// 	Port:     1433,
	// 	User:     os.Getenv("DBNAME"),
	// 	Password: os.Getenv("DBWORD"),
	// 	Database: "dbo.sonichood"}

	if os.Getenv("PANIC") == "true" {
		panic("this is crashing")
	}

	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	go http.Serve(listener, m)
	log.Println("Listening on 0.0.0.0:" + port)

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM)
	<-sigs
	fmt.Println("SIGTERM, time to shutdown")
	listener.Close()
}
