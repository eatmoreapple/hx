package binding

import (
	"net/http"

	"github.com/eatmoreapple/hx/encoding/form"
)

type QueryBinder struct{}

func (q QueryBinder) Bind(r *http.Request, a any) error {
	query := r.URL.Query()
	return form.Unmarshal(query, a)
}
