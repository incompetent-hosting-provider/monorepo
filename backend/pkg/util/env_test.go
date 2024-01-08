package util_test

import (
	"incompetent-hosting-provider/backend/pkg/util"
	"os"
	"testing"
)

func TestGetStringEnvWithDefaultWithSetVariable(t *testing.T) {
	// ARRANGE
	envVarName := "TEST_RUN_ENV_UTIL"
	want := "testing"
	os.Setenv(envVarName, want)
	// Unset env var on test finish
	t.Cleanup(func() { os.Unsetenv(envVarName) })
	// ACT
	val := util.GetStringEnvWithDefault(envVarName, "invalid")
	// ASSERT
	if val != want {
		t.Fatalf(`Want %s, received %s`, want, val)
	}
}

func TestGetStringEnvWithDefaultWithUnsetVariable(t *testing.T) {
	// ARRANGE
	envVarName := "TEST_RUN_ENV_UTIL"
	want := "testing"
	// ACT
	val := util.GetStringEnvWithDefault(envVarName, want)
	if val != want {
		t.Fatalf(`Want %s, received %s`, want, val)
	}
}
