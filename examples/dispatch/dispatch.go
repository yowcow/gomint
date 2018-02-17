package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yowcow/gomint"
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
	app := gomint.New(dir, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.Dispatcher(gomint.HandlerFunc(func(ctx gomint.Context) error {
		return ctx.HTML("This is /")
	})))
	mux.HandleFunc("/hello/", app.Dispatcher(gomint.HandlerFunc(func(ctx gomint.Context) error {
		return ctx.HTML("Hello world")
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
