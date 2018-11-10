package actions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/I1820/backend/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// AuthResource represents login, logout and signup
type AuthResource struct{}

// refresh request payload
type refreshReq struct {
	Token string `json: "token"`
}

// login request payload
type loginReq struct {
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember"`
}

// signup request payload
type signupReq struct {
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Username  string `json:"username" validate:"required,alphanum"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

// AuthMiddleware and getJwtToken are taken from
// https://github.com/gobuffalo/buffalo/blob/master/middleware/tokenauth/tokenauth.go

// AuthMiddleware is an Authorization middleware using JWT. it validates access tokens and checks
// their expiration time.
func AuthMiddleware(next buffalo.Handler) buffalo.Handler {
	key := []byte(envy.Get("JWT_SECRET", "i1820"))

	return func(c buffalo.Context) error {
		// get Authorisation header value
		authString := c.Request().Header.Get("Authorization")

		tokenString, err := getJwtToken(authString)
		// if error on getting the token, return with status unauthorized
		if err != nil {
			return c.Error(http.StatusUnauthorized, err)
		}

		// validating and parsing the tokenString
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validating if algorithm used for signing is same as the algorithm in token
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}
			return key, nil
		})
		if err != nil {
			return c.Error(http.StatusUnauthorized, err)
		}

		// if error validating jwt token, return with status unauthorized
		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
			// set the user as context parameter.
			// so that the actions can use the user object from jwt token
			c.Set("user", claims.U)
		} else {
			return c.Error(http.StatusUnauthorized, err)
		}

		// calling next handler
		return next(c)
	}
}

// getJwtToken gets the token from the Authorisation header
// removes the Bearer part from the authorisation header value.
// returns No token error if Token is not found
// returns Token Invalid error if the token value cannot be obtained by removing `Bearer `
func getJwtToken(authString string) (string, error) {
	if authString == "" {
		return "", errors.New("token not found in request")
	}
	splitToken := strings.Split(authString, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("token invalid")
	}
	tokenString := splitToken[1]
	return tokenString, nil
}

// Signup creates new user with given information amd store it in database.
// Signup do not create any token for new user.
// This function is mapped to the path
// POST /register
func (a AuthResource) Signup(c buffalo.Context) error {
	var rq signupReq

	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}
	if err := validate.Struct(rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	// is there any need for hashing the password before store it in database?
	u := models.User{
		Firstname: rq.Firstname,
		Lastname:  rq.Lastname,
		Username:  rq.Username,
		Email:     rq.Email,
		Password:  rq.Password,
		Projects:  make([]string, 0),
	}
	if _, err := db.Collection("users").InsertOne(c, u); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(u))
}

// Login checks given credentials and generate jwt token
// This function is mapped to the path
// POST /login
func (a AuthResource) Login(c buffalo.Context) error {
	var rq loginReq
	var u models.User

	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}
	if err := validate.Struct(rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	res := db.Collection("users").FindOne(c, bson.NewDocument(
		bson.EC.String("username", rq.Username),
		bson.EC.String("password", rq.Password),
	))

	if err := res.Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("invalid username or password"))
		}
		return c.Error(http.StatusInternalServerError, err)
	}

	accessToken, err := NewAccessToken(u)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	u.AccessToken = accessToken

	d := time.Hour * 1
	if rq.RememberMe {
		d = time.Hour * 24 * 7 // remembers the user for a week
	}
	refreshToken, err := NewRefreshToken(u.Username, d)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	u.RefreshToken = refreshToken

	u.Password = "" // Don't send password

	return c.Render(http.StatusOK, r.JSON(u))
}

// Refresh refreshes given token with new expiration time.
// It updates user information from database.
func (a AuthResource) Refresh(c buffalo.Context) error {
	key := []byte(envy.Get("JWT_SECRET", "i1820"))

	var rq refreshReq
	if err := c.Bind(&rq); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}
	tokenString := rq.Token

	// validating and parsing the tokenString
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validating if algorithm used for signing is same as the algorithm in token
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})
	if err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}

	username := ""
	// if error validating jwt token, return with status unauthorized
	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		username = claims.ID
	} else {
		return c.Error(http.StatusUnauthorized, err)
	}

	var u models.User

	res := db.Collection("users").FindOne(c, bson.NewDocument(
		bson.EC.String("username", username),
	))

	if err := res.Decode(&u); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	accessToken, err := NewAccessToken(u)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	u.AccessToken = accessToken
	u.RefreshToken = tokenString // uses the old refresh token.

	u.Password = "" // Don't send password

	return c.Render(http.StatusOK, r.JSON(u))
}
