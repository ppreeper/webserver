package main

import (
	"flag"
	_ "net/http/pprof"

	"github.com/gofiber/compression"
	"github.com/gofiber/fiber"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"
)

func main() {
	var addr string
	var dir string
	var log bool
	// var bindval string
	flag.StringVar(&addr, "addr", ":8080", "addr to start")
	flag.StringVar(&dir, "d", ".", "directory to serve")
	flag.BoolVar(&log, "l", false, "enable logging")
	flag.Parse()

	app := fiber.New()
	if log {
		app.Use(logger.New())
	}
	app.Use(compression.New())
	app.Use(helmet.New())
	// Optional
	cfg := recover.Config{
		Handler: func(c *fiber.Ctx, err error) {
			c.SendString(err.Error())
			c.SendStatus(500)
		},
	}

	app.Use(recover.New(cfg))

	app.Static("/", dir)
	app.Listen(addr)
}

// MaxAge sets expire headers based on extension
// func MaxAge(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var age time.Duration
// 		ext := filepath.Ext(r.URL.String())

// 		// Timings are based on github.com/h5bp/server-configs-nginx

// 		switch ext {
// 		case ".rss", ".atom":
// 			age = time.Hour / time.Second
// 		case ".css", ".js":
// 			age = (time.Hour * 24 * 365) / time.Second
// 		case ".jpg", ".jpeg", ".gif", ".png", ".ico", ".cur", ".gz", ".svg", ".svgz", ".mp4", ".ogg", ".ogv", ".webm", ".htc":
// 			age = (time.Hour * 24 * 30) / time.Second
// 		default:
// 			age = 0
// 		}

// 		if age > 0 {
// 			w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", age))
// 		}

// 		h.ServeHTTP(w, r)
// 	})
// }
