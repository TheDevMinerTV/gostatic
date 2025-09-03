package main

import (
	"flag"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
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

type userList []string

func (u *userList) String() string {
	return ""
}

func (u *userList) Set(value string) error {
	*u = append(*u, value)
	return nil
}

var users userList

func parseUsers(userList []string) map[string]string {
	userMap := make(map[string]string)
	for _, user := range userList {
		parts := strings.SplitN(user, ":", 2)
		if len(parts) == 2 {
			username := strings.TrimSpace(parts[0])
			password := strings.TrimSpace(parts[1])
			if username != "" && password != "" {
				userMap[username] = password
			}
		}
	}
	return userMap
}

func main() {
	flag.Var(&users, "user", "User credentials in format username:password (can be used multiple times)")
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

	userMap := parseUsers(users)
	if len(userMap) > 0 {
		log.Printf("Enabling basic authentication")

		app.Use(basicauth.New(basicauth.Config{
			Users: userMap,
			Realm: "Restricted Area",
			Authorizer: func(user, pass string) bool {
				expectedPass, exists := userMap[user]
				return exists && expectedPass == pass
			},
		}))
	}

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
