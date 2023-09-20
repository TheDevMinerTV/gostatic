package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	fFilePath      = flag.String("files", "/data/files", "Path to static files")
	fCacheDuration = flag.Duration("cache", 12*time.Hour, "Cache duration for static files")
	fAddr          = flag.String("addr", ":3000", "Address to listen on")
	fCompressLevel = flag.Int("compress-level", 1, "Compression level for static files. 0 = disabled, 1 = default, 2 = best")
	fLogRequests   = flag.Bool("log-requests", false, "Log requests to stdout")
	fSPA           = flag.Bool("spa", false, "Serve index.html for 404 pages (for SPA apps)")
	fIndex         = flag.String("index", "index.html", "Index file relative from the files path")
)

func main() {
	flag.Parse()

	if *fFilePath == "" {
		log.Printf("No files path specified")
	}

	if *fAddr == "" {
		log.Printf("No address specified")
	}

	if *fCompressLevel < 0 || *fCompressLevel > 2 {
		log.Printf("Invalid compression level")
	}

	app := fiber.New(fiber.Config{
		GETOnly:           true,
		EnablePrintRoutes: true,
	})

	if *fLogRequests {
		log.Printf("Enabling request logging")

		app.Use(logger.New())
	}

	if *fCompressLevel > 0 {
		log.Printf("Enabling compression: %d", *fCompressLevel)

		app.Use(compress.New(compress.Config{
			Level: compress.Level(*fCompressLevel),
		}))
	}

	if *fSPA {
		log.Printf("Serving files from %s as SPA", *fFilePath)
	} else {
		log.Printf("Serving files from %s", *fFilePath)
	}

	app.Static("/", *fFilePath, fiber.Static{
		Compress:      true,
		Browse:        false,
		Download:      false,
		CacheDuration: *fCacheDuration,
		Index:         *fIndex,
		MaxAge:        int((*fCacheDuration).Seconds()),
	})

	if *fSPA {
		app.Use(func(c *fiber.Ctx) error {
			return c.SendFile(filepath.Join(*fFilePath, *fIndex))
		})
	}

	log.Printf("Listening on %s", *fAddr)
	if err := app.Listen(*fAddr); err != nil {
		log.Fatal(err)
	}
}
