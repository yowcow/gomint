package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yowcow/gomint/mint"
)

func main() {
	var dir string
	var port int
	var help bool
	flag.StringVar(&dir, "static-dir", "/", "path to static dir")
	flag.IntVar(&port, "port", 5000, "port to bind")
	flag.BoolVar(&help, "help", false, "help")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	app := mint.New(dir, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.Dispatcher(mint.HandlerFunc(func(ctx mint.Context) error {
		return ctx.HTML("This is /")
	})))
	mux.HandleFunc("/hello/", app.Dispatcher(mint.HandlerFunc(func(ctx mint.Context) error {
		return ctx.HTML("Hello world")
	})))
	mux.HandleFunc("/redirect", app.Dispatcher(mint.HandlerFunc(func(ctx mint.Context) error {
		ctx.Redirect("/hello/")
		return nil
	})))
	mux.HandleFunc("/foo/bar/", app.Static())

	server := http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%d", port),
	}
	logger.Println("Starting server listening on address:", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		logger.Fatalln(err)
	}
}
