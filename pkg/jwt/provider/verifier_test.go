package provider

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	fakeVerifier1 *mockVerifier
	fakeVerifier2 *mockVerifier

	r              *http.Request
	verifierBasket Verifier
)

type mockVerifier struct {
	result1 bool
	result2 error
}

func (v *mockVerifier) Verify(r *http.Request) (bool, error) {
	return v.result1, v.result2
}

func (v *mockVerifier) VerifyReturns(result1 bool, result2 error) {
	v.result1 = result1
	v.result2 = result2
}

func TestAllVerifiersFailed(t *testing.T) {
	fakeVerifier1 = new(mockVerifier)
	fakeVerifier2 = new(mockVerifier)
	verifierBasket = NewVerifierBasket(fakeVerifier1, fakeVerifier2)

	fakeVerifier1.VerifyReturns(false, nil)
	fakeVerifier2.VerifyReturns(false, nil)
	result, err := verifierBasket.Verify(r)

	require.NoError(t, err)
	assert.False(t, result)
}

func TestOneVerifierFailed(t *testing.T) {
	fakeVerifier1 = new(mockVerifier)
	fakeVerifier2 = new(mockVerifier)
	verifierBasket = NewVerifierBasket(fakeVerifier1, fakeVerifier2)

	fakeVerifier1.VerifyReturns(false, nil)
	fakeVerifier2.VerifyReturns(true, nil)
	result, err := verifierBasket.Verify(r)

	require.NoError(t, err)
	assert.True(t, result)
}

func TestAllVerifiersError(t *testing.T) {
	fakeVerifier1 = new(mockVerifier)
	fakeVerifier2 = new(mockVerifier)
	verifierBasket = NewVerifierBasket(fakeVerifier1, fakeVerifier2)

	fakeVerifier1.VerifyReturns(false, errors.New("first error"))
	fakeVerifier2.VerifyReturns(false, errors.New("second error"))
	result, err := verifierBasket.Verify(r)

	require.Error(t, err)
	assert.False(t, result)
}

func TestOneVerifierError(t *testing.T) {
	fakeVerifier1 = new(mockVerifier)
	fakeVerifier2 = new(mockVerifier)
	verifierBasket = NewVerifierBasket(fakeVerifier1, fakeVerifier2)

	fakeVerifier1.VerifyReturns(false, errors.New("first error"))
	fakeVerifier2.VerifyReturns(false, nil)
	result, err := verifierBasket.Verify(r)

	require.Error(t, err)
	assert.False(t, result)
}

func TestOneVerifierReturnsTrueAndDoesNotError(t *testing.T) {
	fakeVerifier1 = new(mockVerifier)
	fakeVerifier2 = new(mockVerifier)
	verifierBasket = NewVerifierBasket(fakeVerifier1, fakeVerifier2)

	fakeVerifier1.VerifyReturns(false, errors.New("first error"))
	fakeVerifier2.VerifyReturns(true, nil)
	result, err := verifierBasket.Verify(r)

	require.NoError(t, err)
	assert.True(t, result)
}
