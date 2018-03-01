package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	basejwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseFromRequest_jwtFromHeader(t *testing.T) {
	alg := "HS256"
	key := time.Now().Format(time.RFC3339Nano)

	tokenString, err := generateToken(alg, key)
	require.NoError(t, err)

	req := &http.Request{Header: http.Header{}}

	config := NewParserConfig(SigningMethod{Alg: alg, Key: key})
	parser := NewParser(config)

	_, err = parser.ParseFromRequest(req)
	assert.Error(t, err)

	req.Header.Set("Authorization", "Basic "+tokenString)
	_, err = parser.ParseFromRequest(req)
	assert.Error(t, err)

	req.Header.Set("Authorization", "Bearer "+tokenString)

	assertParseToken(t, parser, req)
}

func TestParser_ParseFromRequest_jwtFromQuery(t *testing.T) {
	alg := "HS256"
	key := time.Now().Format(time.RFC3339Nano)

	tokenString, err := generateToken(alg, key)
	require.NoError(t, err)

	config := NewParserConfig(SigningMethod{Alg: alg, Key: key})
	config.TokenLookup = "query:token"
	parser := NewParser(config)

	req := &http.Request{URL: &url.URL{}}

	_, err = parser.ParseFromRequest(req)
	assert.Error(t, err)

	req.URL.RawQuery = "asd=qwe&token=" + tokenString

	assertParseToken(t, parser, req)
}

func TestParser_ParseFromRequest_jwtFromCookie(t *testing.T) {
	alg := "HS256"
	key := time.Now().Format(time.RFC3339Nano)

	tokenString, err := generateToken(alg, key)
	require.NoError(t, err)

	config := NewParserConfig(SigningMethod{Alg: alg, Key: key})
	config.TokenLookup = "cookie:token"
	parser := NewParser(config)

	req := &http.Request{Header: http.Header{}}

	_, err = parser.ParseFromRequest(req)
	assert.Error(t, err)

	req.Header.Set("Cookie", "qwe=asd;token="+tokenString)

	assertParseToken(t, parser, req)
}

func TestParser_Parse(t *testing.T) {
	alg := "RS256"

	tokenString, err := generateToken(alg, rsa2048Private)
	require.NoError(t, err)

	config := NewParserConfig(SigningMethod{Alg: "HS256", Key: time.Now().Format(time.RFC3339Nano)})
	parser := NewParser(config)

	req := &http.Request{Header: http.Header{"Authorization": {"Bearer " + tokenString}}}

	_, err = parser.ParseFromRequest(req)
	require.Error(t, err)

	parser.Config.SigningMethods = append(parser.Config.SigningMethods, SigningMethod{Alg: alg, Key: rsa2048Public})

	assertParseToken(t, parser, req)
}

func TestParser_Parse_ErrInvalidPEMBlock(t *testing.T) {
	alg := "RS256"

	tokenString, err := generateToken(alg, rsa2048Private)
	require.NoError(t, err)

	config := NewParserConfig(SigningMethod{Alg: alg, Key: "invalid public key"})
	parser := NewParser(config)

	req := &http.Request{Header: http.Header{"Authorization": {"Bearer " + tokenString}}}

	_, err = parser.ParseFromRequest(req)
	assert.Error(t, err)
}

func TestParser_Parse_ErrNotRSAPublicKey(t *testing.T) {
	alg := "RS256"

	tokenString, err := generateToken(alg, rsa2048Private)
	require.NoError(t, err)

	config := NewParserConfig(SigningMethod{Alg: alg, Key: rsa2048Private})
	parser := NewParser(config)

	req := &http.Request{Header: http.Header{"Authorization": {"Bearer " + tokenString}}}

	_, err = parser.ParseFromRequest(req)
	assert.Error(t, err)
}

func TestParser_Parse_ParsePKIXPublicKey(t *testing.T) {
	alg := "RS256"

	tokenString, err := generateToken(alg, rsa2048Private)
	require.NoError(t, err)

	config := NewParserConfig(SigningMethod{Alg: alg, Key: `-----BEGIN PUBLIC KEY-----
AIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvHA+KjSHPzVp7HIVGrQv
xUNnYrfdUR5+MRn3SM/Ts7GwwfifIlTgqjRHjm+sPOhVauF+ZX5PkUmW/HBlOxsj
zA55mCBFymO+nyjl/DhhFNLnKu3IWL8Q3IsnM1EgE9FHgwFqf3X5Eh5h2NdsOWPk
FrOB6BKY1wkOI2E27bNJAat+F059HU9z+jgIwhcm/IciTbqts497x1can4+NFBOl
VWE+yii6tREHF8olVe9a1DA8k7mtOQ2+1bK69kxm7tIde5sWnrlG3dv0gvlF25b3
XhiFAMJ1RfgvQHjXbtmMaZdMxj/Kx4CvBM37eXBGUSDWt9q97g+ywIQ/NrPZWmHo
ewIDAQAB
-----END PUBLIC KEY-----`})
	parser := NewParser(config)

	req := &http.Request{Header: http.Header{"Authorization": {"Bearer " + tokenString}}}

	_, err = parser.ParseFromRequest(req)
	assert.Error(t, err)
}

