package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/ysfada/product-management-system/database"
	"github.com/ysfada/product-management-system/database/migrations"
	_ "github.com/ysfada/product-management-system/docs"
	"github.com/ysfada/product-management-system/handlers"
)

const (
	defaultHost      = "localhost"
	defaultPort      = "8080"
	defaultSourceURL = "file://database/migrations"
)

// @title Inventory Management
// @version 1.0
// @description This is an API for inventory management

// @contact.name Yusuf Ada
// @contact.url https://github.com/ysfada
// @contact.email yusufadaa@gmail.com

// @license.name MIT
// @license.url https://spdx.org/licenses/MIT.html

// @BasePath /api/v1
func main() {
	// debug := flag.Bool("debug", false, "run in debug mode")
	flag.Parse()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("HOST")
	if len(host) == 0 {
		host = defaultHost
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if len(databaseURL) == 0 {
		log.Fatal("Error 'DATABASE_URL' environment variable not found")
	}
	sourceURL := os.Getenv("MIGRATION_SOURCE_URL")
	if len(sourceURL) == 0 {
		sourceURL = defaultSourceURL
	}

	runMigrations, err := strconv.ParseBool(os.Getenv("RUN_MIGRATIONS"))
	if err != nil {
		runMigrations = true
	}

	database.CreateConnection(databaseURL)
	defer database.DbConn.Close()

	if runMigrations {
		migrations.Run(sourceURL, databaseURL)
	}

	fiberConf := fiber.Config{
		Prefork: false,
	}

	// if *debug {
	// 	fiberConf.Prefork = false
	// }

	app := fiber.New(fiberConf)
	app.Use(recover.New())

	app.Get("/docs/*", swagger.Handler)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	handlers.Use(v1)

	app.Static("/public", "./public", fiber.Static{
		Compress: true,
	})

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", host, port)))
}
