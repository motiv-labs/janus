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
func (v *PasswordVerifier) Verify(r *http.Request) (bool, error) {
	currentUser, err := v.getUserFromRequest(r)
	if err != nil {
		log.Debug("Could not get user from request")
		return false, errors.Wrap(err, "could not get user from request")
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

func (v *PasswordVerifier) getUserFromRequest(r *http.Request) (*User, error) {
	var user *User

	//checks basic auth
	username, password, ok := r.BasicAuth()
	if ok {
		user = &User{
			Username: username,
			Password: password,
		}
	}

	// checks if the content is json otherwise just get from the form params
	if r.Header.Get("Content-Type") == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			return user, errors.Wrap(err, "could not parse the json body")
		}
	} else {
		r.ParseForm()

		user = &User{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}
	}

	return user, nil
}
