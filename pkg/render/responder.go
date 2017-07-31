package render

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// M is a simple abstraction for a map interface
type M map[string]interface{}

// JSON marshals 'v' to JSON, automatically escaping HTML and setting the
// Content-Type as application/json.
func JSON(w http.ResponseWriter, code int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(buf.Bytes())
}
