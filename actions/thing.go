/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 18-10-2018
 * |
 * | File Name:     thing.go
 * +===============================================
 */

package actions

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

// ThingsHandler proxies each thing request to pm component.
func ThingsHandler(c buffalo.Context) error {
	path := c.Param("path")

	url, err := url.Parse(envy.Get("PM_URL", "http://127.0.0.1:1375"))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	c.Request().URL.Path = path
	return buffalo.WrapHandler(
		httputil.NewSingleHostReverseProxy(url),
	)(c)
}
