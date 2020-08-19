package service

import (
	"log"
	"pangolin/app/pangolin/config"
	"pangolin/app/pangolin/model/db"
	"pangolin/app/pangolin/server/dal"
	"pangolin/app/pangolin/utils"
	"sync"
)

type Service struct {
	currentConfig *config.DBConfig
	dbAccess      *dal.DBAccess
	logger        log.Logger
	tunnelMapLock *sync.RWMutex
	tunnelMap     map[uint64]*utils.SSHTunnel
	portMap       map[int]bool
	tunnelConfig  map[uint64]*db.TunnelConfig
}
