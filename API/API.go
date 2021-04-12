package API

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/efvincent/archex5/commands"
	"github.com/efvincent/archex5/eventStore"
	"github.com/efvincent/archex5/eventStore/MemoryEventStore"
	"github.com/efvincent/archex5/processor"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var cmdProc = processor.MakeCmdProc()

const COMMAND_TYPE_ATTRIB = "commandType"

func Run(host string, port string) {
	router := mux.NewRouter()
	r := router.HandleFunc("/api/command", commandHandler)
	r.Methods("POST")

	r = router.HandleFunc("/api/{namespace}/products", getProductsHandler)
	r.Methods("GET")

	r = router.HandleFunc("/api/{namespace}/products/{sku}", getProductHandler)

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Server running. Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ns := vars["namespace"]
	sku := vars["sku"]
	if len(ns) == 0 || len(sku) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if p, err := cmdProc.GetProduct(ns, sku); err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Could not retrieve: %v", err)
	}

}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ns := vars["namespace"]
	if len(ns) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var es eventStore.EventStore
	// some kind of DI might make sense here ... demo approach only
	es = MemoryEventStore.SingletonMemoryEventStore
	streamIds, err := es.GetStreams(ns)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"namespace": ns,
		"skus":      streamIds,
	})
}

// Reads the raw body, unmarshals it as generic json, looks for a field called
// commandType, and sends the raw json and raw event type to commands.UnmarshalAsTypedCommand
// to get a typed command, and then forwards that to the command processor
func commandHandler(w http.ResponseWriter, r *http.Request) {
	var raw map[string]interface{}
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Could not read from the request body")
		return
	}

	rawStr := buf.String()
	if err := json.Unmarshal([]byte(rawStr), &raw); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Could not unmarshal request body as json")
		return
	}

	// add a timestamp and unique ID to the incoming command
	raw["ts"] = time.Now().Unix()
	raw["uid"] = uuid.New().String()

	if typeKey, tOk := raw[COMMAND_TYPE_ATTRIB]; tOk {
		switch typeKey.(type) {
		case string:
			cmd, err := commands.UnmarshalAsTypedCommand(fmt.Sprintf("%v", typeKey), []byte(rawStr))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Could not unmarshal request body as a valid command: %v", err)
				return
			}

			// there is such a type mapping. Attempt to decode it.
			if err := cmdProc.ProcessProductCommand(cmd); err != nil {
				log.Printf("API error: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "API error: %v", err)
			} else {
				w.WriteHeader(http.StatusOK)
			}
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "'%s' attribute should be a string with a valid command type value", COMMAND_TYPE_ATTRIB)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "request body json does not contain an attribute '%s'", COMMAND_TYPE_ATTRIB)
}
