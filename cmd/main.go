package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/hoanguyen1998/crypto-payment-system/internal/drivers"
	"github.com/hoanguyen1998/crypto-payment-system/internal/handlers"
	"github.com/hoanguyen1998/crypto-payment-system/internal/repository/postgresRepo"
	"github.com/hoanguyen1998/crypto-payment-system/internal/services"
	"github.com/robfig/cron/v3"
)

func main() {
	r := mux.NewRouter()

	db := setupDB()

	repo := postgresRepo.NewPostgresRepo(db)

	services := services.NewAppService(repo)

	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	app := handlers.NewServerHandler(services, redis)

	r.HandleFunc("/register", app.Register)
	r.HandleFunc("/login", app.Login)
	r.Use(app.CheckToken)
	r.HandleFunc("/master-public-keys", app.CreateMasterPublicKey).Methods("POST")
	r.HandleFunc("/apps", app.CreateApp)
	r.HandleFunc("/app-keys", app.CreateAppKey)
	r.HandleFunc("/orders", app.CreateOrder)
	r.HandleFunc("/send-tx", app.SendTransaction)
	r.HandleFunc("/withdraw", app.GetOrdersToWithdraw)

	c := cron.New()
	c.AddFunc("* * * * *", func() { app.ProcessEthBlock() })
	c.Start()

	http.ListenAndServe(":8080", r)
}

func setupDB() *sql.DB {
	dbHost := flag.String("dbhost", "localhost", "database host")
	dbPort := flag.String("dbport", "5432", "database port")
	dbUser := flag.String("dbuser", "", "database user")
	dbPass := flag.String("dbpass", "", "database password")
	databaseName := flag.String("db", "vigilate", "database name")
	dbSsl := flag.String("dbssl", "disable", "database ssl setting")

	flag.Parse()

	if *dbUser == "" || *dbHost == "" || *dbPort == "" || *databaseName == "" {
		fmt.Println("Missing required flags.")
		os.Exit(1)
	}

	log.Println("Connecting to database....")
	dsnString := ""

	// when developing locally, we often don't have a db password
	if *dbPass == "" {
		dsnString = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			*dbHost,
			*dbPort,
			*dbUser,
			*databaseName,
			*dbSsl)
	} else {
		dsnString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			*dbHost,
			*dbPort,
			*dbUser,
			*dbPass,
			*databaseName,
			*dbSsl)
	}

	db, err := drivers.ConnectPostgres(dsnString)
	if err != nil {
		log.Fatal("Cannot connect to database!", err)
	}

	return db
}
