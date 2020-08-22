package dal

import (
	"github.com/jinzhu/gorm"
	"pangolin/app/pangolin/model/db"
	"time"
)

func (dba *DBAccess) ListAllConfig(pageNum, pageSize int) ([]*db.TunnelConfig, int, error) {
	result := []*db.TunnelConfig{}
	totalCount := 0
	query := dba.db.Model(&db.TunnelConfig{})
	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}
	if pageNum > 0 && pageSize > 0 {
		query = query.Offset((pageNum - 1) * pageSize)
	}
	if pageSize > 0 {
		query = query.Limit(pageSize)
	}
	err = query.Find(&result).Error
	if err != nil {
		return nil, 0, err
	}
	return result, totalCount, nil
}

func (dba *DBAccess) CreateConfig(config *db.TunnelConfig) error {
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()
	return dba.db.Create(config).Error
}

func (dba *DBAccess) UpdateConfig(config *db.TunnelConfig) error {
	config.UpdatedAt = time.Now()
	return dba.db.Save(config).Error
}

func (dba *DBAccess) GetConfigById(id uint64) (*db.TunnelConfig, error) {
	result := &db.TunnelConfig{}
	err := dba.db.Model(result).Where("id = ?", id).First(result).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
