/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 15-10-2018
 * |
 * | File Name:     auth_test.go
 * +===============================================
 */

package actions

import (
	"context"

	"github.com/I1820/backend/models"
	"github.com/mongodb/mongo-go-driver/bson"
)

const uFName = "پرهام"
const uLName = "الوانی"
const uuName = "1995parham"
const uEmail = "parham.alvani@gmail.com"
const uPass = "123123"

func (as *ActionSuite) Test_AuthResource_Signup_Login() {
	var ur models.User

	// Create (POST /api/v1/auth/register)
	resr := as.JSON("/api/v1/auth/register").Post(signupReq{
		Firstname: uFName,
		Lastname:  uLName,
		Username:  uuName,
		Email:     uEmail,
		Password:  uPass,
	})
	as.Equalf(200, resr.Code, "Error: %s", resr.Body.String())
	resr.Bind(&ur)

	// check database for project existence
	var ud models.User
	dr := db.Collection("users").FindOne(context.Background(), bson.NewDocument(
		bson.EC.String("username", uuName),
	))

	as.NoError(dr.Decode(&ud))

	as.Equal(ud, ur)

	// Login (POST /api/v1/auth/login)
	resl := as.JSON("/api/v1/auth/login").Post(loginReq{
		Username: uuName,
		Password: uPass,
	})
	as.Equalf(200, resl.Code, "Error: %s", resl.Body.String())
	resr.Bind(&ur)
	as.NotNil(ur.AccessToken)  // token must be there
	as.NotNil(ur.RefreshToken) // token must be there
}
