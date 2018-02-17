package gomint

import (
	"log"
	"net/http"
)

type H map[string]interface{}

type App struct {
	Logger *log.Logger
}

func New(logger *log.Logger) *App {
	app := &App{logger}
	return app
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
