package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrSigningMethodMismatch is the error returned when token is signed with the method other than verified
	ErrSigningMethodMismatch = errors.New("signing method mismatch")
	// ErrFailedToParseToken is the error returned when token is failed to parse and validate against secret and expiration date
	ErrFailedToParseToken = errors.New("failed to parse token")
	// ErrUnsupportedSigningMethod is the error returned when token is signed with unsupported by the library method
	ErrUnsupportedSigningMethod = errors.New("unsupported signing method")
	// ErrInvalidPEMBlock is the error returned for keys expected to be PEM-encoded
	ErrInvalidPEMBlock = errors.New("invalid RSA: not PEM-encoded")
	// ErrNotRSAPublicKey is the error returned for invalid RSA public key
	ErrNotRSAPublicKey = errors.New("invalid RSA: expected PUBLIC KEY block type")
	// ErrBadPublicKey is the error returned for invalid RSA public key
	ErrBadPublicKey = errors.New("invalid RSA: failed to assert public key")
)

// SigningMethod defines signing method algorithm and key
type SigningMethod struct {
	// Alg defines JWT signing algorithm. Possible values are: HS256, HS384, HS512, RS256, RS384, RS512
	Alg string `json:"alg"`
	Key string `json:"key"`
}

// ParserConfig configures the way JWT Parser gets and validates token
type ParserConfig struct {
	// SigningMethods defines chain of token signature verification algorithm/key pairs.
	SigningMethods []SigningMethod

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// Leeway is the time in seconds to account for clock skew when checking nbf, iat or expiration times
	Leeway int64
}

// NewParserConfig creates a new instance of ParserConfig
func NewParserConfig(leeway int64, signingMethod ...SigningMethod) ParserConfig {
	return ParserConfig{
		SigningMethods: signingMethod,
		TokenLookup:    "header:Authorization",
		Leeway:         leeway,
	}
}

// Parser struct
type Parser struct {
	Config ParserConfig
}

// NewParser creates a new instance of Parser
func NewParser(config ParserConfig) *Parser {
	return &Parser{config}
}

// ParseFromRequest tries to extract and validate token from request.
// See "Guard.TokenLookup" for possible ways to pass token in request.
func (jp *Parser) ParseFromRequest(r *http.Request) (*jwt.Token, error) {
	var token string
	var err error

	parts := strings.Split(jp.Config.TokenLookup, ":")
	switch parts[0] {
	case "header":
		token, err = jp.jwtFromHeader(r, parts[1])
	case "query":
		token, err = jp.jwtFromQuery(r, parts[1])
	case "cookie":
		token, err = jp.jwtFromCookie(r, parts[1])
	}

	if err != nil {
		return nil, err
	}

	return jp.Parse(token)
}

// Parse a JWT token and validates it
func (jp *Parser) Parse(tokenString string) (*jwt.Token, error) {
	for _, method := range jp.Config.SigningMethods {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != method.Alg {
				return nil, ErrSigningMethodMismatch
			}

			switch token.Method.(type) {
			case *jwt.SigningMethodHMAC:
				return []byte(method.Key), nil
			case *jwt.SigningMethodRSA:
				block, _ := pem.Decode([]byte(method.Key))
				if block == nil {
					return nil, ErrInvalidPEMBlock
				}
				if got, want := block.Type, "PUBLIC KEY"; got != want {
					return nil, ErrNotRSAPublicKey
				}
				pub, err := x509.ParsePKIXPublicKey(block.Bytes)
				if nil != err {
					return nil, err
				}

				if _, ok := pub.(*rsa.PublicKey); !ok {
					return nil, ErrBadPublicKey
				}

				return pub, nil
			default:
				return nil, ErrUnsupportedSigningMethod
			}
		}, jwt.WithLeeway(time.Duration(jp.Config.Leeway)))

		if err != nil {
			if errors.Is(err, ErrSigningMethodMismatch) {
				continue
			}

			if errors.Is(err, jwt.ErrTokenUnverifiable) || errors.Is(err, jwt.ErrSignatureInvalid) {
				continue
			}
		}

		return token, err
	}

	return nil, ErrFailedToParseToken
}

// GetMapClaims returns a map version of Claims Section
func (jp *Parser) GetMapClaims(token *jwt.Token) (jwt.MapClaims, bool) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}, ok
	}
	return claims, ok
}

func (jp *Parser) jwtFromHeader(r *http.Request, key string) (string, error) {
	authHeader := r.Header.Get(key)

	if authHeader == "" {
		return "", errors.New("auth header empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("invalid auth header")
	}

	return parts[1], nil
}

func (jp *Parser) jwtFromQuery(r *http.Request, key string) (string, error) {
	token := r.URL.Query().Get(key)

	if token == "" {
		return "", errors.New("query token empty")
	}

	return token, nil
}

func (jp *Parser) jwtFromCookie(r *http.Request, key string) (string, error) {
	cookie, _ := r.Cookie(key)

	if nil == cookie {
		return "", errors.New("cookie token empty")
	}

	return cookie.Value, nil
}
