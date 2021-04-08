package API

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/efvincent/archex5/commands"
	"github.com/efvincent/archex5/models"
	"github.com/efvincent/archex5/processor"
	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("<h1>Product: %s</h1>", models.SampleProduct.SKU)))
}

// Reads the raw body, unmarshals it as generic json, looks for a field called
// eventType, and sends the raw json and event type to commands.UnmarshalAsTypedCommand
// to get a typed command, and then forwards that to the command processor
func commandHandler(w http.ResponseWriter, r *http.Request) {
	var raw map[string]interface{}
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rawStr := buf.String()
	json.Unmarshal([]byte(rawStr), &raw)
	if typeKey, tOk := raw["eventType"]; tOk {
		switch typeKey.(type) {
		case string:
			if cmd := commands.UnmarshalAsTypedCommand(fmt.Sprintf("%v", typeKey), []byte(rawStr)); cmd != nil {
				// there is such a type mapping. Attempt to decode it.
				if err := processor.ProcessProductCommand(cmd); err != nil {
					log.Printf("command processor error: %s", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}

func Run(host string, port string) {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	r := router.HandleFunc("/api", commandHandler)
	r.Methods("POST")
	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Server running. Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
