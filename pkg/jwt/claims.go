package jwt

import (
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
)

// JanusClaims is the temporary solution for JWT claims validation with leeway support,
// should be removed as soon as github.com/dgrijalva/jwt-go 4.0 will be released with
// leeway support out of the box. This code is loosely based on the solution from
// https://github.com/dgrijalva/jwt-go/issues/131
type JanusClaims struct {
	jwt.MapClaims

	leeway int64
}

// NewJanusClaims instantiates new JanusClaims
func NewJanusClaims(leeway int64) *JanusClaims {
	return &JanusClaims{MapClaims: jwt.MapClaims{}, leeway: leeway}
}

// UnmarshalJSON is Unmarshaler interface implementation for JanusClaims to unmarshal nested map claims correctly
func (c *JanusClaims) UnmarshalJSON(text []byte) error {
	return json.Unmarshal(text, &c.MapClaims)
}

// VerifyExpiresAt overrides jwt.StandardClaims.VerifyExpiresAt() to use leeway for check
func (c *JanusClaims) VerifyExpiresAt(cmp int64, req bool) bool {
	return c.MapClaims.VerifyExpiresAt(cmp-c.leeway, req)
}

// VerifyIssuedAt overrides jwt.StandardClaims.VerifyIssuedAt() to use leeway for check
func (c *JanusClaims) VerifyIssuedAt(cmp int64, req bool) bool {
	return c.MapClaims.VerifyIssuedAt(cmp+c.leeway, req)
}

// VerifyNotBefore overrides jwt.StandardClaims.VerifyNotBefore() to use leeway for check
func (c *JanusClaims) VerifyNotBefore(cmp int64, req bool) bool {
	return c.MapClaims.VerifyNotBefore(cmp+c.leeway, req)
}
