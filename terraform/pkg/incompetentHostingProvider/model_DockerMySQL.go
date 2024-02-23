package incompetenthostingprovider

type DockerMySQL struct {
	index               int
	uid                 int
	external_port       int
	mysql_root_password string
}

// Getter
func (d DockerMySQL) GetUid() int {
	return d.uid
}

func (d DockerMySQL) GetMySqlRootPassword() string {
	return d.mysql_root_password
}

// Setter
func (d *DockerMySQL) SetUid(uid int) {
	d.uid = uid
}

func (d *DockerMySQL) SetMySqlRootPassword(mysql_root_password string) {
	d.mysql_root_password = mysql_root_password
}
