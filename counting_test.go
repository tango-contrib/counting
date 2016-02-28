package counting

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/lunny/log"
	"github.com/lunny/tango"
)

var (
	responseText    = "counting"
	postRequestText = "a=1&b=2"
)

type CountingAction struct {
	tango.Ctx
}

func (a *CountingAction) Get() string {
	return responseText
}

func (a *CountingAction) Post() string {
	body, _ := a.Body()
	fmt.Println(body)
	return responseText
}

func TestCounting(t *testing.T) {
	afterCounting := func(req *http.Request, read, write int) {
		if req.Method == "POST" {
			if read != len(postRequestText) {
				t.Error("expect read", len(postRequestText), "bytes, but read", read, "bytes.")
				return
			}
		}

		if write != len(responseText) {
			t.Error("expect write", len(responseText), "bytes, but write", write, "bytes.")
			return
		}
		log.Info("Read from request", read, "bytes, write to response", write, "bytes.")
	}

	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := tango.New()
	tg.Use(New(Options{AfterCounting: afterCounting}))
	tg.Use(tango.ClassicHandlers...)
	tg.Any("/", new(CountingAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "counting")

	req, err = http.NewRequest("POST", "http://localhost:8000/", strings.NewReader(postRequestText))
	if err != nil {
		t.Error(err)
	}

	buff.Reset()
	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "counting")
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
