package mint

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Stash map[string]interface{}

type Context struct {
	w   http.ResponseWriter
	req *http.Request
	Stash
	*log.Logger
}

func NewContext(w http.ResponseWriter, req *http.Request, logger *log.Logger) *Context {
	return &Context{w, req, Stash{}, logger}
}

func (ctx Context) ResponseWriter() http.ResponseWriter {
	return ctx.w
}

func (ctx Context) Request() *http.Request {
	return ctx.req
}

func (ctx Context) JSON(data interface{}) error {
	w := ctx.w
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}

func (ctx Context) HTML(data string) error {
	w := ctx.w
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := fmt.Fprintln(w, data)
	return err
}

func (ctx Context) Redirect(location string) {
	w := ctx.w
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusFound)
}

func (ctx Context) RedirectPermanently(location string) {
	w := ctx.w
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusMovedPermanently)
}
