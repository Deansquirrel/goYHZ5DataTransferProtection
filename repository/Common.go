package repository

import (
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goToolMSSql2000"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/global"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
)

import log "github.com/Deansquirrel/goToolLog"

type common struct {
}

func NewCommon() *common {
	return &common{}
}

//获取配置库连接配置
func (c *common) getConfigDbConfig() *goToolMSSql.MSSqlConfig {
	return &goToolMSSql.MSSqlConfig{
		Server: global.SysConfig.DB.Server,
		Port:   global.SysConfig.DB.Port,
		DbName: global.SysConfig.DB.DbName,
		User:   global.SysConfig.DB.User,
		Pwd:    global.SysConfig.DB.Pwd,
	}
}

func (c *common) getTdDbConfig() (*goToolMSSql.MSSqlConfig, error) {
	//TODO 获取通道库连接配置
	log.Debug("获取通道库连接配置")
	return nil, nil
}

func (c *common) getZbDbConfig() (*goToolMSSql.MSSqlConfig, error) {
	//TODO 获取总部库连接配置
	log.Debug("获取总部库连接配置")
	return nil, nil
}

func (c *common) getDbConfigByStr(config string) (*goToolMSSql.MSSqlConfig, error) {
	//TODO 根据配置字符号串获取数据库连接配置
	log.Debug("根据配置字符号串获取数据库连接配置")
	return nil, nil
}

//将普通数据库连接配置转换为Sql2000可用的配置
func (c *common) convertDbConfigTo2000(dbConfig *goToolMSSql.MSSqlConfig) *goToolMSSql2000.MSSqlConfig {
	return &goToolMSSql2000.MSSqlConfig{
		Server: dbConfig.Server,
		Port:   dbConfig.Port,
		DbName: dbConfig.DbName,
		User:   dbConfig.User,
		Pwd:    dbConfig.Pwd,
	}
}

//获取任务执行时间配置
func (c *common) GetTaskCron() ([]*object.TaskCron, error) {
	// TODO 获取任务执行时间配置
	log.Debug("获取任务执行时间配置")
	return nil, nil
}

//根据Key获取任务执行时间配置
func (c *common) GetTaskCronByKey(key string) (*object.TaskCron, error) {
	// TODO 获取任务执行时间配置
	log.Debug("获取任务执行时间配置")
	return nil, nil
}
