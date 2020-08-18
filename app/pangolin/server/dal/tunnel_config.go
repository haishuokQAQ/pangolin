package dal

import (
	"context"
	"pangolin/app/pangolin/model/db"
)

func (dba *DBAccess) ListAllConfig(ctx context.Context, pageNum, pageSize int) ([]*db.TunnelConfig, error) {
	dbConn := GetDB(ctx)
	result := []*db.TunnelConfig{}
	query := dbConn.Model(&db.TunnelConfig{})
	if pageNum > 0 && pageSize > 0 {
		query = query.Offset((pageNum - 1) * pageSize)
	}
	if pageSize > 0 {
		query = query.Limit(pageSize)
	}
	err := query.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (dba *DBAccess) CreateConfig(ctx context.Context, config *db.TunnelConfig) error {
	dbConn := GetDB(ctx)
	return dbConn.Create(config).Error
}

func (dba *DBAccess) UpdateConfig(ctx context.Context, config *db.TunnelConfig) error {
	return GetDB(ctx).Save(config).Error
}
