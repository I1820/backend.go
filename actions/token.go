/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 10-11-2018
 * |
 * | File Name:     token.go
 * +===============================================
 */

package actions

import (
	"fmt"
	"time"

	"github.com/I1820/backend/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/envy"
)

// NewAccessToken creates new access token for given user.
func NewAccessToken(u models.User) (string, error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		U:   u,                               // logged in user information
		Exp: time.Now().Add(time.Minute * 1), // tokens expire in 30 minutes
	})

	// Generate encoded token and send it as response
	encodedToken, err := token.SignedString([]byte(envy.Get("JWT_SECRET", "i1820")))
	if err != nil {
		return "", err
	}

	return encodedToken, nil
}

// NewRefreshToken creates new refresh token for given identification that does not expire until
// time.Now() + untilExpire.
func NewRefreshToken(id string, untilExpire time.Duration) (string, error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims{
		ID:  id,
		Exp: time.Now().Add(untilExpire),
	})

	// Generate encoded token and send it as response
	encodedToken, err := token.SignedString([]byte(envy.Get("JWT_SECRET", "i1820")))
	if err != nil {
		return "", err
	}

	return encodedToken, nil
}

// RefreshClaims contains required information in I1820 platform to
// refresh a user access token.
type RefreshClaims struct {
	ID  string
	Exp time.Time
}

// Valid checks claims expiration time
func (rc RefreshClaims) Valid() error {
	// if claims is expired, return with status unathorized
	if rc.Exp.Before(time.Now()) {
		return fmt.Errorf("Token in expired %g minutes ago", time.Now().Sub(rc.Exp).Minutes())
	}

	return nil
}

// UserClaims contains required information in I1820 platform
// for logged in user.
type UserClaims struct {
	U   models.User
	Exp time.Time
}

// Valid checks claims expiration time
func (uc UserClaims) Valid() error {
	// if claims is expired, return with status unathorized
	if uc.Exp.Before(time.Now()) {
		return fmt.Errorf("Token in expired %g minutes ago", time.Now().Sub(uc.Exp).Minutes())
	}

	return nil
}
