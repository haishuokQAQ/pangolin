package dal

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	ctxUtil "pangolin/app/pangolin/utils/context"
	gormUtil "pangolin/app/pangolin/utils/gorm"
	"pangolin/app/pangolin/utils/log"
	"pangolin/app/pangolin/utils/log/logger"
)

var stdDBAccess *DBAccess

type DBAccess struct {
	db     *gorm.DB
	Logger logger.Logger
}

func GetDB(ctx context.Context) *gorm.DB {
	var tx *gorm.DB
	txObj, exists := ctxUtil.GetTransaction(ctx)
	if exists {
		tx = txObj
	}
	if tx != nil {
		// tx 在初始化的时候已经注入过 trace logger
		return tx
	}
	//_, exists = ctxUtil.GetTrace(ctx)
	if exists {
		// 存在 trace 对象，注入trace logger
		// 克隆，不会克隆底层的 sqldb 对象，不用担心连接池问题
		db := stdDBAccess.db.New()
		db.SetLogger(gormUtil.NewGormLoggerWithLevel(stdDBAccess.Logger.WithTraceInCtx(ctx), log.LevelInfo))
		return db
	} else {
		// 没有 trace 对象，不处理
		return stdDBAccess.db
	}
}

func (da *DBAccess) DB() *gorm.DB {
	return da.db
}

func (da *DBAccess) BeginTransaction() *gorm.DB {
	return da.db.Begin()
}

// ConnectDB is used to open database connection
func ConnectDB(ip string, port int, username string, password string, dbname string) error {

	if stdDBAccess != nil {
		stdDBAccess.db.Close()
	}

	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, ip, port, dbname)

	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		return err
	}
	db.SingularTable(true)
	db.LogMode(false)
	db = db.Set("gorm:save_associations", false).Set("gorm:association_save_reference", false)
	db = db.Set("gorm:association_autoupdate", false)

	stdDBAccess = &DBAccess{
		db: db,
	}

	return nil
}

func CurrentDBAccess() *DBAccess {
	return stdDBAccess
}
