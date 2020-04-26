package utils

import (
	"encoding/json"
	"net/http"
)

//Message returns the code, message, hasError, data as a map
func Message(code int, message string, hasError bool, data interface{}) map[string]interface{} {
	responseData := map[string]interface{}{"response": data}
	return map[string]interface{}{"code": code, "message": message, "hasError": hasError, "data": responseData}
}

func MessageWithTokens(code int, message string, hasError bool, data interface{}, tokens interface{}) map[string]interface{} {
	responseData := map[string]interface{}{"response": data, "tokens": tokens}
	return map[string]interface{}{"code": code, "message": message, "hasError": hasError, "data": responseData}
}

//Respond adds headers and encodes as json
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
