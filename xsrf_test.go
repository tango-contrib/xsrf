package xsrf

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"bytes"
	"compress/gzip"

	"github.com/lunny/tango"
)

type XsrfAction struct {
}

func (xsrf *XsrfAction) Do() {
}

func TestXsrf(t *testing.T) {
	go func() {
		tg := tango.Classic()
		tg.Use(NewXsrf(time.Minute * 20))
		tg.Post("/", new(XsrfAction))
		tg.Run("0.0.0.0:9996")
	}()

	resp, bs, err := post("http://localhost:9996/?id=1&name=lllll")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(bs)

	if resp.StatusCode == http.StatusOK {
		t.Error("should say xsrf error.")
		return
	}
}

type NoCheckXsrfAction struct {
}

func (NoCheckXsrfAction) CheckXsrf() bool {
	return false
}

func (NoCheckXsrfAction) Do() string {
	return "this action will not check xsrf"
}

func TestNoCheckXsrf(t *testing.T) {
	go func() {
		tg := tango.Classic()
		tg.Use(NewXsrf(time.Minute * 20))
		tg.Post("/", new(NoCheckXsrfAction))
		tg.Run("0.0.0.0:9995")
	}()

	resp, bs, err := post("http://localhost:9995/?id=1&name=lllll")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(bs)

	if resp.StatusCode != http.StatusOK {
		t.Error("should say ok.")
		return
	}
}

func gzipDecode(src []byte) ([]byte, error) {
	rd := bytes.NewReader(src)
	b, err := gzip.NewReader(rd)
	if err != nil {
		return nil, err
	}

	defer b.Close()

	data, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func get(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.Header.Get("Content-Encoding") == "gzip" {
		data, err := gzipDecode(bs)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return string(bs), nil
}

func post(url string) (*http.Response, string, error) {
	resp, err := http.Post(url, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return resp, "", err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, "", err
	}

	if resp.Header.Get("Content-Encoding") == "gzip" {
		data, err := gzipDecode(bs)
		if err != nil {
			return resp, "", err
		}
		return resp, string(data), nil
	}
	return resp, string(bs), nil
}
