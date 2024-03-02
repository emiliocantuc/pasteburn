package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var store = MessageStore{
	messages:    make(map[string]Message),
	maxMessages: 100,
}

var baseurl = "http://localhost:8080"

var expireStr2Duration = map[string]time.Duration{
	"5 minutes":  5 * time.Minute,
	"30 minutes": 30 * time.Minute,
	"1 hour":     time.Hour,
	"5 hours":    5 * time.Hour,
}

func respondWithTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "error")
	} else {
		t.Execute(w, data)
	}
}

func send(w http.ResponseWriter, req *http.Request, messageText string, expires string) {
	fmt.Println(expires)
	message := Message{
		text:   messageText,
		expire: expireStr2Duration[expires],
	}
	if id, ok := store.Add(message); ok {
		respondWithTemplate(w, "templates/msg_sent.html", MessageSentPage{
			BaseURL: baseurl,
			Expire:  time.Now().Add(message.expire).Format("15:04:05 02 Jan 2006 MST"),
			Id:      id,
		})
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "error")
	}
}

func read(w http.ResponseWriter, req *http.Request, id string) {
	if message, ok := store.Pop(id); ok {
		respondWithTemplate(w, "templates/msg_view.html", MessageViewPage{Text: message.text, Id: id})
	} else {
		fmt.Fprint(w, "Message not found")
	}

}

func index(w http.ResponseWriter, req *http.Request) {

	fmt.Println(req.Form, req.FormValue("m"))

	if id := req.FormValue("id"); req.Method == "GET" && id != "" {

		read(w, req, id)

	} else if (req.Method == "POST" || req.Method == "PUT") && req.FormValue("text") != "" && req.FormValue("expire") != "" {

		send(w, req, req.FormValue("text"), req.FormValue("expire"))

	} else {

		http.ServeFile(w, req, "index.html")
	}

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
