package binding

import "net/http"

type QueryBinder struct{}

func (q QueryBinder) Bind(r *http.Request, a any) error {
	query := r.URL.Query()
	return mapTo(query, a)
}
