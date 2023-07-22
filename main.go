package main

import (
	"flag"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	fFilePath      = flag.String("files", "/data/files", "Path to static files")
	fCacheDuration = flag.Duration("cache", 12*time.Hour, "Cache duration for static files")
	fAddr          = flag.String("addr", ":3000", "Address to listen on")
)

func main() {
	flag.Parse()

	if *fFilePath == "" {
		log.Printf("No files path specified")
	}

	if *fAddr == "" {
		log.Printf("No address specified")
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
		CacheDuration: *fCacheDuration,
		MaxAge:        int((*fCacheDuration).Seconds()),
	})

	if err := app.Listen(*fAddr); err != nil {
		log.Fatal(err)
	}
}
