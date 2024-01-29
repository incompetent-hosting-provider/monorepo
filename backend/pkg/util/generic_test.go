package util_test

import (
	"incompetent-hosting-provider/backend/pkg/util"
	"testing"
)

func TestIsTestRunInTest(t *testing.T) {
	// ACT && ASSERT
	if !util.IsTestRun() {
		t.Fatal(`IsTestRun did not realize this was a test run`)
	}
}
