/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-07-2018
 * |
 * | File Name:     actions/project.go
 * +===============================================
 */

package actions

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

// ProjectsHandler sends request (proxies) to pm
func ProjectsHandler(c buffalo.Context) error {
	path := c.Param("path")
	user := c.Value("username").(string)

	url, err := url.Parse(fmt.Sprintf("http://%s:8080", envy.Get("PM_HOST", "127.0.0.1")))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	c.Request().URL.Path = fmt.Sprintf("api/%s/projects/%s", user, path)
	return buffalo.WrapHandler(
		httputil.NewSingleHostReverseProxy(url),
	)(c)
}
