package db_instances

func getInstanceId(usersub string, containerUUID string) string {
	return usersub + containerUUID
}
