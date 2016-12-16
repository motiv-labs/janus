package response

import (
	"encoding/json"
	"net/http"
)

type H map[string]interface{}

func JSON(w http.ResponseWriter, code int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if nil != obj || http.StatusNoContent == code {
		json.NewEncoder(w).Encode(obj)
	}
}
