package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yowcow/gomint/mint"
)

type MyApp struct{}

func (a MyApp) DispatchRoot(ctx mint.Context) error {
	return ctx.HTML("This is root")
}

func (a MyApp) DispatchHello(ctx mint.Context) error {
	return ctx.HTML("Hello world")
}

func (a MyApp) DispatchRedirect(ctx mint.Context) error {
	ctx.Redirect("/hello/")
	return nil
}

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
	a := MyApp{}
	mt := mint.New(dir, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/", mt.Dispatcher(a.DispatchRoot))
	mux.HandleFunc("/hello/", mt.Dispatcher(a.DispatchHello))
	mux.HandleFunc("/redirect", mt.Dispatcher(a.DispatchRedirect))
	mux.HandleFunc("/foo/bar/", mt.Static())

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
