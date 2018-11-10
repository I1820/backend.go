/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 31-08-2018
 * |
 * | File Name:     user.go
 * +===============================================
 */

package models

// User represents users that signup into I1820 platform.
// Each user can have many projects. projects can exist for multiple users at the same time.
// I1820 do not store tokens in database, tokens just check by hand in auth middleware.
type User struct {
	Username  string `json:"username" bson:"username"`   // like 1995parham without any restriction and must be unique in the system
	Firstname string `json:"firstname" bson:"firstname"` // UTF-8 support
	Lastname  string `json:"lastname" bson:"lastname"`   // UTF-8 support
	Email     string `json:"email" bson:"email"`         // valid email address
	Password  string `json:"password,omitempty" bson:"password"`

	AccessToken  string `json:"access_token,omitempty" bson:"-"`  // token is uesed to access I1820 resources
	RefreshToken string `json:"refresh_token,omitempty" bson:"-"` // token is used to generate a new token when the old one is expired.

	Projects []string `json:"projects" bson:"projects"` // list of user owned projects
}
