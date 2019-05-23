package repository

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
	"strconv"
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
		"FROM %s"
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

//获取门店数据传递信道状态
func (r *td) GetMdDataTransState() ([]*object.MdDataTransState, error) {
	log.Debug("获取门店数据传递信道状态(rep)")
	c := common{}
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
	for rows.Next() {
		err = rows.Scan(&mdCode, &chanId, &datLstUp, &datLstUpRcv)
		if err != nil {
			errMsg := fmt.Sprintf("read data err: %s", err.Error())
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		state := object.MdDataTransState{
			RecordTime:   currTime,
			MdCode:       mdCode,
			ChanId:       chanId,
			DataRowCount: -1,
			LastUp:       datLstUp,
			LastUpRcv:    datLstUpRcv,
		}
		resultList = append(resultList, &state)
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("read data err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	var lastErr error
	for _, s := range resultList {
		tableName := fmt.Sprintf("zl" + object.TdUserCode + s.MdCode + "tfupv3_" + strconv.Itoa(s.ChanId))
		count, err := r.getDataRowCount(tableName)
		if err != nil {
			log.Error(err.Error())
			lastErr = err
		} else {
			s.DataRowCount = count
		}
	}
	return resultList, lastErr
}

//获取通道库表中的数据行数
func (r *td) getDataRowCount(tableName string) (int, error) {
	log.Debug(tableName)
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

//获取门店营业日开闭店记录传递时间（通道库）
func (r *td) GetMdYyStateTransTime() (*object.MdYyStateTransTimeTd, error) {
	//TODO 获取门店营业日开闭店记录传递时间（通道库）
	return nil, nil
}

//获取通道库连接配置
func (r *td) getDbConnConfig() (*goToolMSSql.MSSqlConfig, error) {
	repConfig, err := NewConfig()
	if err != nil {
		return nil, err
	}
	return repConfig.GetDbConfig(object.ConnTypeKeyTd)
}
