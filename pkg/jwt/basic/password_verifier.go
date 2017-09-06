package basic

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const (
	contentTypeJSON string = "application/json"
)

// user represents a user that wants to login
type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Equals compares one user to another
func (u *user) Equals(user *user) bool {
	return user.Username == u.Username && user.Password == u.Password
}

// PasswordVerifier checks if the current user `matches any of the given passwords
type PasswordVerifier struct {
	users []*user
}

// NewPasswordVerifier creates a new instance of PasswordVerifier
func NewPasswordVerifier(users []*user) *PasswordVerifier {
	return &PasswordVerifier{users}
}

// Verify makes a check and return a boolean if the check was successful or not
func (v *PasswordVerifier) Verify(r *http.Request, httpClient *http.Client) (bool, error) {
	currentUser, err := v.getUserFromRequest(r)
	if err != nil {
		return false, errors.Wrap(err, "could not get user from request")
	}

	for _, user := range v.users {
		if user.Equals(currentUser) {
			return true, nil
		}
	}

	return false, errors.New("incorrect username or password")
}

func (v *PasswordVerifier) getUserFromRequest(r *http.Request) (*user, error) {
	var u *user

	//checks basic auth
	username, password, ok := r.BasicAuth()
	if ok {
		u = &user{
			Username: username,
			Password: password,
		}
	}

	// checks if the content is json otherwise just get from the form params
	if r.Header.Get("Content-Type") == contentTypeJSON {
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			return u, errors.Wrap(err, "could not parse the json body")
		}
	} else {
		r.ParseForm()

		u = &user{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}
	}

	return u, nil
}
