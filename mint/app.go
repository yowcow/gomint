package mint

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type H map[string]interface{}

type App struct {
	StaticDir string
	Logger    *log.Logger
}

func New(dir string, logger *log.Logger) *App {
	app := &App{dir, logger}
	return app
}

func (app App) Static() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		path := filepath.Join(app.StaticDir, req.URL.Path)
		fi, err := os.Stat(path)
		if err != nil || fi.IsDir() {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "<h1>File Not Found</h1>")
			return
		}
		f, err := os.Open(path)
		if err != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "<h1>Internal Server Error</h1>")
			return
		}
		defer f.Close()
		io.Copy(w, f)
	}
}

func (app App) Dispatcher(fns ...HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := NewContext(w, req, app.Logger)
		app.Logger.Println(req.Method, req.URL)
		for _, fn := range fns {
			err := fn(*ctx)
			if err != nil {
				app.Logger.Printf("[ERROR] %s", err)
				return
			}
		}
	}
}
