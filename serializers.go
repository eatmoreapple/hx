package hx

import "github.com/eatmoreapple/hx/internal/serializer"

// SetJSONSerializer sets the JSON serializer used by the framework.
// This function allows you to customize the JSON serialization behavior.
// such as jsonit, easyjson, or any other custom serializer.
func SetJSONSerializer(s serializer.Serializer) {
	serializer.SetJSONSerializer(s)
}
