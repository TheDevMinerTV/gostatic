package main

import (
	"flag"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	fFilePath      = flag.String("files", "/data/files", "Path to static files")
	fCacheDuration = flag.Duration("cache", 12*time.Hour, "Cache duration for static files")
	fAddr          = flag.String("addr", ":3000", "Address to listen on")
	fCompressLevel = flag.Int("compress-level", 1, "Compression level for static files. Setting this to -1 will entirely remove the middleware.")
	fLogRequests   = flag.Bool("log-requests", false, "Log requests to stdout")
)

func main() {
	flag.Parse()

	if *fFilePath == "" {
		log.Printf("No files path specified")
	}

	if *fAddr == "" {
		log.Printf("No address specified")
	}

	if *fCompressLevel < -2 || *fCompressLevel > 2 {
		log.Printf("Invalid compression level")
	}

	app := fiber.New(fiber.Config{
		GETOnly:           true,
		EnablePrintRoutes: true,
	})

	if *fLogRequests {
		app.Use(logger.New())
	}

	if *fCompressLevel >= 0 {
		app.Use(compress.New(compress.Config{
			Level: compress.Level(*fCompressLevel),
		}))
	}

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
