package basic

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// User represents a user that wants to login
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Equals compares one user to another
func (u *User) Equals(user *User) bool {
	return user.Username == u.Username && user.Password == u.Password
}

// PasswordVerifier checks if the current user `matches any of the given passwords
type PasswordVerifier struct {
	users []*User
}

// NewPasswordVerifier creates a new instance of PasswordVerifier
func NewPasswordVerifier(users []*User) *PasswordVerifier {
	return &PasswordVerifier{users}
}

// Verify makes a check and return a boolean if the check was successful or not
func (v *PasswordVerifier) Verify(r *http.Request, httpClient *http.Client) (bool, error) {
	var currentUser *User
	err := json.NewDecoder(r.Body).Decode(&currentUser)
	if err != nil {
		return false, errors.Wrap(err, "could not parse the json body")
	}

	for _, user := range v.users {
		if user.Equals(currentUser) {
			return true, nil
		}
	}

	log.WithFields(log.Fields{
		"have": currentUser,
		"want": v.users,
	}).Debug("not in the user list")

	return false, nil
}
