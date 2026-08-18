package hx

import "github.com/eatmoreapple/hx/httpx/extractor"

// Single-value request extractor shortcuts.
type (
	FromPath[T extractor.Value]   = extractor.PathValueExtractor[T]
	FromHeader[T extractor.Value] = extractor.HeaderValueExtractor[T]
	FromQuery[T extractor.Value]  = extractor.QueryValueExtractor[T]
	FromForm[T extractor.Value]   = extractor.FormValueExtractor[T]
	FromCookie[T extractor.Value] = extractor.CookieValueExtractor[T]
)

// Complete request data extractor shortcuts.
type (
	Header  = extractor.HeaderExtractor
	Cookies = extractor.CookieExtractor
	Query   = extractor.QueryExtractor
	Form    = extractor.FormExtractor
	Empty   = extractor.Empty
)
