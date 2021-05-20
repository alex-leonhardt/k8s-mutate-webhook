package web

import (
	"encoding/json"
	"net/http"
)

func HandleHealthz(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
