/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 12-11-2018
 * |
 * | File Name:     wf.go
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

// WFHandler handles weather forcecasting requests by proxies them to wf component.
// This function is mapped to the path POST /wf/{service}
func WFHandler(c buffalo.Context) error {
	var rq interface{}

	// generic request so there is no restriction on wf.
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	wfclient := resty.New().SetHostURL(envy.Get("WF_URL", "http://127.0.0.1:6976")).SetError(types.Error{})

	// send request to wf
	var w interface{}
	resp, err := wfclient.R().SetBody(rq).SetResult(&w).SetPathParams(map[string]string{
		"service": c.Param("service"),
	}).Post("api/{service}")
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if resp.IsError() {
		return c.Render(resp.StatusCode(), r.JSON(resp.Error()))
	}

	return c.Render(http.StatusOK, r.JSON(w))
}
