package handlers

import (
	"net/http"
	"encoding/json"
	"log"
)

func validate(data map[string]interface{}, key string) (string, bool) {
	field, ok := data[key].(string)
	return field, ok
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is main page.\nWlecome!"))
}

func PostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	name, ok := validate(data, "name")
	if !ok {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	surname, ok := validate(data, "surname")
	if !ok {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Hello, %s %s", name, surname)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, " + name + " " + surname))

}