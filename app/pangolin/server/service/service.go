package service

import (
	"pangolin/app/pangolin/model/db"
	"pangolin/app/pangolin/utils"
	"pangolin/app/pangolin/utils/log/logger"
	"sync"
)

type Service struct {
	logger        logger.Logger
	tunnelMapLock *sync.RWMutex
	tunnelMap     map[uint64]*utils.SSHTunnel
	portMap       map[int]bool
	tunnelConfig  map[uint64]*db.TunnelConfig
}
