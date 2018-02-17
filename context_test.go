package gomint

import (
	"bytes"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	_ = NewContext(w, req, logger)
}

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	data := H{
		"Hello": "World",
		"ID":    1234,
	}

	ctx := NewContext(w, req, logger)
	err := ctx.JSON(data)

	assert.Nil(t, err)
	assert.Equal(t, "application/json; charset=UTF-8", w.Header().Get("Content-Type"))
	assert.Equal(t, `{"Hello":"World","ID":1234}`+"\n", w.Body.String())
}

func TestHTML(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	ctx := NewContext(w, req, logger)
	err := ctx.HTML("hogehoge fugafuga")

	assert.Nil(t, err)
	assert.Equal(t, "text/html; charset=UTF-8", w.Header().Get("Content-Type"))
	assert.Equal(t, `hogehoge fugafuga`+"\n", w.Body.String())
}
