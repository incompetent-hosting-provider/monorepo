package db_instances

func getInstanceId(userSub string, containerUUID string) string {
	return userSub + containerUUID
}
