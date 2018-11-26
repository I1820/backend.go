/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 26-11-2018
 * |
 * | File Name:     add_info.go
 * +===============================================
 */

package actions

import (
	"fmt"
	"net/http"

	"github.com/I1820/backend/models"
	"github.com/gobuffalo/buffalo"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/mongoopt"
)

// AdditionalsResource contronls additional information on users.
type AdditionalsResource struct {
}

// Create creates new key on user additional information.
// Please note that you cannot store array on additional information field.
// This function is mapped to the path POST /info/{key}
func (AdditionalsResource) Create(c buffalo.Context) error {
	key := c.Param("key")

	// applications can store everything on additional_info
	var rq interface{}
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	// get user from request context
	u, ok := c.Value("user").(models.User)
	if !ok {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("There is no valid user in request context"))
	}

	dr := db.Collection("users").FindOneAndUpdate(c, bson.NewDocument(
		bson.EC.String("username", u.Username),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.Interface(fmt.Sprintf("additional_info.%s", key), rq)),
	), findopt.ReturnDocument(mongoopt.After))

	if err := dr.Decode(&u); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(true))
}

// Show returns the value of given key.
// This function is mapped to the path GET /info/{key}
func (AdditionalsResource) Show(c buffalo.Context) error {
	key := c.Param("key")

	// get user from request context
	u, ok := c.Value("user").(models.User)
	if !ok {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("There is no valid user in request context"))
	}

	dr := db.Collection("users").FindOne(c, bson.NewDocument(
		bson.EC.String("username", u.Username),
	))

	if err := dr.Decode(&u); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	fmt.Println(u.AdditionalInfo[key])

	return c.Render(http.StatusOK, r.JSON(u.AdditionalInfo[key]))

}
