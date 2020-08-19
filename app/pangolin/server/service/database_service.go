package service

import (
	"pangolin/app/pangolin/config"
	"pangolin/app/pangolin/server/dal"
)

func (srv *Service) GetCurrentDBConfig() *config.DBConfig {
	return srv.currentConfig
}

func (srv *Service) ConnectDatabase(conf *config.DBConfig) error {
	err := dal.ConnectDB(conf.Host, conf.Port, conf.UserName, conf.Password, conf.DbName)
	if err != nil {
		return err
	}
	srv.currentConfig = conf
	return nil
}
