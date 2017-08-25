package basic

import (
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
)

var (
	// ErrNotAuthorized is used when the the access is not permisted
	ErrNotAuthorized = errors.New(http.StatusUnauthorized, "not authorized")
	// ErrUserNotFound is used when an user is not found
	ErrUserNotFound = errors.New(http.StatusNotFound, "user not found")
	// ErrUserExists is used when an user already exists
	ErrUserExists = errors.New(http.StatusNotFound, "user already exists")
	// ErrInvalidMongoDBSession is used when mongodb is not beeing used
	ErrInvalidMongoDBSession = errors.New(http.StatusNotFound, "invalid mongodb session given")
	// ErrInvalidAdminRouter is used when an invalid admin router is given
	ErrInvalidAdminRouter = errors.New(http.StatusNotFound, "invalid admin router given")
)
