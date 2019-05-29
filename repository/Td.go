package repository

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
	"strings"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

const (
	sqlGetMdDataTransState = "" +
		"SELECT [mdcode],[chanid],[datlstup],[datlstuprcv] " +
		"FROM [zltfmngdv3] " +
		"WHERE [usercode] = ? " +
		"	AND [chanid] < 5" +
		"Order by [mdcode] ASC,[chanid] ASC"
	sqlGetMdDataTransStateDataRowCount = "" +
		"SELECT COUNT(*) AS NUM " +
		"FROM %s WITH(NOLOCK)"
	sqlGetLstMdYyStateTransTimeTdOpr = "" +
		"SELECT TOP 1 [oprsn],[mdid],[mdyydate],[chanid],[oprtype],[oprtime] " +
		"FROM [mdyystatetranstime] " +
		"ORDER BY [oprsn] ASC"
	sqlDelMdYyStateTransTimeTdOpr = "" +
		"DELETE FROM [mdyystatetranstime] " +
		"WHERE [oprsn] = ?"
)

//门店信道数据
type td struct {
	dbConfig       *goToolMSSql.MSSqlConfig
	configDbConfig *goToolMSSql.MSSqlConfig
}

func NewRepTd() (*td, error) {
	repConfig, err := NewConfig()
	if err != nil {
		return nil, err
	}
	dbConfig, err := repConfig.GetDbConfig(object.ConnTypeKeyTd)
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
	return &td{
		dbConfig:       dbConfig,
		configDbConfig: configDbConfig,
	}, nil
}

//获取门店营业日开闭店记录传递首条记录（不存在则返回nil）
func (r *td) GetLstMdYyStateTransTimeTdOpr() (*object.OprMdYyStateTransTimeTd, error) {
	comm := NewCommon()
	rows, err := comm.GetRowsBySQL(r.dbConfig, sqlGetLstMdYyStateTransTimeTdOpr)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	rList := make([]*object.OprMdYyStateTransTimeTd, 0)
	var oprSn, mdId, chanId, oprType int
	var mdYyDate string
	var oprTime time.Time
	for rows.Next() {
		err = rows.Scan(&oprSn, &mdId, &mdYyDate, &chanId, &oprType, &oprTime)
		if err != nil {
			errMsg := fmt.Sprintf("get last MdYyStateTransTimeTdOpr opr err: %s", err.Error())
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		rList = append(rList, &object.OprMdYyStateTransTimeTd{
			OprSn:    oprSn,
			MdId:     mdId,
			MdYyDate: mdYyDate,
			ChanId:   chanId,
			OprType:  oprType,
			OprTime:  oprTime,
		})
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("get last MdYyStateTransTimeTdOpr opr err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(rList) > 1 {
		errMsg := fmt.Sprintf("get last MdYyStateTransTimeTdOpr opr list count err: exp 1 act %d", len(rList))
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(rList) == 0 {
		return nil, nil
	}
	return rList[0], nil
}

//删除门店财务公司关系操作记录
func (r *td) DelMdYyStateTransTimeTdOpr(sn int) error {
	comm := NewCommon()
	return comm.SetRowsBySQL(r.dbConfig, sqlDelMdYyStateTransTimeTdOpr, sn)
}

//获取门店数据传递信道状态
func (r *td) GetMdDataTransState() ([]*object.MdDataTransState, error) {
	log.Debug("获取门店数据传递信道状态(rep)")
	c := NewCommon()
	rows, err := c.GetRowsBySQL(r.dbConfig, sqlGetMdDataTransState, object.TdUserCode)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	currTime := time.Now()
	resultList := make([]*object.MdDataTransState, 0)
	var mdCode string
	var chanId int
	var datLstUp, datLstUpRcv time.Time
	batchNo := r.getBatchNo()
	for rows.Next() {
		err = rows.Scan(&mdCode, &chanId, &datLstUp, &datLstUpRcv)
		if err != nil {
			errMsg := fmt.Sprintf("read data err: %s", err.Error())
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		state := object.MdDataTransState{
			BatchNo:      batchNo,
			MdCode:       mdCode,
			ChanId:       chanId,
			DataRowCount: -1,
			LastUp:       datLstUp,
			LastUpRcv:    datLstUpRcv,
			RecordTime:   currTime,
		}
		resultList = append(resultList, &state)
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("read data err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return resultList, nil
}

//获取通道库表中的数据行数
func (r *td) GetDataRowCount(tableName string) (int, error) {
	c := common{}
	rows, err := c.GetRowsBySQL(r.dbConfig, fmt.Sprintf(sqlGetMdDataTransStateDataRowCount, tableName))
	if err != nil {
		errMsg := fmt.Sprintf("get data row count error,table %s,err: %s", tableName, err.Error())
		log.Error(errMsg)
		return -1, errors.New(errMsg)
	}
	defer func() {
		_ = rows.Close()
	}()
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			errMsg := fmt.Sprintf("get data row count,read data error,table %s,err: %s", tableName, err.Error())
			log.Error(errMsg)
			return -1, errors.New(errMsg)
		}
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("get data row count,read data error,table %s,err: %s", tableName, rows.Err().Error())
		log.Error(errMsg)
		return -1, errors.New(errMsg)
	}
	return count, nil
}

//获取通道库连接配置
func (r *td) getDbConnConfig() (*goToolMSSql.MSSqlConfig, error) {
	repConfig, err := NewConfig()
	if err != nil {
		return nil, err
	}
	return repConfig.GetDbConfig(object.ConnTypeKeyTd)
}

//产生操作批次号
func (r *td) getBatchNo() string {
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
