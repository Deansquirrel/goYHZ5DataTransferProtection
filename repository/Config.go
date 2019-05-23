package repository

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
)

import log "github.com/Deansquirrel/goToolLog"

const (
	sqlGetConnStrByKey = "" +
		"SELECT [connstr] " +
		"FROM [conninfo] " +
		"WHERE connKey = ?"

	sqlGetCronByKey = "" +
		"SELECT [cron] " +
		"FROM [taskcron] " +
		"WHERE taskkey = ?"

	sqlUpdateHeartBeat = "" +
		"DELETE FROM [heartbeat] " +
		"WHERE 1=1 " +
		"INSERT INTO [heartbeat] " +
		"SELECT getdate()"
)

type config struct {
	dbConfig *goToolMSSql.MSSqlConfig
}

func NewConfig() (*config, error) {
	c := common{}
	return &config{
		dbConfig: c.GetConfigDbConfig(),
	}, nil
}

func (c *config) UpdateHeartBeat() error {
	comm := common{}
	err := comm.SetRowsBySQL(c.dbConfig, sqlUpdateHeartBeat)
	if err != nil {
		errMsg := fmt.Sprintf("update heart beat err: %s", err.Error())
		log.Error(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

//获取通道库连接配置
func (c *config) GetDbConfig(key object.ConnTypeKey) (*goToolMSSql.MSSqlConfig, error) {
	log.Debug(fmt.Sprintf("获取通道库连接配置[%s]", key))
	comm := common{}
	rows, err := comm.GetRowsBySQL(comm.GetConfigDbConfig(), sqlGetConnStrByKey, key)
	if err != nil {
		errLs := errors.New(fmt.Sprintf("get db config err,key = %s,err = %s", key, err.Error()))
		log.Error(errLs.Error())
		return nil, errLs
	}
	defer func() {
		_ = rows.Close()
	}()
	connStr := make([]string, 0)
	var str string
	for rows.Next() {
		err = rows.Scan(&str)
		if err != nil {
			errLs := errors.New(fmt.Sprintf("read db config val,key = %s,err = %s", key, err.Error()))
			log.Error(errLs.Error())
			return nil, errLs
		} else {
			connStr = append(connStr, str)
		}
	}
	if rows.Err() != nil {
		errLs := errors.New(fmt.Sprintf("read db config val,key = %s,err = %s", key, rows.Err().Error()))
		log.Error(errLs.Error())
		return nil, errLs
	}
	if len(connStr) != 1 {
		errLs := errors.New(fmt.Sprintf("db config num error,exp 1,act %d", len(connStr)))
		log.Error(errLs.Error())
		return nil, errLs
	}
	log.Debug(connStr[0])
	return comm.getDBConfigByStr(connStr[0])
}

//根据Key获取任务执行时间配置
func (c *config) GetTaskCronByKey(key object.TaskKey) (*object.TaskCron, error) {
	log.Debug(fmt.Sprintf("获取任务执行时间[%s]", key))
	comm := common{}
	rows, err := comm.GetRowsBySQL(comm.GetConfigDbConfig(), sqlGetCronByKey, key)
	if err != nil {
		errLs := errors.New(fmt.Sprintf("get task cron err,key = %s,err = %s", key, err.Error()))
		log.Error(errLs.Error())
		return nil, errLs
	}
	defer func() {
		_ = rows.Close()
	}()
	cronStr := make([]string, 0)
	var str string
	for rows.Next() {
		err = rows.Scan(&str)
		if err != nil {
			errLs := errors.New(fmt.Sprintf("get task cron err,key = %s,err = %s", key, err.Error()))
			log.Error(errLs.Error())
			return nil, errLs
		} else {
			cronStr = append(cronStr, str)
		}
	}
	if rows.Err() != nil {
		errLs := errors.New(fmt.Sprintf("get task cron err,key = %s,err = %s", key, rows.Err().Error()))
		log.Error(errLs.Error())
		return nil, errLs
	}
	if len(cronStr) != 1 {
		errLs := errors.New(fmt.Sprintf("task cron num error,exp 1,act %d", len(cronStr)))
		log.Error(errLs.Error())
		return nil, errLs
	}
	return &object.TaskCron{
		TaskKey: key,
		Cron:    cronStr[0],
	}, nil
}
