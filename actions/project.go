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
	"github.com/go-resty/resty"
	"github.com/gobuffalo/buffalo"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/mongoopt"
)

// ProjectsResource controls the users access to projects and proxies their request to pm
type ProjectsResource struct {
	buffalo.Resource
	pmclient *resty.Client
}

// project request payload
type projectReq struct {
	Name string            `json:"name" validate:"required"` // project name
	Envs map[string]string `json:"envs"`                     // project environment variables
}

// Create creates new project in pm and if it successful then adds newly created project to user projects
// path POST /projects
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
	resp, err := v.pmclient.R().SetBody(map[string]interface{}{
		"name":  rq.Name,
		"owner": u.Email,
		"envs":  rq.Envs,
	}).SetResult(&p).Post("api/projects")
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
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
