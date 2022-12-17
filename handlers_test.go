package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
)

func TestLoadHandler(t *testing.T) {
	server, err := LoadServerFromSettings(`settings.json`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	if value := server.CertFolder; value != `./certs` {
		t.Fatalf(`expected './certs' but got '%v'`, value)
	}
	if value := server.ServeFolder; value != `./serveTest` {
		t.Fatalf(`expected './serveTest' but got '%v'`, value)
	}
}

func TestLoadHomepage(t *testing.T) {
	server, _ := LoadServerFromSettings(`settings.json`)
	router := server.GenerateRouter()
	w := &testWriter{}
	r := getRequestFor(`https://www.host1.com/`)
	router.ServeHTTP(w, r)

	err := checkResponse(w, 200, `./serveTest/host1/index.html`)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestLoadFile(t *testing.T) {
	server, _ := LoadServerFromSettings(`settings.json`)
	router := server.GenerateRouter()
	w := &testWriter{}
	r := getRequestFor(`https://www.host1.com/scripts.js`)
	router.ServeHTTP(w, r)

	err := checkResponse(w, 200, `./serveTest/host1/scripts.js`)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test404Response(t *testing.T) {
	server, _ := LoadServerFromSettings(`settings.json`)
	router := server.GenerateRouter()
	w := &testWriter{}
	r := getRequestFor(`https://www.host1.com/invalid_file`)
	router.ServeHTTP(w, r)

	err := checkResponse(w, 404, `./serveTest/host1/404.html`)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
func TestHostWhitelist(t *testing.T) {
	server, err := LoadServerFromSettings(`settings.json`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	whitelist := server.CollectHostWhitelist()
	expected := []string{"host1.com", "www.host1.com", "something.somewhere.com"}
	if !reflect.DeepEqual(whitelist, expected) {
		t.Fatalf(`expected %v but got %v`, expected, whitelist)
	}
}

func TestRouter(t *testing.T) {
	server, err := LoadServerFromSettings(`settings.json`)
	if err != nil {
		t.Fatalf(`expected no error but got: %v`, err.Error())
	}
	router := server.GenerateRouter()
	if err = checkRoute(router, `https://www.host1.com/`); err != nil {
		t.Fatalf(err.Error())
	}
	if err = checkRoute(router, `https://www.host1.com/index.html`); err != nil {
		t.Fatalf(err.Error())
	}
	if err = checkRoute(router, `https://host1.com/`); err != nil {
		t.Fatalf(err.Error())
	}
	if err = checkRoute(router, `https://host1.com/index.html`); err != nil {
		t.Fatalf(err.Error())
	}
	if err = checkRoute(router, `https://something.somewhere.com/`); err != nil {
		t.Fatalf(err.Error())
	}
	if err = checkRoute(router, `https://something.somewhere.com/index.html`); err != nil {
		t.Fatalf(err.Error())
	}
}

func checkRoute(router *mux.Router, url string) error {
	match := &mux.RouteMatch{}
	r := getRequestFor(url)
	if success := router.Match(r, match); !success {
		return match.MatchErr
	}
	return nil
}

func checkResponse(w *testWriter, expectedStatus int, expectedFile string) error {
	if w.status != expectedStatus {
		return fmt.Errorf(`expected status %v but got %v`, expectedStatus, w.status)
	}
	expected, _ := os.ReadFile(expectedFile)
	if !reflect.DeepEqual(w.content, expected) {
		return fmt.Errorf("expected %v content but got:\n%v", expectedFile, w.content)
	}
	return nil
}

func getRequestFor(testUrl string) *http.Request {
	return requestFor(testUrl, `GET`)
}

func requestFor(testUrl string, method string) *http.Request {
	u, _ := url.Parse(testUrl)
	return &http.Request{
		Method: method,
		URL:    u,
	}
}

type testWriter struct {
	content []byte
	status  int
	header  http.Header
}

func (w *testWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *testWriter) Write(content []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	if w.content == nil {
		w.content = make([]byte, len(content))
		copy(w.content, content)
		return len(content), nil
	}
	w.content = append(w.content, content...)
	return len(content), nil
}

func (w *testWriter) WriteHeader(status int) {
	w.status = status
}
