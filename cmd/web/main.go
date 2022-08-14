package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/magier25/go-bookings/pkg/config"
	"github.com/magier25/go-bookings/pkg/handlers"
	"github.com/magier25/go-bookings/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot read template cache")
	}
	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)
	mux := routes(&app)

	fmt.Printf("Running http server on port %s...\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: mux,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
