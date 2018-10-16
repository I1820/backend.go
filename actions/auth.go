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

// login request payload
type loginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
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

// AuthMiddleware is an Authorization middleware using JWT. it validates JWT tokens and checks
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
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validating if algorithm used for signing is same as the algorithm in token
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}
			return key, nil
		})
		// if error validating jwt token, return with status unauthorized
		if err != nil {
			return c.Error(http.StatusUnauthorized, err)
		}

		claims := token.Claims.(jwt.MapClaims)
		// if token expired, return with status unathorized
		if int64(claims["exp"].(float64)) < time.Now().Unix() {
			return c.Error(http.StatusUnauthorized, err)
		}

		// set the user as context parameter.
		// so that the actions can use the user object from jwt token
		c.Set("user", claims["user"])

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

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = u
	claims["exp"] = time.Now().Add(time.Second * 120).Unix() // tokens expire in 2 minutes

	// Generate encoded token and send it as response
	encodedToken, err := token.SignedString([]byte(envy.Get("JWT_SECRET", "i1820")))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	u.Token = encodedToken

	u.Password = "" // Don't send password

	return c.Render(http.StatusOK, r.JSON(u))
}
