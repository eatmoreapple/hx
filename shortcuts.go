package hx

import "github.com/eatmoreapple/hx/httpx/extractor"

type Value = extractor.Value

// Single-value request extractor shortcuts.
type (
	FromPath[T Value]   = extractor.PathValueExtractor[T]
	FromHeader[T Value] = extractor.HeaderValueExtractor[T]
	FromQuery[T Value]  = extractor.QueryValueExtractor[T]
	FromForm[T Value]   = extractor.FormValueExtractor[T]
	FromCookie[T Value] = extractor.CookieValueExtractor[T]
)

// Complete request data extractor shortcuts.
type (
	Header  = extractor.HeaderExtractor
	Cookies = extractor.CookieExtractor
	Query   = extractor.QueryExtractor
	Form    = extractor.FormExtractor
	Empty   = extractor.Empty
)
