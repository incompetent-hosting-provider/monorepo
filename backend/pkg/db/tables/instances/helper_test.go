package db_instances

import "testing"

func TestInstanceIdGenerator(t *testing.T) {
	// ARRANGE
	userSub := "123"
	containerUUID := "456"

	// ACT
	res := getInstanceId(userSub, containerUUID)
	// Assert

	if res != "123456" {
		t.Fatalf(`Want %s, received %v`, "123456", res)
	}
}
