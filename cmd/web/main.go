package main

import (
	"encoding/gob"
	"flag"
	"fmt"
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

	// read flags
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	flag.Parse()
	if *dbName == "" || *dbUser == "" {
		log.Println("Missing required flags")
		os.Exit(1)
	}

	// connect to database
	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database! Exiting...")
	}
	defer db.SQL.Close()

	log.Println("Connected to database.")

	err = run(*inProduction, *useCache, db)
	if err != nil {
		log.Fatal(err)
	}
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

func run(inProduction, useCache bool, db *driver.DB) error {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = inProduction
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

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot read template cache")
		return err
	}
	app.TemplateCache = tc
	app.UseCache = useCache

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)
	render.NewRenderer(&app)

	return nil
}
