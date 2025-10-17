package config

type MysqlConf struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Charset  string
}

type DatabaseConf struct {
	Driver  string    `yaml:"driver"`
	Mysql   MysqlConf `yaml:"mysql"`
	Sqllite struct {
		Path string `yaml:"path"`
	} `yaml:"sqllite"`
}
