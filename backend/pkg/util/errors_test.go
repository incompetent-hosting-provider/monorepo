package util_test

import (
	"encoding/json"
	"incompetent-hosting-provider/backend/pkg/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestThrowInternalServerException(t *testing.T) {
	// ARRANGE
	g := gin.New()
	err_want := "test internal server error"
	code_want := 500
	g.GET("/", func(ctx *gin.Context) {
		util.ThrowInternalServerErrorException(ctx, err_want)
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	// ACT
	g.ServeHTTP(w, req)
	// ARRANGE
	if w.Code != code_want {
		t.Fatalf(`Want %d, received %d`, code_want, w.Code)
	}
	var target map[string]string
	_ = json.NewDecoder(w.Result().Body).Decode(&target)
	if target["error"] != err_want {
		t.Fatalf(`Want %d, received %d`, code_want, w.Code)
	}
}

func TestThrowUnauthorizedException(t *testing.T) {
	// ARRANGE
	g := gin.New()
	err_want := "test unauthorized exception"
	code_want := 401
	g.GET("/", func(ctx *gin.Context) {
		util.ThrowUnauthorizedException(ctx, err_want)
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	// ACT
	g.ServeHTTP(w, req)
	// ARRANGE
	if w.Code != code_want {
		t.Fatalf(`Want %d, received %d`, code_want, w.Code)
	}
	var target map[string]string
	_ = json.NewDecoder(w.Result().Body).Decode(&target)
	if target["error"] != err_want {
		t.Fatalf(`Want %d, received %d`, code_want, w.Code)
	}
}

func TestThrowNotFoundException(t *testing.T) {
	// ARRANGE
	g := gin.New()
	err_want := "test not found exception"
	code_want := 404
	g.GET("/", func(ctx *gin.Context) {
		util.ThrowNotFoundException(ctx, err_want)
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	// ACT
	g.ServeHTTP(w, req)
	// ARRANGE
	if w.Code != code_want {
		t.Fatalf(`Want %d, received %d`, code_want, w.Code)
	}
	var target map[string]string
	_ = json.NewDecoder(w.Result().Body).Decode(&target)
	if target["error"] != err_want {
		t.Fatalf(`Want %d, received %d`, code_want, w.Code)
	}
}
