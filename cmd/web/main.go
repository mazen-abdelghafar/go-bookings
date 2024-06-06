package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"github.com/kidkever/go-bookings/internal/config"
	"github.com/kidkever/go-bookings/internal/driver"
	"github.com/kidkever/go-bookings/internal/handlers"
	"github.com/kidkever/go-bookings/internal/helpers"
	"github.com/kidkever/go-bookings/internal/models"
	"github.com/kidkever/go-bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	log.Println("Starting mail listener...")
	listenForMail()
	defer close(app.MailChan)

	msg := fmt.Sprintf("Staring application on port %s", portNumber)
	fmt.Println(msg)

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

	// what am i going to store in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// read flags
	// inProduction := flag.Bool("production", true, "Application is in production")
	// useCache := flag.Bool("cache", true, "Use template cache")
	// dbHost := flag.String("dbhost", "localhost", "Database host")
	// dbName := flag.String("dbname", "", "Database name")
	// dbUser := flag.String("dbuser", "", "Database user")
	// dbPassword := flag.String("dbpassword", "", "Database password")
	// dbPort := flag.String("dbport", "5432", "Database port")
	// dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")
	// flag.Parse()
	// if *dbName == "" || *dbUser == "" {
	// 	fmt.Println("Missing required flags")
	// 	os.Exit(1)
	// }

	// using env variables instead of cmd flags
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	inProduction, err := strconv.ParseBool(os.Getenv("PRODUCTION"))
	if err != nil {
		log.Fatal("Can't parse PRODUCTION env")
	}
	useCache, err := strconv.ParseBool(os.Getenv("CACHE"))
	if err != nil {
		log.Fatal("Can't parse CACHE env")
	}

	dbHost := os.Getenv("DBHOST")
	dbName := os.Getenv("DBNAME")
	dbUser := os.Getenv("DBUSER")
	dbPassword := os.Getenv("DBPASSWORD")
	dbPort := os.Getenv("DBPORT")
	dbSSL := os.Getenv("DBSSL")

	if dbName == "" || dbUser == "" {
		fmt.Println("Missing required env variables")
		os.Exit(1)
	}

	// creating a mail channel
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = inProduction
	app.UseCache = useCache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPassword, dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database...")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
