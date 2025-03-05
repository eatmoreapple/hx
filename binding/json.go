package binding

import (
	"encoding/json"
	"net/http"
)

type JSONBinder struct{}

func (j JSONBinder) Bind(r *http.Request, a any) error {
	return json.NewDecoder(r.Body).Decode(a)
}
