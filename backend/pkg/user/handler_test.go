package user_test

import (
	"encoding/json"
	"fmt"
	"incompetent-hosting-provider/backend/pkg/constants"
	"incompetent-hosting-provider/backend/pkg/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type userResponse struct {
	Balance int
	Email   string
}

func TestUserFetchHandlerAllHeadersPresent(t *testing.T) {
	// ARRANGE

	userEmail := "test@test.test"
	userId := "123abc456"

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request.Header.Set(constants.USER_EMAIL_HEADER, userEmail)
	c.Request.Header.Set(constants.USER_ID_HEADER, userId)
	// ACT
	user.UserFetchHandler(c)
	// Assert
	var response userResponse
	json.NewDecoder(w.Result().Body).Decode(&response)

	fmt.Printf("%v", response)

	if response.Email != userEmail {
		t.Fatalf(`Want %v, received %v`, userEmail, response.Email)
	}

	if response.Balance != 1000 {
		t.Fatalf(`Want %d, received %v`, 1000, response.Balance)
	}
}
