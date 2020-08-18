package dal

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

var stdDBAccess *DBAccess

type DBAccess struct {
	db     *gorm.DB
	Logger log.Logger
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
