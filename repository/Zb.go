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

	sqlGetLstCwGsSetOpr = "" +
		"SELECT TOP 1 [oprsn],[gsid],[gsname],[oprtype],[oprtime] " +
		"FROM [cwgsset] " +
		"ORDER BY [oprsn] ASC"

	sqlDelCwGsSetOpr = "" +
		"DELETE FROM [cwgsset] " +
		"WHERE [oprsn] = ?"

	sqlGetLstMdSetOpr = "" +
		"SELECT TOP 1 [oprsn] ,[mdid],[mdname],[mdcode],[oprtype],[oprtime] " +
		"FROM [mdset] " +
		"ORDER BY [oprsn] ASC"

	sqlDelMdSetOpr = "" +
		"DELETE FROM [mdset] " +
		"WHERE [oprsn] = ?"

	sqlGetLstMdCwGsRefOpr = "" +
		"SELECT TOP 1 [oprsn],[mdid],[gsid],[oprtype],[oprtime] " +
		"FROM [mdcwgsref] " +
		"ORDER BY [oprsn] ASC"
	sqlDelMdCwGsRefOpr = "" +
		"DELETE FROM [mdcwgsref] " +
		"WHERE [oprsn] = ?"

	sqlGetLstMdYyStateOpr = "" +
		"SELECT TOP 1 [oprsn],[mdid],[mdyydate],[mdyyopentime],[mdyyclosetime]," +
		"[mdyysjtype],[oprtype],[oprtime] " +
		"FROM [mdyystate]" +
		"ORDER BY [oprsn] ASC"
	sqlDelMdYyStateOpr = "" +
		"DELETE FROM [mdyystate] " +
		"WHERE [oprsn] = ?"
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
	comm := NewCommon()
	rows, err := comm.GetRowsBySQL2000(r.dbConfig, sqlGetLstCwGsSetOpr)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	rList := make([]*object.OprCwGsSet, 0)
	var oprSn, gsId, oprType int
	var gsName string
	var oprTime time.Time
	for rows.Next() {
		err = rows.Scan(&oprSn, &gsId, &gsName, &oprType, &oprTime)
		if err != nil {
			errMsg := fmt.Sprintf("get last cwgsset opr err: %s", err.Error())
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		rList = append(rList, &object.OprCwGsSet{
			OprSn:   oprSn,
			GsId:    gsId,
			GsName:  gsName,
			OprType: oprType,
			OprTime: oprTime,
		})
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("get last cwgsset opr err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(rList) > 1 {
		errMsg := fmt.Sprintf("get last cwgsset opr list count err: exp 1 act %d", len(rList))
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(rList) == 0 {
		return nil, nil
	}
	return rList[0], nil
}

//删除财务公司设置操作记录
func (r *zb) DelCwGsSetOpr(sn int) error {
	comm := NewCommon()
	return comm.SetRowsBySQL2000(r.dbConfig, sqlDelCwGsSetOpr, sn)
}

//获取门店设置操作首条记录（不存在则返回nil）
func (r *zb) GetLstMdSetOpr() (*object.OprMdSet, error) {
	comm := NewCommon()
	rows, err := comm.GetRowsBySQL2000(r.dbConfig, sqlGetLstMdSetOpr)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	rList := make([]*object.OprMdSet, 0)
	var oprSn, mdId, oprType int
	var mdName, mdCode string
	var oprTime time.Time
	for rows.Next() {
		err = rows.Scan(&oprSn, &mdId, &mdName, &mdCode, &oprType, &oprTime)
		if err != nil {
			errMsg := fmt.Sprintf("get last mdset opr err: %s", err.Error())
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		rList = append(rList, &object.OprMdSet{
			OprSn:   oprSn,
			MdId:    mdId,
			MdName:  mdName,
			MdCode:  mdCode,
			OprType: oprType,
			OprTime: oprTime,
		})
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("get last mdset opr err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(rList) > 1 {
		errMsg := fmt.Sprintf("get last mdset opr list count err: exp 1 act %d", len(rList))
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(rList) == 0 {
		return nil, nil
	}
	return rList[0], nil
}

//删除门店设置操作记录
func (r *zb) DelMdSetOpr(sn int) error {
	comm := NewCommon()
	return comm.SetRowsBySQL2000(r.dbConfig, sqlDelMdSetOpr, sn)
}

//获取门店财务公司关系操作首条记录（不存在则返回nil）
func (r *zb) GetLstMdCwGsRefOpr() (*object.OprMdCwGsRef, error) {
	comm := NewCommon()
	rows, err := comm.GetRowsBySQL2000(r.dbConfig, sqlGetLstMdCwGsRefOpr)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	rList := make([]*object.OprMdCwGsRef, 0)
	var oprSn, mdId, gsId, oprType int
	var oprTime time.Time
	for rows.Next() {
		err = rows.Scan(&oprSn, &mdId, &gsId, &oprType, &oprTime)
		if err != nil {
			errMsg := fmt.Sprintf("get last mdset opr err: %s", err.Error())
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		rList = append(rList, &object.OprMdCwGsRef{
			OprSn:   oprSn,
			MdId:    mdId,
			GsId:    gsId,
			OprType: oprType,
			OprTime: oprTime,
		})
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("get last mdcwgsref opr err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(rList) > 1 {
		errMsg := fmt.Sprintf("get last mdcwgsref opr list count err: exp 1 act %d", len(rList))
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(rList) == 0 {
		return nil, nil
	}
	return rList[0], nil
}

//删除门店财务公司关系操作记录
func (r *zb) DelMdCwGsRefOpr(sn int) error {
	comm := NewCommon()
	return comm.SetRowsBySQL2000(r.dbConfig, sqlDelMdCwGsRefOpr, sn)
}

//获取门店营业日开闭店操作首条记录（不存在则返回nil）
func (r *zb) GetLstMdYyStateOpr() (*object.OprMdYyState, error) {
	comm := NewCommon()
	rows, err := comm.GetRowsBySQL2000(r.dbConfig, sqlGetLstMdYyStateOpr)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	rList := make([]*object.OprMdYyState, 0)
	var oprSn, mdId, mdYySjType, oprType int
	var mdYyDate string
	var mdYyOpenTime, mdYyCloseTime, oprTime time.Time
	for rows.Next() {
		err = rows.Scan(&oprSn, &mdId, &mdYyDate, &mdYyOpenTime, &mdYyCloseTime, &mdYySjType, &oprType, &oprTime)
		if err != nil {
			errMsg := fmt.Sprintf("get last mdset opr err: %s", err.Error())
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		rList = append(rList, &object.OprMdYyState{
			OprSn:         oprSn,
			MdId:          mdId,
			MdYyDate:      mdYyDate,
			MdYyOpenTime:  mdYyOpenTime,
			MdYyCloseTime: mdYyCloseTime,
			MdYySjType:    mdYySjType,
			OprType:       oprType,
			OprTime:       oprTime,
		})
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("get last mdcwgsref opr err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(rList) > 1 {
		errMsg := fmt.Sprintf("get last mdcwgsref opr list count err: exp 1 act %d", len(rList))
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(rList) == 0 {
		return nil, nil
	}
	return rList[0], nil
}

//删除门店营业日开闭店操作记录
func (r *zb) DelMdYyStateOpr(sn int) error {
	comm := NewCommon()
	return comm.SetRowsBySQL2000(r.dbConfig, sqlDelMdYyStateOpr, sn)
}
