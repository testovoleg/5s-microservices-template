package httpUtils

import "github.com/labstack/echo/v4"

// QueryAwareBinder always binds query params in addition to path params and body.
// Echo's DefaultBinder only binds query params for GET/DELETE/HEAD, which silently
// drops `query:"..."` tagged fields on POST/PATCH handlers that carry identifiers
// in the query string alongside a JSON body.
type QueryAwareBinder struct {
	echo.DefaultBinder
}

func (b *QueryAwareBinder) Bind(i interface{}, c echo.Context) error {
	if err := b.BindPathParams(c, i); err != nil {
		return err
	}
	if err := b.BindQueryParams(c, i); err != nil {
		return err
	}
	return b.BindBody(c, i)
}
