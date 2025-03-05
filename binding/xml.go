package binding

import (
	"encoding/xml"
	"net/http"
)

type XMLBinder struct{}

func (b XMLBinder) Bind(r *http.Request, obj any) error {
	return xml.NewDecoder(r.Body).Decode(obj)
}
