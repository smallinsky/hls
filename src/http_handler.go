package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func handler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[1:]
	data, err := ioutil.ReadFile(url)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header()["Access-Control-Allow-Orgin"] = []string{"*"}
	w.Write(data)
}

type handlers struct {
	mtx sync.Mutex
	reg map[string]struct{}
}

func (h *handlers) register(name string) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	if _, ok := h.reg[name]; ok {
		log.Printf("%s handler exist, registration failed", name)
		return
	}
	url := "/" + name + "/"
	log.Printf("Handler added: %s", url)
	http.HandleFunc(url, handler)
}

var reghandlers handlers

func addHandler(url string) {
	reghandlers.register(url)
}
