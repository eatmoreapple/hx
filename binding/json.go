package binding

import (
	"net/http"

	"github.com/eatmoreapple/hx/internal/serializer"
)

type JSONBinder struct{}

func (j JSONBinder) Bind(r *http.Request, a any) error {
	return serializer.JSONSerializer().Deserialize(r.Body, a)
}
