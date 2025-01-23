package db

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type DBConfig struct {
	Master      string
	Slaves      []string
	MaxIdle     int
	MaxOpen     int
	MaxLifetime time.Duration
}

func NewDB(config DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.Master), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 配置读写分离
	resolverConfig := dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(config.Master)},
		Replicas: make([]gorm.Dialector, len(config.Slaves)),
		Policy:   dbresolver.RandomPolicy{},
	}

	for i, slave := range config.Slaves {
		resolverConfig.Replicas[i] = mysql.Open(slave)
	}

	err = db.Use(dbresolver.Register(resolverConfig))
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(config.MaxIdle)
	sqlDB.SetMaxOpenConns(config.MaxOpen)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	return db, nil
}
