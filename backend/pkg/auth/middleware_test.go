package auth_test

import (
	"incompetent-hosting-provider/backend/pkg/auth"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func getExampleEngine() *gin.Engine {
	ginEngine := gin.New()
	ginEngine.GET("/", auth.AuthMiddleware, func(c *gin.Context) {
		c.String(200, "ok")
	})
	return ginEngine
}

func TestMiddleWareWithNoHeader(t *testing.T) {
	// ARRANGE
	g := getExampleEngine()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	// ACT
	g.ServeHTTP(w, req)
	// ARRANGE
	if w.Code != 401 {
		t.Fatalf(`Want %d, received %d`, 401, w.Code)
	}
}
