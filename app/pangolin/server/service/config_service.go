package service

import (
	"context"
	"errors"
	"fmt"
	"pangolin/app/pangolin/model/db"
)

func (srv *Service) ListAllConfig(ctx context.Context, pageNum, pageSize int) ([]*db.TunnelConfig, int, error) {
	srv.dbLock.RLock()
	defer srv.dbLock.RUnlock()
	configs, totalCount, err := srv.dbAccess.ListAllConfig(pageNum, pageSize)
	if err != nil {
		fmt.Println(err)
		//log.Errorf(ctx, "Fail to list all config.Error %+v", err)
		return nil, 0, err
	}
	return configs, totalCount, nil
}

func (srv *Service) CreateConfig(ctx context.Context, config *db.TunnelConfig) (*db.TunnelConfig, error) {
	srv.dbLock.RLock()
	defer srv.dbLock.RUnlock()
	err := srv.dbAccess.CreateConfig(config)
	if err != nil {
		fmt.Println(err)
		//log.Errorf(ctx, "Fail to create config.Error %+v", err)
		return nil, err
	}
	return config, nil
}

func (srv *Service) GetConfigById(ctx context.Context, id uint64) (*db.TunnelConfig, error) {
	config, err := srv.dbAccess.GetConfigById(id)
	if err != nil {
		fmt.Println(err)
		//log.Errorf(ctx, "Fail to get config by id.Error %+v", err)
		return nil, err
	}
	if config == nil {
		fmt.Println(err)
		//log.Errorf(ctx, "Cannot find config by id %+v", id)
		return nil, errors.New("Not found")
	}
	return config, nil
}
