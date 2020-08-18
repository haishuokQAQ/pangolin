package db

const (
	LogModePassword   = "password"
	LogModePrivateKey = "private_key"
)

type TunnelConfig struct {
	Id         uint64 `json:"id"`
	LocalHost  string `json:"local_host"`
	LocalPort  int    `json:"local_port"`
	ServerHost string `json:"server_host"`
	ServerPort int    `json:"server_port"`
	RemoteHost string `json:"remote_host"`
	RemotePort int    `json:"remote_port"`
	LogMode    string `json:"log_mode"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
	Deleted    int    `json:"deleted"`
	DeletedAt  int64  `json:"deleted_at"`
}
