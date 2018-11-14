/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 14-11-2018
 * |
 * | File Name:     health.go
 * +===============================================
 */

package actions

import (
	"net/http"

	"github.com/I1820/types"
	"github.com/go-resty/resty"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

// PMHealthHandler checks status of pm service.
// This function is mapped to the path GET /health/pm
func PMHealthHandler(c buffalo.Context) error {
	pmclient := resty.New().SetHostURL(envy.Get("PM_URL", "http://127.0.0.1:1375")).SetError(types.Error{})

	var a string
	resp, err := pmclient.R().SetResult(&a).Get("about")
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if resp.IsError() {
		return c.Render(resp.StatusCode(), r.JSON(resp.Error()))
	}

	return c.Render(http.StatusOK, r.JSON(a == "18.20 is leaving us"))
}
