package service

import (
	"pangolin/app/pangolin/config"
	"pangolin/app/pangolin/server/dal"
)

func (srv *Service) GetCurrentDBConfig() *config.DBConfig {
	if srv.currentConfig == nil {
		return &config.DBConfig{}
	}
	return srv.currentConfig
}

func (srv *Service) ConnectDatabase(conf *config.DBConfig) error {
	srv.dbLock.Lock()
	defer srv.dbLock.Unlock()
	err := dal.ConnectDB(conf.Host, conf.Port, conf.UserName, conf.Password, conf.DbName)
	if err != nil {
		return err
	}
	srv.dbAccess = dal.CurrentDBAccess()
	srv.currentConfig = conf
	return nil
}
