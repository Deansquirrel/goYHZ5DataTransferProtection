package repository

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

const (
	sqlGetConnStrByKey = "" +
		"SELECT [connStr] " +
		"FROM [connInfo] " +
		"WHERE connKey = ?"

	sqlGetCronByKey = "" +
		"SELECT [cron] " +
		"FROM [taskCron] " +
		"WHERE taskKey = ?"

	sqlUpdateHeartBeat = "" +
		"DELETE FROM [heartbeat] " +
		"WHERE 1=1 " +
		"INSERT INTO [heartbeat] " +
		"SELECT getDate()"

	sqlUpdateWaitRestoreDataCount = "" +
		"UPDATE WaitRestoreDataCount " +
		"SET TfRcv%d=?,getDataTime%d=? " +
		"WHERE BatchNo = ?"

	sqlAddWaitRestoreDataCountRecord = "" +
		"INSERT INTO WaitRestoreDataCount([batchNo],[tfRcv1],[getDataTime1],[tfRcv2],[getDataTime2]," +
		"[tfRcv3],[getDataTime3],[tfRcv4],[getDataTime4],[recordDate]) " +
		"SELECT ?,?,?,?,?,?,?,?,?,getDate()"

	sqlAddMdDataTransStateCreateTempTable = "" +
		"CREATE TABLE #mddatatransstate(" +
		"[batchno] [char](12) NOT NULL," +
		"[mdcode] [char](3) NOT NULL," +
		"[chanid] [tinyint] NOT NULL," +
		"[datrowcount] [int] NOT NULL," +
		"[datlstup] [datetime] NOT NULL," +
		"[datlstuprcv] [datetime] NOT NULL," +
		"[recorddate] [datetime] NOT NULL" +
		")"
	sqlAddMdDataTransStateInsertTempTable = "" +
		"INSERT INTO #mddatatransstate(" +
		"[batchno],[mdcode],[chanid],[datrowcount],[datlstup]," +
		"[datlstuprcv],[recorddate]" +
		") VALUES (?,?,?,?,?,?,?)"
	sqlAddMdDataTransStateInsertData = "" +
		"insert into mddatatransstate(" +
		"[batchno],[mdcode],[chanid],[datrowcount],[datlstup]," +
		"[datlstuprcv],[recorddate]) " +
		"select [batchno],[mdcode],[chanid],[datrowcount],[datlstup]," +
		"[datlstuprcv],[recorddate] " +
		"from #mddatatransstate"
	sqlUpdateMdDataTransState = "" +
		"UPDATE mddatatransstate " +
		"SET DATROWCOUNT=? " +
		"WHERE batchno=? AND MDCODE=? AND CHANID = ?"
	sqlInsertCwGsSet = "" +
		"INSERT INTO [cwgsset]([gsid],[gsname],[lastupdate]) " +
		"VALUES (?,?,GETDATE())"
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

//添加待恢复信道数据总量记录
func (c *config) AddWaitRestoreDataCountRecord(d *object.WaitRestoreDataCount) error {
	comm := common{}
	return comm.SetRowsBySQL(c.dbConfig,
		sqlAddWaitRestoreDataCountRecord,
		d.BatchNo, d.TfRcv1, d.GetDataTime1, d.TfRcv2, d.GetDataTime2, d.TfRcv3, d.GetDataTime3, d.TfRcv4, d.GetDataTime4)
}

//更新待恢复信道数据总量
func (c *config) UpdateWaitRestoreDataCount(batchNo string, chanId int, getDatetime time.Time, rowCount int) error {
	if batchNo == "" {
		return errors.New("batch no can not be empty")
	}
	if chanId < 0 || chanId > 4 {
		return errors.New(fmt.Sprintf("chan id must between 0 and 4,got %d", chanId))
	}
	if rowCount < 0 {
		return errors.New(fmt.Sprintf("row count must > 0,got %d", rowCount))
	}
	comm := common{}
	return comm.SetRowsBySQL(c.dbConfig,
		fmt.Sprintf(sqlUpdateWaitRestoreDataCount, chanId, chanId),
		rowCount, goToolCommon.GetDateTimeStr(getDatetime), batchNo)
}

//添加门店数据信道传递状态
func (c *config) AddMdDataTransState(dList []*object.MdDataTransState) error {
	db, err := goToolMSSql.GetConn(c.dbConfig)
	if err != nil {
		errMsg := fmt.Sprintf("get db conn err: %s", err.Error())
		log.Error(errMsg)
		return errors.New(errMsg)
	}
	tx, err := db.Begin()
	if err != nil {
		errMsg := fmt.Sprintf("get db tx err: %s", err.Error())
		log.Error(errMsg)
		return errors.New(errMsg)
	}
	_, err = tx.Exec(sqlAddMdDataTransStateCreateTempTable)
	if err != nil {
		errMsg := fmt.Sprintf("tx create temp table err: %s", err.Error())
		log.Error(errMsg)
		_ = tx.Rollback()
		return errors.New(errMsg)
	}
	for _, d := range dList {
		_, err = tx.Exec(sqlAddMdDataTransStateInsertTempTable,
			d.BatchNo, d.MdCode, d.ChanId, d.DataRowCount,
			goToolCommon.GetDateTimeStr(d.LastUp),
			goToolCommon.GetDateTimeStr(d.LastUpRcv),
			goToolCommon.GetDateTimeStr(d.RecordTime))
		if err != nil {
			errMsg := fmt.Sprintf("tx insert temp table err: %s", err.Error())
			log.Error(errMsg)
			_ = tx.Rollback()
			return errors.New(errMsg)
		}
	}
	_, err = tx.Exec(sqlAddMdDataTransStateInsertData)
	if err != nil {
		errMsg := fmt.Sprintf("tx insert data from temp table err: %s", err.Error())
		log.Error(errMsg)
		_ = tx.Rollback()
		return errors.New(errMsg)
	}
	err = tx.Commit()
	if err != nil {
		errMsg := fmt.Sprintf("tx commit err: %s", err.Error())
		log.Error(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

//更新门店数据信道暂存数据量
func (c *config) UpdateMdDataTransState(batchNo string, mdCode string, chanId int, rowCount int) error {
	comm := common{}
	return comm.SetRowsBySQL(c.dbConfig, sqlUpdateMdDataTransState, rowCount, batchNo, mdCode, chanId)
}

//根据操作更新财务公司设置
func (c *config) UpdateCwGsSet(opr *object.OprCwGsSet) error {
	switch opr.OprType {
	case 1:
		return c.cwGsSetInsert(opr)
	case 2:
		return c.cwGsSetUpdate(opr)
	case 3:
		return c.cwGsSetDelete(opr)
	default:
		errMsg := fmt.Sprintf("cwgs set opr type error,exp 1 or 2 or 3 got %d", opr.OprType)
		log.Error(errMsg)
		return errors.New(errMsg)
	}
}

//新增财务公司设置
func (c *config) cwGsSetInsert(opr *object.OprCwGsSet) error {
	comm := NewCommon()
	return comm.SetRowsBySQL(c.dbConfig, sqlInsertCwGsSet, opr.GsId, opr.GsName)
}

//更新财务公司设置
func (c *config) cwGsSetDelete(opr *object.OprCwGsSet) error {
	//TODO 更新财务公司设置
	return nil
}

//删除财务公司设置
func (c *config) cwGsSetUpdate(opr *object.OprCwGsSet) error {
	//TODO 删除财务公司设置
	return nil
}
