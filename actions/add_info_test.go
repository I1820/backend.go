/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 26-11-2018
 * |
 * | File Name:     add_info_test.go
 * +===============================================
 */

package actions

import (
	"fmt"

	"github.com/I1820/backend/models"
)

const iKey = "hello"
const iValue = "elie"

func (as *ActionSuite) Test_AdditionalsResource_Create() {
	var ur models.User

	// Signup (POST /api/v1/auth/register)
	resr := as.JSON("/api/v1/auth/register").Post(signupReq{
		Firstname: uFName,
		Lastname:  uLName,
		Username:  fmt.Sprintf("additiona%s", uuName),
		Email:     uEmail,
		Password:  uPass,
	})
	as.Equalf(200, resr.Code, "Error: %s", resr.Body.String())
	resr.Bind(&ur)

	// Login (POST /api/v1/auth/login)
	resl := as.JSON("/api/v1/auth/login").Post(loginReq{
		Username: fmt.Sprintf("additiona%s", uuName),
		Password: uPass,
	})
	as.Equalf(200, resl.Code, "Error: %s", resl.Body.String())
	resl.Bind(&ur)

	// Create (POST /api/v1/info/{key})
	reqc := as.JSON("/api/v1/info/%s", iKey)
	reqc.Headers["Authorization"] = fmt.Sprintf("Bearer %s", ur.AccessToken)
	resc := reqc.Post(iValue)
	as.Equalf(200, resc.Code, "Error: %s", resc.Body.String())

	// Show (GET /api/v1/info/{key})
	var v string
	reqs := as.JSON("/api/v1/info/%s", iKey)
	reqs.Headers["Authorization"] = fmt.Sprintf("Bearer %s", ur.AccessToken)
	ress := reqs.Get()
	as.Equalf(200, ress.Code, "Error: %s", ress.Body.String())
	ress.Bind(&v)
	as.Equal(v, iValue)
}
