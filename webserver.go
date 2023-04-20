package main

import (
	"flag"
	_ "net/http/pprof"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	app.Use(compress.New())
	app.Use(recover.New())

	app.Static("/", dir, fiber.Static{
		Compress: true,
		Browse:   true,
	})
	app.Listen(addr)
}
