package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"subscribe/data"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore" //redis base session store for scs
	"github.com/alexedwards/scs/v2"         // http session management
	"github.com/gomodule/redigo/redis"      //client for redis data base
	_ "github.com/jackc/pgconn"             // postgresDB Driver and ToolKit
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func main() {
	//create a database -- postgres
	db := initDB()
	db.Ping()

	// create session -- redis
	session := initSession()

	// create logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	//create channels

	//create waitgrrop
	wg := sync.WaitGroup{}

	//set up the application config

	// type Config struct {  config.go
	// 		Session  *scs.SessionManager
	// 		DB       *sql.DB
	// 		Wait     *sync.WaitGroup
	// 		InfoLog  *log.Logger
	// 		ErrorLog *log.Logger
	//
	// }
	app := Config{
		Session:  session,
		DB:       db,
		Wait:     &wg,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Models:   data.New(db),
	}
	// set up mail
	app.Mailer = app.createMail()
	go app.listenForMail()

	// listen for signals
	go app.listenForShutdown() // running in the background

	//listen for web connection
	app.serve()
}

func (app *Config) serve() {
	//start http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	app.InfoLog.Println("Starting Web Server")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0) // gracefully shutdown
}

func (app *Config) shutdown() { // gracefully shutdown
	// perform any cleanup tasks
	app.InfoLog.Println("run cleanup tasks")

	// block until waitgroup is empty
	app.Wait.Wait()
	app.Mailer.DoneChan <- true

	app.InfoLog.Println("closing Channels and shutting down application")
	close(app.Mailer.MailerChan)
	close(app.Mailer.ErrorChan)
	close(app.Mailer.DoneChan)
}

func initSession() *scs.SessionManager {
	//store information in the session
	gob.Register(data.User{})

	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
}

func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}
	return redisPool
}

func (app *Config) createMail() Mail {
	errorChan := make(chan error)
	mailerChan := make(chan Message, 100)
	mailerDoneChan := make(chan bool)

	m := Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025,
		Encryption:  "none",
		FromName:    "Info",
		FromAddress: "info@mycompany.com",
		Wait:        app.Wait,
		ErrorChan:   errorChan,
		MailerChan:  mailerChan,
		DoneChan:    mailerDoneChan,
	}
	return m
}
