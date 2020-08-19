package service

import (
	"context"
	"google.golang.org/appengine/log"
	"pangolin/app/pangolin/model/db"
)

func (srv *Service) ListAllConfig(ctx context.Context, pageNum, pageSize int) ([]*db.TunnelConfig, int, error) {
	configs, totalCount, err := srv.dbAccess.ListAllConfig(pageNum, pageSize)
	if err != nil {
		log.Errorf(ctx, "Fail to list all config.Error %+v", err)
		return nil, 0, err
	}
	return configs, totalCount, nil
}

func (srv *Service) CreateConfig(ctx context.Context, config *db.TunnelConfig) error {
	err := srv.dbAccess.CreateConfig(config)
	if err != nil {
		log.Errorf(ctx, "Fail to create config.Error %+v", err)
		return err
	}
	return nil
}
