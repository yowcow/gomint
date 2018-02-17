package gomint

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	_ = New(logger)
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
			app := New(logger)

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