func TestParser_Parse_ErrBadPublicKey(t *testing.T) {
	alg := "RS256"

	tokenString, err := generateToken(alg, rsa2048Private)
	require.NoError(t, err)

	// DSA public key
	// ssh-keygen -t dsa
	// openssl dsa -in ~/.ssh/id_dsa -outform pem > dsa_priv.pem
	// openssl dsa -in dsa_priv.pem -pubout -out dsa_pub.pem
	config := NewParserConfig(SigningMethod{Alg: alg, Key: `-----BEGIN PUBLIC KEY-----
MIIBtjCCASsGByqGSM44BAEwggEeAoGBAIATdqQyUyprc8NtzuttJz8JahT+vwVK
4d2eVufm5IHyuqUyYroPQYpjQ1AfHOTE3ntmNFJcF+KhyqCTVdnaWwrmfiH2H6+D
2b+O8J50QFONKgktxy5LSBpUIIJfbhJG6vWW5GETnKOC7unoMh7yWDkDBYx+sSdg
ePBlI+Lq2+V1AhUA5e2ydAoXe0xa2lmoQDq09s+YsJUCgYBAMk28rTCuPLw+a7LF
++ouIDVMfxc/r8+/L4RCX+B5IScsk8SyzPeYdFtnCSGklQMdMw6YXCPdHGcexK/F
F7i9t5vxpD98aWRrJBW5fE99CPUXuMO6Gn8kV+1flRoBeBjPCd807BH/VEgcPGB4
ipAeCcl1yxfxyM6xARg4Fm/L7AOBhAACgYBESpculbUlOxvLK8tnYNI55T3eKGXw
9oSpxhgEzczq98PhaDu+ajjOqdD7DrM/VyvQuOwvhChPDTOlhazRZwyPCX1lUWnY
gXWfkeyb1H3jXz0cOQe2iHCSSMSYr2sH/E7kOMknJClemVFjWC7KO1F1yFXAspPs
V2pT9Twi0IeXmw==
-----END PUBLIC KEY-----`})
	parser := NewParser(config)

	req := &http.Request{Header: http.Header{"Authorization": {"Bearer " + tokenString}}}

	_, err = parser.ParseFromRequest(req)
	assert.Error(t, err)
}

func TestParser_Parse_ErrUnsupportedSigningMethod(t *testing.T) {
	alg := "PS256"

	tokenString, err := generateToken(alg, rsa2048Private)
	require.NoError(t, err)

	config := NewParserConfig(SigningMethod{Alg: alg, Key: rsa2048Public})
	parser := NewParser(config)

	req := &http.Request{Header: http.Header{"Authorization": {"Bearer " + tokenString}}}

	_, err = parser.ParseFromRequest(req)
	assert.Error(t, err)
}

