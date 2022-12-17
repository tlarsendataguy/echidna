package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type HostInfo struct {
	Folder        string
	HostWhitelist []string
}

type Server struct {
	CertFolder  string
	ServeFolder string
	Hosts       []HostInfo
}

func LoadServerFromSettings(settingsFilePath string) (*Server, error) {
	settings := &Server{}

	content, err := os.ReadFile(settingsFilePath)
	if err != nil {
		log.Printf(`error reading settings file: %v`, err.Error())
		return nil, err
	}
	err = json.Unmarshal(content, settings)
	if err != nil {
		log.Printf(`error parsing settings file: %v`, err.Error())
		return nil, err
	}
	return settings, err
}

func (s *Server) homepageHandler(hostFolder string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fullPath := path.Join(s.ServeFolder, hostFolder, `index.html`)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			w.WriteHeader(404)
			return
		}
		w.Header().Add("Content-Type", `text/html`)
		_, _ = w.Write(content)
	}
}

func (s *Server) fileHandler(hostFolder string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fullPath := path.Join(s.ServeFolder, hostFolder, r.URL.Path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			err404, _ := os.ReadFile(path.Join(s.ServeFolder, hostFolder, `404.html`))
			w.WriteHeader(404)
			_, _ = w.Write(err404)
			return
		}
		mimeType := mime.TypeByExtension(filepath.Ext(fullPath))
		w.Header().Add("Content-Type", mimeType)
		_, _ = w.Write(content)
	}
}

func (s *Server) CollectHostWhitelist() []string {
	whitelist := make([]string, 0, 10)
	for _, hosts := range s.Hosts {
		whitelist = append(whitelist, hosts.HostWhitelist...)
	}
	return whitelist
}

func (s *Server) GenerateRouter() *mux.Router {
	e := mux.NewRouter()

	for _, info := range s.Hosts {
		for _, host := range info.HostWhitelist {
			sub := e.Host(host).Subrouter()
			sub.Path(`/`).HandlerFunc(s.homepageHandler(info.Folder))
			sub.PathPrefix(`/`).Methods(`GET`).HandlerFunc(s.fileHandler(info.Folder))
		}
	}

	return e
}
