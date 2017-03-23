package response

import (
	"encoding/json"
	"net/http"
)

// H is a helper for creating a JSON response
type H map[string]interface{}

// JSON writes a JSON response to ResponseWriter
func JSON(w http.ResponseWriter, code int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if nil != obj || http.StatusNoContent == code {
		json.NewEncoder(w).Encode(obj)
	}
}