const (
	clientID = "test-client-id"
	userName = "test@hellofresh.com"

	rsa2048Private = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAvHA+KjSHPzVp7HIVGrQvxUNnYrfdUR5+MRn3SM/Ts7Gwwfif
IlTgqjRHjm+sPOhVauF+ZX5PkUmW/HBlOxsjzA55mCBFymO+nyjl/DhhFNLnKu3I
WL8Q3IsnM1EgE9FHgwFqf3X5Eh5h2NdsOWPkFrOB6BKY1wkOI2E27bNJAat+F059
HU9z+jgIwhcm/IciTbqts497x1can4+NFBOlVWE+yii6tREHF8olVe9a1DA8k7mt
OQ2+1bK69kxm7tIde5sWnrlG3dv0gvlF25b3XhiFAMJ1RfgvQHjXbtmMaZdMxj/K
x4CvBM37eXBGUSDWt9q97g+ywIQ/NrPZWmHoewIDAQABAoIBAFLFbtj1F89Q9AUT
G2gOa8lXUStQnhtKrJ1+zVsjRtdwnralMalP5Rt+OUw8i0h5uUNoZy/HqsWjsHmU
GTM8OZ4hYZHL4zwCUjHxMgx261XNShNWPSGWU568VOy6nr91tta5oYD5Xf1ycQJh
pb0TvpWmJdK9kHssFBTAV/NTRCdB3klSSQ0t9gIfsa7ILYylaQQMyEtO0u6mTDxT
JjAeIWhYrALU12gLQD4jndF9ouzzgut0mcFnQbNt9vhXTEC1ZRghlRL95ELlrSi0
8AMxgtaiMcIeRezDo4Y+SAAPkVzUprlGEts5TBWcP53/BPfo9Mf0WBBDKhjPsIcx
cFKjp5ECgYEA7z691lJPdj9A0xEMZia5ZcCu2yL8DCGFLDEyDAsCfts9RpIAobb5
X/jOvklwXSvgtvkaiZcfMbgeR3KallWYQFN/q28CX/KPFLm7iA+ON5/nEqgOYlTS
/dLb1JQUs2qfNPjpWAzVL4KLO+fZyUXVYo15uw+M/CFMqZxmUoh3LnMCgYEAyaKf
33RUYhO8vZj6oumAddAOVg3t4jqEJ7IkvrbXIyEPQT5P1DmJHWSmeca8b+Y8pvz9
hIeSuygqWLDS3U5y4MfBYFBsQLQqM0KntjOItW0G/1KqM9YkNBIG7Gfk7fGxH2f8
sOEIVA8V9i2HM62k7ZJ+9lxBFJ7BsCq5UBwzE9kCgYAzyVb6T3LX27VCesw+SF+V
QPIYiSgZ0B+tgzCcHr35i6dl4TC10I+GUKsf0XG7GUZZFO7Dnayo7HvRZ2NC62A7
fFeEWlEfR7fk+pc3Sna0X657AVmru0S4oK3pA+y/MXMo2kBYSN7Um+NbokIoKS+Z
V5pj/We9I9AeXrZfYx65NQKBgGd7yjdhuckYPh7Ee6XO1zofzKvHvFYGGDtTR16F
8kY6Ol0OwOO3n7JxLKuFHsMDVA+T+fzho6HgTFN2dNJV58mLW6i1vck7bgke5Xoy
WrBaQ2QYpfeyqKP8uIbuD2U7TN9EfEC/TYnusCPHXANe1C2FqRmBYXlWvStP0gnW
XzSJAoGBALVfBV/WXcTArvbsT2KvwJovZG9kmSiR3ba3iXIeGwvBxtuDyeodz5k8
9S3m+ev58TCf1lYad+FAavr1ro8fbFyZZV6HItz4v377VcljxlvN739ST6R1RY36
PX7Lwn6YrQ2gk9efgKAEmcxBenq6UKkNXEGeiv4vxGG/cuTepQme
-----END RSA PRIVATE KEY-----`
	rsa2048Public = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvHA+KjSHPzVp7HIVGrQv
xUNnYrfdUR5+MRn3SM/Ts7GwwfifIlTgqjRHjm+sPOhVauF+ZX5PkUmW/HBlOxsj
zA55mCBFymO+nyjl/DhhFNLnKu3IWL8Q3IsnM1EgE9FHgwFqf3X5Eh5h2NdsOWPk
FrOB6BKY1wkOI2E27bNJAat+F059HU9z+jgIwhcm/IciTbqts497x1can4+NFBOl
VWE+yii6tREHF8olVe9a1DA8k7mtOQ2+1bK69kxm7tIde5sWnrlG3dv0gvlF25b3
XhiFAMJ1RfgvQHjXbtmMaZdMxj/Kx4CvBM37eXBGUSDWt9q97g+ywIQ/NrPZWmHo
ewIDAQAB
-----END PUBLIC KEY-----`
)

func assertParseToken(t *testing.T, parser *Parser, r *http.Request) {
	token, err := parser.ParseFromRequest(r)
	require.NoError(t, err)

	claims, ok := parser.GetMapClaims(token)
	assert.True(t, ok)
	assert.Equal(t, clientID, claims["iss"])

	customClaims, ok := parser.GetMapClaims(token)
	assert.True(t, ok)
	assert.Equal(t, userName, customClaims["username"])
}

func generateToken(alg, key string) (string, error) {
	type userClaims struct {
		Username string `json:"username"`
		basejwt.StandardClaims
	}

	token := basejwt.NewWithClaims(basejwt.GetSigningMethod(alg), userClaims{
		userName,
		basejwt.StandardClaims{
			Issuer:    clientID,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	})

	var signingKey interface{}
	switch token.Method.(type) {
	case *basejwt.SigningMethodHMAC:
		signingKey = []byte(key)
	case *basejwt.SigningMethodRSA, *basejwt.SigningMethodRSAPSS:
		block, _ := pem.Decode([]byte(key))
		if block == nil {
			return "", ErrInvalidPEMBlock
		}
		if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
			return "", errors.New("Invalid RSA: expected RSA PRIVATE KEY block type")
		}

		var err error
		signingKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if nil != err {
			return "", err
		}
	default:
		return "", ErrUnsupportedSigningMethod
	}

	return token.SignedString(signingKey)
}
