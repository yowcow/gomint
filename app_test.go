package gomint

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	_ = New("/", logger)
}

func TestStatic(t *testing.T) {
	f, err := os.Open("static/foo/bar/test.jpg")
	if err != nil {
		t.Fatal("expected no error but got", err)
	}
	defer f.Close()
	jpgbody, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal("expected no error but got", err)
	}

	type Case struct {
		Subtest             string
		Path                string
		ExpectedStatus      int
		ExpectedContentType string
		ExpectedBody        []byte
	}
	cases := []Case{
		{
			"non existing file",
			"/foo/bar.html",
			404,
			"text/html; charset=utf-8",
			[]byte("<h1>File Not Found</h1>\n"),
		},
		{
			"existing HTML file",
			"/foo/bar/index.html",
			200,
			"text/html; charset=utf-8",
			[]byte("<h1>Hello world</h1>\n"),
		},
		{
			"existing text file",
			"/foo/bar/index.txt",
			200,
			"text/plain; charset=utf-8",
			[]byte("Hello world\n"),
		},
		{
			"existing jpg file",
			"/foo/bar/test.jpg",
			200,
			"image/jpeg",
			jpgbody,
		},
		{
			"existing dir",
			"/foo/bar/",
			404,
			"text/html; charset=utf-8",
			[]byte("<h1>File Not Found</h1>\n"),
		},
	}

	for _, c := range cases {
		t.Run(c.Subtest, func(t *testing.T) {
			logbuf := new(bytes.Buffer)
			logger := log.New(logbuf, "", 0)
			app := New("static", logger)

			handler := app.Static()
			svr := httptest.NewServer(handler)
			defer svr.Close()

			req, err := http.NewRequest("GET", svr.URL+c.Path, nil)
			if err != nil {
				t.Fatal("expected no error but got ", err)
			}
			client := http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal("expected no error but got ", err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal("expected no error but got ", err)
			}

			assert.Equal(t, c.ExpectedContentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, c.ExpectedStatus, resp.StatusCode)
			assert.Equal(t, c.ExpectedBody, body)
		})
	}
}

func TestDispatch(t *testing.T) {
	type Case struct {
		Subtest        string
		Funcs          []HandlerFunc
		ExpectedStatus int
		ExpectedBody   string
	}
	cases := []Case{
		{
			"simple dispatch",
			[]HandlerFunc{
				func(ctx Context) error {
					return ctx.JSON(H{"hoge": "fuga"})
				},
			},
			200,
			`{"hoge":"fuga"}` + "\n",
		},
		{
			"cascading dispatch",
			[]HandlerFunc{
				func(ctx Context) error {
					ctx.Stash["count"] = 0
					return nil
				},
				func(ctx Context) error {
					ctx.Stash["count"] = ctx.Stash["count"].(int) + 1
					return nil
				},
				func(ctx Context) error {
					ctx.Stash["count"] = ctx.Stash["count"].(int) + 2
					return nil
				},
				func(ctx Context) error {
					return ctx.JSON(H{
						"total": ctx.Stash["count"],
					})
				},
			},
			200,
			`{"total":3}` + "\n",
		},
		{
			"error while cascading dispatch",
			[]HandlerFunc{
				func(ctx Context) error {
					ctx.ResponseWriter().WriteHeader(http.StatusForbidden)
					ctx.HTML("forbidden!!")
					return errors.New("not authorized")
				},
				func(ctx Context) error {
					return ctx.HTML("hoge")
				},
			},
			403,
			"forbidden!!\n",
		},
	}

	for _, c := range cases {
		t.Run(c.Subtest, func(t *testing.T) {
			logbuf := new(bytes.Buffer)
			logger := log.New(logbuf, "", 0)
			app := New("/", logger)

			handler := app.Dispatcher(c.Funcs...)
			svr := httptest.NewServer(handler)
			defer svr.Close()

			req, err := http.NewRequest("GET", svr.URL, nil)
			if err != nil {
				t.Fatal("expected no error but got ", err)
			}
			client := http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal("expected no error but got ", err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal("expected no error but got ", err)
			}

			assert.Equal(t, c.ExpectedStatus, resp.StatusCode)
			assert.Equal(t, c.ExpectedBody, string(body))
		})
	}
}
