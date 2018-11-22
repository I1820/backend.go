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

	"github.com/gobuffalo/buffalo"
)

// PMHealthHandler checks status of pm service.
// This function is mapped to the path GET /health/pm
func PMHealthHandler(c buffalo.Context) error {
	var a string
	resp, err := pmclient.R().SetResult(&a).Get("about")
	if err != nil { // there is no connection to the component
		return c.Render(http.StatusOK, r.JSON(false))
	}

	if resp.IsError() { // something bad is happing
		return c.Render(http.StatusOK, r.JSON(false))
	}

	return c.Render(http.StatusOK, r.JSON(a == "18.20 is leaving us"))
}

// WFHealthHandler checks status of wf service.
// This function is mapped to the path GET /health/wf
func WFHealthHandler(c buffalo.Context) error {
	var a string
	resp, err := wfclient.R().SetResult(&a).Get("about")
	if err != nil { // there is no connection to the component
		return c.Render(http.StatusOK, r.JSON(false))
	}

	if resp.IsError() { // something bad is happing
		return c.Render(resp.StatusCode(), r.JSON(resp.Error()))
	}

	return c.Render(http.StatusOK, r.JSON(a == "18.20 is leaving us"))
}

// DMHealthHandler checks status of dm service.
// This function is mapped to the path GET /health/dm
func DMHealthHandler(c buffalo.Context) error {
	var a string
	resp, err := dmclient.R().SetResult(&a).Get("about")
	if err != nil { // there is no connection to the component
		return c.Render(http.StatusOK, r.JSON(false))
	}

	if resp.IsError() { // something bad is happing
		return c.Render(resp.StatusCode(), r.JSON(resp.Error()))
	}

	return c.Render(http.StatusOK, r.JSON(a == "18.20 is leaving us"))
}
