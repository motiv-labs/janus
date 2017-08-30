package provider

import (
	"net/http"

	"github.com/pkg/errors"
)

// Verifier contains the methods for verification of providers
type Verifier interface {
	Verify(r *http.Request, httpClient *http.Client) (bool, error)
}

// VerifierBasket acts as a collection of verifier
type VerifierBasket struct {
	verifiers []Verifier
}

//NewVerifierBasket creates a new instace of VerifierBasket
func NewVerifierBasket(verifiers ...Verifier) VerifierBasket {
	return VerifierBasket{verifiers: verifiers}
}

// Verify checks is the provider is valid
func (vb VerifierBasket) Verify(r *http.Request, httpClient *http.Client) (bool, error) {
	var wrappedErrors error

	for _, verifier := range vb.verifiers {
		var verified bool

		verified, err := verifier.Verify(r, httpClient)
		if err != nil {
			wrappedErrors = errors.Wrap(err, "verification failed")
			continue
		}
		if verified {
			return true, nil
		}
	}

	return false, wrappedErrors
}
