package provider

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	r          *http.Request
	httpClient *http.Client
)

type mockVerifier struct {
	result1 bool
	result2 error
}

func (v *mockVerifier) Verify(r *http.Request, httpClient *http.Client) (bool, error) {
	return v.result1, v.result2
}

func (v *mockVerifier) VerifyReturns(result1 bool, result2 error) {
	v.result1 = result1
	v.result2 = result2
}

func TestVerifiersScenarios(t *testing.T) {
	tests := []struct {
		scenario string
		function func(*testing.T, *mockVerifier, *mockVerifier, *VerifierBasket)
	}{
		{
			scenario: "when all verifiers fails, it should return false",
			function: testAllVerifiersFailed,
		},
		{
			scenario: "when one verifier fails, it should return false",
			function: testOneVerifierFailed,
		},
		{
			scenario: "when all verifiers have an error, it should return false and error",
			function: testAllVerifiersError,
		},
		{
			scenario: "when one verifier has an error, it should return false and error",
			function: testOneVerifierError,
		},
		{
			scenario: "when one verifier returns true and does not error, it should return true",
			function: testOneVerifierReturnsTrueAndDoesNotError,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			t.Parallel()
			fakeVerifier1 := new(mockVerifier)
			fakeVerifier2 := new(mockVerifier)
			verifierBasket := NewVerifierBasket(fakeVerifier1, fakeVerifier2)

			test.function(t, fakeVerifier1, fakeVerifier2, verifierBasket)
		})
	}
}

func testAllVerifiersFailed(t *testing.T, fakeVerifier1 *mockVerifier, fakeVerifier2 *mockVerifier, verifierBasket *VerifierBasket) {
	fakeVerifier1.VerifyReturns(false, nil)
	fakeVerifier2.VerifyReturns(false, nil)
	result, err := verifierBasket.Verify(r, httpClient)

	require.NoError(t, err)
	assert.False(t, result)
}

func testOneVerifierFailed(t *testing.T, fakeVerifier1 *mockVerifier, fakeVerifier2 *mockVerifier, verifierBasket *VerifierBasket) {
	fakeVerifier1.VerifyReturns(false, nil)
	fakeVerifier2.VerifyReturns(true, nil)
	result, err := verifierBasket.Verify(r, httpClient)

	require.NoError(t, err)
	assert.True(t, result)
}

func testAllVerifiersError(t *testing.T, fakeVerifier1 *mockVerifier, fakeVerifier2 *mockVerifier, verifierBasket *VerifierBasket) {
	fakeVerifier1.VerifyReturns(false, errors.New("first error"))
	fakeVerifier2.VerifyReturns(false, errors.New("second error"))
	result, err := verifierBasket.Verify(r, httpClient)

	require.Error(t, err)
	assert.False(t, result)
}

func testOneVerifierError(t *testing.T, fakeVerifier1 *mockVerifier, fakeVerifier2 *mockVerifier, verifierBasket *VerifierBasket) {
	fakeVerifier1.VerifyReturns(false, errors.New("first error"))
	fakeVerifier2.VerifyReturns(false, nil)
	result, err := verifierBasket.Verify(r, httpClient)

	require.Error(t, err)
	assert.False(t, result)
}

func testOneVerifierReturnsTrueAndDoesNotError(t *testing.T, fakeVerifier1 *mockVerifier, fakeVerifier2 *mockVerifier, verifierBasket *VerifierBasket) {
	fakeVerifier1.VerifyReturns(false, errors.New("first error"))
	fakeVerifier2.VerifyReturns(true, nil)
	result, err := verifierBasket.Verify(r, httpClient)

	require.NoError(t, err)
	assert.True(t, result)
}
