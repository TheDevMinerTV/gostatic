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
	fDownload      = flag.Bool("download", false, "Enables direct Download for Served Files")
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

	index := filepath.Join(*fFilePath, *fIndex)

	app := fiber.New(fiber.Config{
		GETOnly:           true,
		EnablePrintRoutes: true,
	})

	if *fLogRequests {
		log.Printf("Enabling request logging")

		app.Use(logger.New())
	}

	if *fDownload {
		log.Printf("Enabling downloads")
	}

	if *fCompressLevel > 0 {
		log.Printf("Enabling compression: %d", *fCompressLevel)

		app.Use(compress.New(compress.Config{
			Level: compress.Level(*fCompressLevel),
		}))
	}

	if *fSPA {
		log.Printf("Serving files from %s as SPA with index file %s", *fFilePath, index)
	} else {
		log.Printf("Serving files from %s with index file %s", *fFilePath, index)
	}

	app.Static("/", *fFilePath, fiber.Static{
		Compress:      true,
		Download:      *fDownload,
		CacheDuration: *fCacheDuration,
		MaxAge:        int((*fCacheDuration).Seconds()),
		ByteRange:     true,
	})

	if *fSPA {
		app.Use(func(c *fiber.Ctx) error {
			log.Printf("Serving index.html for %s", c.Path())
			return c.SendFile(index)
		})
	}

	log.Printf("Listening on %s", *fAddr)
	if err := app.Listen(*fAddr); err != nil {
		log.Fatal(err)
	}
}
