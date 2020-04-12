package utils

import (
	"encoding/json"
	"net/http"
)

//Message returns the code, message, hasError, data as a map
func Message(code int, message string, hasError bool, data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{"code": code, "message": message, "hasError": hasError, "data": data}
}

//Respond adds headers and encodes as json
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
