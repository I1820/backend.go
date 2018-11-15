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
	"fmt"
	"net/http"

	"github.com/I1820/backend/models"
	"github.com/I1820/types"
	"github.com/gobuffalo/buffalo"
)

// ThingsResource controls the users access to things and proxies their request to pm
type ThingsResource struct {
	buffalo.Resource
}

// List lists things of given project. This function is mapped to the path
// GET /projects/{project_id}/things
func (v ThingsResource) List(c buffalo.Context) error {
	projectID := c.Param("project_id")

	// get user from request context
	u, ok := c.Value("user").(models.User)
	if !ok {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("There is no valid user in request context"))
	}

	for _, p := range u.Projects {
		if p == projectID {
			var ts []types.Thing

			// gets project things from pm
			// I1820/pm/ThingsResource.List
			resp, err := pmclient.R().SetResult(&ts).SetPathParams(map[string]string{
				"projectID": projectID,
			}).Get("api/projects/{projectID}/things")
			if err != nil {
				return c.Error(http.StatusInternalServerError, err)
			}

			if resp.IsError() {
				return c.Render(resp.StatusCode(), r.JSON(resp.Error()))
			}

			return c.Render(http.StatusOK, r.JSON(ts))
		}
	}

	return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", projectID))
}

// Create creates a new thing in pm. This function is mapped to the
// path POST /projects/{project_id}/things
func (v ThingsResource) Create(c buffalo.Context) error {
	projectID := c.Param("project_id")

	// generic request so we can change it when pm `thingReq` is changed.
	var rq interface{}
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	// get user from request context
	u, ok := c.Value("user").(models.User)
	if !ok {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("There is no valid user in request context"))
	}

	for _, p := range u.Projects {
		if p == projectID {
			var t types.Thing

			// creates a thing in pm
			// I1820/pm/ThingsResource.Create
			resp, err := pmclient.R().SetBody(rq).SetResult(&t).SetPathParams(map[string]string{
				"projectID": projectID,
			}).Post("api/projects/{projectID}/things")
			if err != nil {
				return c.Error(http.StatusInternalServerError, err)
			}

			if resp.IsError() {
				return c.Render(resp.StatusCode(), r.JSON(resp.Error()))
			}

			return c.Render(http.StatusOK, r.JSON(t))
		}
	}

	return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", projectID))
}

// Show shows given thing details that are fetched from pm.
// This function is mapped to the path GET /projects/{proejct_id}/things/{thing_id}
func (v ThingsResource) Show(c buffalo.Context) error {
	projectID := c.Param("project_id")
	thingID := c.Param("thing_id")

	// get user from request context
	u, ok := c.Value("user").(models.User)
	if !ok {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("There is no valid user in request context"))
	}

	for _, p := range u.Projects {
		if p == projectID {
			var t types.Thing

			// fetches a thing from pm
			// I1820/pm/ThingsResource.Show
			resp, err := pmclient.R().SetResult(&t).SetPathParams(map[string]string{
				"projectID": projectID,
				"thingID":   thingID,
			}).Get("api/projects/{projectID}/things/{thingID}")
			if err != nil {
				return c.Error(http.StatusInternalServerError, err)
			}

			if resp.IsError() {
				return c.Render(resp.StatusCode(), r.JSON(resp.Error()))
			}

			return c.Render(http.StatusOK, r.JSON(t))
		}
	}

	return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", projectID))

}
