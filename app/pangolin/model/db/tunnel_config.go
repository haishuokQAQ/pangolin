package db

import (
	"time"
)

const (
	LogModePassword   = "password"
	LogModePrivateKey = "private_key"
)

type TunnelConfig struct {
	Id         uint64     `json:"id"`
	LocalHost  string     `json:"local_host"`
	LocalPort  int        `json:"local_port"`
	ServerHost string     `json:"server_host"`
	ServerPort int        `json:"server_port"`
	RemoteHost string     `json:"remote_host"`
	RemotePort int        `json:"remote_port"`
	LogMode    string     `json:"log_mode"`
	UserName   string     `json:"user_name"`
	Password   string     `json:"password"`
	PrivateKey string     `json:"private_key"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
