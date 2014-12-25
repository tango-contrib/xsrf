package xsrf

import (
	"reflect"
	"net/http/httptest"
	"net/http"
	"testing"
	"time"
	"bytes"

	"github.com/lunny/tango"
)

type XsrfAction struct {
	Checker
}

func (xsrf *XsrfAction) Post() {
}

func TestXsrf(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := tango.Classic()
	tg.Use(New(time.Minute * 20))
	tg.Post("/", new(XsrfAction))

	req, err := http.NewRequest("POST", "http://localhost:8000/?id=1&name=lllll", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	refute(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
}

type NoCheckXsrfAction struct {
	NoCheck
}

func (NoCheckXsrfAction) Post() string {
	return "this action will not check xsrf"
}

func TestNoCheckXsrf(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := tango.Classic()
	tg.Use(New(time.Minute * 20))
	tg.Post("/", new(NoCheckXsrfAction))

	req, err := http.NewRequest("POST", "http://localhost:8000/?id=1&name=lllll", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "this action will not check xsrf")
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
