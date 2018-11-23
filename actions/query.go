/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 23-11-2018
 * |
 * | File Name:     query.go
 * +===============================================
 */

package actions

import (
	"fmt"
	"net/http"

	"github.com/I1820/backend/models"
	"github.com/gobuffalo/buffalo"
)

// QueryHandler handles data queries in dm.
// This function is mapped to the path ANY /projects/{project_id}/things/{thing_id}/queries/{path:.+}
func QueryHandler(c buffalo.Context) error {
	projectID := c.Param("project_id")
	thingID := c.Param("thing_id")
	path := c.Param("path")
	method := c.Value("current_route").(buffalo.RouteInfo).Method

	// generic request so we can change it when pm `connectivityReq` is changed.
	var rq interface{}
	if err := c.Bind(&rq); err != nil {
		// do noting on binding errors
	}

	// get user from request context
	u, ok := c.Value("user").(models.User)
	if !ok {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("There is no valid user in request context"))
	}

	for _, p := range u.Projects {
		if p == projectID {
			var result interface{}
			// do query
			// I1820/dm/
			resp, err := dmclient.R().SetBody(rq).SetResult(&result).SetPathParams(map[string]string{
				"projectID": projectID,
				"thingID":   thingID,
				"path":      path,
			}).Execute(method, "api/projects/{projectID}/things/{thingID}/queries/{path}")
			if err != nil {
				return c.Error(http.StatusInternalServerError, err)
			}
			if resp.IsError() {
				return c.Render(resp.StatusCode(), r.JSON(resp.Error()))
			}

			return c.Render(http.StatusOK, r.JSON(result))
		}
	}

	return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", projectID))
}
