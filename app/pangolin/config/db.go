package config

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
}
