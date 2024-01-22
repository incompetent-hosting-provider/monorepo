package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"incompetent-hosting-provider/backend/pkg/auth"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const keyId = "my-key"

func getExampleEngine(raw *rsa.PrivateKey) (*gin.Engine, jwk.Key) {
	ginEngine := gin.New()

	if raw == nil {
		raw, _ = rsa.GenerateKey(rand.Reader, 2048)
	}

	key, _ := jwk.FromRaw(raw)

	key.Set(jwk.KeyIDKey, keyId)

	buf, err := json.MarshalIndent(key, "", "  ")

	jwkJson, _ := keyfunc.NewJWKJSON(buf)

	fmt.Printf("%v", key.Algorithm().String())

	// Construct artificial middleware because we do not ned to test http fetching from underlying library.
	// This is already pushing it in terms of only testing my code and not some random library...however there is no way I am smart enough to write a secure jwt validator
	a := auth.AuthMiddleware{JWKS: jwkJson}
	fmt.Printf("%v", err)
	ginEngine.GET("/", a.AuthFunc, func(c *gin.Context) {
		c.String(200, "ok")
	})
	return ginEngine, key
}

// Positive Tests
func TestWithValidJWT(t *testing.T) {
	// ARRANGE
	raw, _ := rsa.GenerateKey(rand.Reader, 2048)
	g, key := getExampleEngine(raw)
	w := httptest.NewRecorder()

	tok, _ := jwt.NewBuilder().
		Claim("sub", "abc-def").
		Claim("email", "test@test.tets").
		JwtID(keyId).
		IssuedAt(time.Now()).
		Build()

	s, err := jwt.Sign(tok, jwt.WithKey(jwa.RS512, key))

	fmt.Sprintln(s)

	fmt.Printf("%v", err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s))
	// ACT
	g.ServeHTTP(w, req)
	// ARRANGE
	if w.Code != 200 {
		t.Fatalf(`Want %d, received %d`, 200, w.Code)
	}
}

// Negative Tests
func TestMiddleWareWithNoHeader(t *testing.T) {
	// ARRANGE
	g, _ := getExampleEngine(nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	// ACT
	g.ServeHTTP(w, req)
	// ARRANGE
	if w.Code != 401 {
		t.Fatalf(`Want %d, received %d`, 401, w.Code)
	}
}

func TestWithMissingClaims(t *testing.T) {
	// ARRANGE
	raw, _ := rsa.GenerateKey(rand.Reader, 2048)
	g, key := getExampleEngine(raw)
	w := httptest.NewRecorder()

	tok, _ := jwt.NewBuilder().
		Claim("a", "abc-def").
		Claim("b", "test@test.tets").
		JwtID(keyId).
		IssuedAt(time.Now()).
		Build()

	s, err := jwt.Sign(tok, jwt.WithKey(jwa.RS512, key))

	fmt.Sprintln(s)

	fmt.Printf("%v", err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s))
	// ACT
	g.ServeHTTP(w, req)
	// ARRANGE
	if w.Code != 500 {
		t.Fatalf(`Want %d, received %d`, 500, w.Code)
	}
}

func TestWithInvalidJWT(t *testing.T) {
	// ARRANGE
	raw, _ := rsa.GenerateKey(rand.Reader, 2048)
	g, key := getExampleEngine(raw)
	w := httptest.NewRecorder()

	tok, _ := jwt.NewBuilder().
		Claim("sub", "abc-def").
		Claim("email", "test@test.tets").
		JwtID(keyId).
		IssuedAt(time.Now()).
		Build()

	// Different algorithm => invalid jwt
	s, err := jwt.Sign(tok, jwt.WithKey(jwa.ES256, key))

	fmt.Sprintln(s)

	// Add a random character to the jwt
	fmt.Printf("%va", err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s))
	// ACT
	g.ServeHTTP(w, req)
	// ARRANGE
	if w.Code != 401 {
		t.Fatalf(`Want %d, received %d`, 500, w.Code)
	}
}
