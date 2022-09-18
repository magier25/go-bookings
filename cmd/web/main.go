package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/magier25/go-bookings/internal/config"
	"github.com/magier25/go-bookings/internal/driver"
	"github.com/magier25/go-bookings/internal/handlers"
	"github.com/magier25/go-bookings/internal/helpers"
	"github.com/magier25/go-bookings/internal/models"
	"github.com/magier25/go-bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	log.Println("Starting mail listener...")
	listenForMail()

	log.Printf("Running http server on port %s...\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = false
	infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=127.0.0.1 port=5432 dbname=bookings user=postgres password=postgres")
	if err != nil {
		log.Fatal("Cannot connect to database! Exiting...")
	}
	log.Println("Connected to database.")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot read template cache")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
	render.NewRenderer(&app)

	return db, nil
}
