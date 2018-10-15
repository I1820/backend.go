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
// each user can have many projects. projects can exist for multiple users at the same time.
type User struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password,omitempty" bson:"password"`
	Token    string `json:"token,omitempty" bson:"-"`
}
