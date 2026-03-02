package json

import (
	"encoding/json"
	"io"
	"net/http"
)

func Read(r *http.Request, v any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func Write(w http.ResponseWriter, status int, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(data)
	return err
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	data, err := json.Marshal(map[string]string{"error": msg})
	if err != nil {
		_, _ = w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	_, _ = w.Write(data)
}
