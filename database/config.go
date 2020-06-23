package database

type DBConfig struct {
	Hosts             string
	AuthDatabase      string
	Database          string
	AuthUserName      string
	AuthPassword      string
	ConnectionTimeout int64
	Env               string
	Replica           string
}
