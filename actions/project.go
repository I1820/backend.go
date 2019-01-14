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

	"github.com/I1820/backend/models"
	pmmodels "github.com/I1820/pm/models"
	"github.com/gobuffalo/buffalo"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/mongoopt"
)

// ProjectsResource controls the users access to projects and proxies their request to pm
type ProjectsResource struct {
	buffalo.Resource
}

// project request payload
type projectReq struct {
	Name string            `json:"name" validate:"required"` // project name
	Envs map[string]string `json:"envs"`                     // project environment variables
}

// List lists user projects. It first gets projects list from pm then returns only the projects that are exist in
// user porject list.
// This function is mapped to the path GET /projects
func (v ProjectsResource) List(c buffalo.Context) error {
	var ps []pmmodels.Project
	var ups []pmmodels.Project

	// get user from request context
	u, ok := c.Value("user").(models.User)
	if !ok {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("There is no valid user in request context"))
	}

	// gets all projects from pm
	// I1820/pm/ProjectsResource.List
	resp, err := pmclient.R().SetResult(&ps).Get("api/projects")
	if err != nil {
		return c.Error(http.StatusServiceUnavailable, fmt.Errorf("PM Service is not available"))
	}

	if resp.IsError() {
		return c.Render(resp.StatusCode(), r.JSON(resp.Error()))
	}

	// collects projects that are exist in user project list
	for _, p := range ps {
		for _, id := range u.Projects {
			if p.ID == id {
				ups = append(ups, p)
			}
		}
	}

	return c.Render(http.StatusOK, r.JSON(ups))
}

// Create creates new project in pm and if it successful then adds newly created project to user projects
// This function is mapped to the path POST /projects
func (v ProjectsResource) Create(c buffalo.Context) error {
	var rq projectReq
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	if err := validate.Struct(rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	// get user from request context
	u, ok := c.Value("user").(models.User)
	if !ok {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("There is no valid user in request context"))
	}

	var p pmmodels.Project

	// creates a project in pm with `projectReq`
	// I1820/pm/ProjectsResource.Create
	resp, err := pmclient.R().SetBody(map[string]interface{}{
		"name":  rq.Name,
		"owner": u.Email,
		"envs":  rq.Envs,
	}).SetResult(&p).Post("api/projects")
	if err != nil {
		return c.Error(http.StatusServiceUnavailable, fmt.Errorf("PM Service is not available"))
	}

	if resp.IsError() {
		return c.Render(resp.StatusCode(), r.JSON(resp.Error()))
	}

	// adds newly created project into user projects
	dr := db.Collection("users").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("username", u.Username),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$addToSet", bson.EC.String("projects", p.ID)),
	), findopt.ReturnDocument(mongoopt.After))
	if err := dr.Decode(&u); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	// and then update jwt token but how?
	// it is on client for now

	return c.Render(http.StatusOK, r.JSON(p))
}

// Show shows given project details that are fetched from pm.
// This function is mapped to the path GET /projects/{proejct_id}
func (v ProjectsResource) Show(c buffalo.Context) error {
	projectID := c.Param("project_id")

	// get user from request context
	u, ok := c.Value("user").(models.User)
	if !ok {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("There is no valid user in request context"))
	}

	for _, p := range u.Projects {
		if p == projectID {
			var p pmmodels.Project

			// shows a project from pm
			// I1820/pm/ProjectsResource.Show
			resp, err := pmclient.R().SetResult(&p).SetPathParams(map[string]string{
				"projectID": projectID,
			}).Get("api/projects/{projectID}")
			if err != nil {
				return c.Error(http.StatusServiceUnavailable, fmt.Errorf("PM Service is not available"))
			}

			if resp.IsError() {
				return c.Render(resp.StatusCode(), r.JSON(resp.Error()))
			}

			return c.Render(http.StatusOK, r.JSON(p))
		}
	}

	return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", projectID))
}

// Destroy deletes given project from pm and if it successful then removes it from user projects
// This function is mapped to the path DELETE /projects/{proejct_id}
func (v ProjectsResource) Destroy(c buffalo.Context) error {
	projectID := c.Param("project_id")

	// get user from request context
	u, ok := c.Value("user").(models.User)
	if !ok {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("There is no valid user in request context"))
	}

	for _, p := range u.Projects {
		if p == projectID {
			var p pmmodels.Project

			// removes a thing from pm
			// I1820/pm/ThingsResource.Destroy
			resp, err := pmclient.R().SetResult(&p).SetPathParams(map[string]string{
				"projectID": projectID,
			}).Delete("api/projects/{projectID}")
			if err != nil {
				return c.Error(http.StatusServiceUnavailable, fmt.Errorf("PM Service is not available"))
			}

			if resp.IsError() {
				return c.Render(resp.StatusCode(), r.JSON(resp.Error()))
			}

			// removes given project from user projects
			dr := db.Collection("users").FindOneAndUpdate(c, bson.NewDocument(
				bson.EC.String("username", u.Username),
			), bson.NewDocument(
				bson.EC.SubDocumentFromElements("$pull", bson.EC.String("projects", projectID)),
			), findopt.ReturnDocument(mongoopt.After))
			if err := dr.Decode(&u); err != nil {
				return c.Error(http.StatusInternalServerError, err)
			}

			// and then update jwt token but how?
			// it is on client for now

			return c.Render(http.StatusOK, r.JSON(p))
		}
	}

	return c.Error(http.StatusNotFound, fmt.Errorf("Project %s not found", projectID))
}
