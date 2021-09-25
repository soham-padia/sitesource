package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"

	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/solow-crypt/bookings/internal/config"
	"github.com/solow-crypt/bookings/internal/driver"
	"github.com/solow-crypt/bookings/internal/handlers"
	"github.com/solow-crypt/bookings/internal/helpers"
	"github.com/solow-crypt/bookings/internal/models"
	"github.com/solow-crypt/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)
	fmt.Println("Starting mail listener...")
	ListenForMail()

	fmt.Println("starting at port", portNumber)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	gob.Register(models.Registration{})
	gob.Register(models.User{})

	inProgress := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "postgres", "Database name")
	dbUser := flag.String("dbuser", "postgres", "Database user")
	dbPass := flag.String("dbpass", "", "Database pass")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database SSl settings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("missing required flags")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	//change this to true when in  production
	app.InProduction = *inProgress
	app.UseCache = *useCache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.Infolog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction //change it when going for release

	app.Session = session

	//connext to db
	log.Println("connecting to Database...")

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)

	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("cannot connect to db , dying")
	}
	log.Println("Connected to Database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Println(err)
		log.Fatal("cannot create template")
		return nil, err
	}

	app.TemplateCache = tc

	//for developers use false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
