package main

import (
	"flag"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var fFilePath = flag.String("files", "/data/files", "Path to static files")

func main() {
	flag.Parse()

	if *fFilePath == "" {
		log.Printf("No files path specified")
	}

	app := fiber.New(fiber.Config{
		GETOnly:           true,
		EnablePrintRoutes: true,
	})

	app.Use(logger.New())

	app.Static("/", *fFilePath, fiber.Static{
		Compress:      true,
		Browse:        false,
		Download:      false,
		CacheDuration: 12 * time.Hour,
		// let the browser cache this item for 12 hours, after that it should refresh
		MaxAge: 43200,
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
