package repository

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goToolMSSql2000"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
	"strings"
	"sync"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

const (
	sqlGetChanWaitRestoreDataCount = "" +
		"SELECT count(*) AS num,Isnull(min(TrRcvTime),'1900-01-01') AS d " +
		"FROM TfRcv%d"
)

var rppZbSyncLock sync.Mutex

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

//获取待恢复信道数据总量
func (r *zb) GetWaitRestoreDataCount() (*object.WaitRestoreDataCount, error) {
	c := common{}
	return &object.WaitRestoreDataCount{
		BatchNo:      r.getBatchNo(),
		TfRcv1:       -1,
		GetDataTime1: c.GetDefaultOprTime(),
		TfRcv2:       -1,
		GetDataTime2: c.GetDefaultOprTime(),
		TfRcv3:       -1,
		GetDataTime3: c.GetDefaultOprTime(),
		TfRcv4:       -1,
		GetDataTime4: c.GetDefaultOprTime(),
	}, nil
}

//获取信道待恢复信道数据总量
func (r *zb) GetChanWaitRestoreDataCount(chanId int) (int, time.Time, error) {
	comm := common{}
	num := -1
	t := time.Now()

	rows, err := comm.GetRowsBySQL2000(r.dbConfig, fmt.Sprintf(sqlGetChanWaitRestoreDataCount, chanId))
	if err != nil {
		errMsg := fmt.Sprintf("get chan wait restore data count err: chanId %d %s", chanId, err.Error())
		log.Error(errMsg)
		return num, t, errors.New(errMsg)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		err = rows.Scan(&num, &t)
		if err != nil {
			errMsg := fmt.Sprintf("read chan wait restore data count err: chanId %d %s", chanId, err.Error())
			log.Error(errMsg)
			return num, t, errors.New(errMsg)
		}
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("read chan wait restore data count err: chanId %d %s", chanId, rows.Err().Error())
		log.Error(errMsg)
		return num, t, errors.New(errMsg)
	}
	return num, t, nil
}

//产生操作批次号
func (r *zb) getBatchNo() string {
	rppZbSyncLock.Lock()
	defer rppZbSyncLock.Unlock()
	bn := goToolCommon.GetDateTimeStr(time.Now())
	bn = strings.Replace(bn, " ", "", -1)
	bn = strings.Replace(bn, ":", "", -1)
	bn = strings.Replace(bn, "-", "", -1)
	bn = bn[2:]
	time.Sleep(time.Second)
	return bn
}

//获取财务公司设置操作首条记录（不存在则返回nil）
func (r *zb) GetLstCwGsSetOpr() (*object.OprCwGsSet, error) {
	//TODO 获取财务公司设置操作首条记录（不存在则返回nil）
	log.Debug("获取财务公司设置操作首条记录（不存在则返回nil）")
	return nil, nil
}

//删除财务公司设置操作记录
func (r *zb) DelCwGsSetOpr(sn int) error {
	//TODO 删除财务公司设置操作记录
	log.Debug("删除财务公司设置操作记录")
	return nil
}
