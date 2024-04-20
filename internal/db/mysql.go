package db

import (
	"gorm.io/gorm"
	tgorm "trpc.group/trpc-go/trpc-database/gorm"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/config"
)

// InitDatabase 初始化数据库并返回 gorm.DB 实例
func InitDatabase() (Db *gorm.DB) {
	var err error
	if Db, err = tgorm.NewClientProxy(config.DbServiceName); err != nil || Db == nil {
		if err != nil {
			log.Error(map[string]interface{}{"mysql connect error": err.Error()})
		} else {
			log.Error(map[string]interface{}{"database error": Db.Error})
		}
	}
	log.Info("Database initialization completed.")
	return
}
