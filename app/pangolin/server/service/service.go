package service

import (
	"pangolin/app/pangolin/config"
	"pangolin/app/pangolin/model/db"
	"pangolin/app/pangolin/server/dal"
	"pangolin/app/pangolin/utils"
	"sync"
)

type Service struct {
	currentConfig *config.DBConfig
	dbAccess      *dal.DBAccess
	dbLock        *sync.RWMutex
	dbConnected   bool
	tunnelMapLock *sync.RWMutex
	tunnelMap     map[uint64]*utils.SSHTunnel
	portMap       map[int]bool
	tunnelConfig  map[uint64]*db.TunnelConfig
}

func (srv *Service) InitResources() {
	srv.dbLock = &sync.RWMutex{}
	srv.tunnelMap = map[uint64]*utils.SSHTunnel{}
	srv.portMap = map[int]bool{}
	srv.tunnelConfig = map[uint64]*db.TunnelConfig{}
	srv.tunnelMapLock = &sync.RWMutex{}
}
