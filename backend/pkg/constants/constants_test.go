package constants_test

import (
	"incompetent-hosting-provider/backend/pkg/constants"
	"testing"
)

func TestUserIdHeaderConstDefined(t *testing.T) {
	// ACT && ASSERT
	if constants.USER_ID_HEADER == "" {
		t.Fatalf(`USER_ID_HEADER const is not defined`)
	}
}

func TestUserEmailHeaderConstDefined(t *testing.T) {
	// ACT && ASSERT
	if constants.USER_EMAIL_HEADER == "" {
		t.Fatalf(`USER_ID_HEADER const is not defined`)
	}
}
