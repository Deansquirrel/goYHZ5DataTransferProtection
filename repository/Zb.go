package repository

import (
	"errors"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goToolMSSql2000"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
)

import log "github.com/Deansquirrel/goToolLog"

//门店信道数据
type zb struct {
	dbConfig       *goToolMSSql2000.MSSqlConfig
	configDbConfig *goToolMSSql.MSSqlConfig
}

func NewRepZb() (*zb, error) {
	repConfig, err := NewConfig()
	if err != nil {
		return nil, err
	}
	dbConfig, err := repConfig.GetDbConfig(object.ConnTypeKeyZb)
	if err != nil {
		return nil, err
	}
	if dbConfig == nil {
		errMsg := "td db config is nil"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	c := common{}
	configDbConfig := c.GetConfigDbConfig()
	if configDbConfig == nil {
		errMsg := "config db config is nil"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &zb{
		dbConfig:       c.ConvertDbConfigTo2000(dbConfig),
		configDbConfig: configDbConfig,
	}, nil
}

//获取门店营业日开闭店记录恢复时间（总部库）
func (r *zb) GetMdYyStateTransTime() (*object.MdYyStateTransTimeZb, error) {
	//TODO 获取门店营业日开闭店记录恢复时间（总部库）
	return nil, nil
}
