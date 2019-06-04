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
	sqlUpdateCwGsSet = "" +
		"UPDATE [cwgsset] " +
		"SET [gsname] = ?,[lastupdate]=getdate() " +
		"WHERE [gsid] = ?"
	sqlDeleteCwGsSet = "" +
		"DELETE FROM [cwgsset] " +
		"WHERE [gsid] = ?"
	sqlGetCwGsSet = "" +
		"SELECT [gsid],[gsname],[lastupdate] " +
		"FROM [cwgsset]"

	sqlMdSetInsert = "" +
		"IF EXISTS(SELECT * FROM [mdset] WHERE [mdid] = ?) " +
		"BEGIN	" +
		"UPDATE [mdset] " +
		"SET [isforbidden] = 0,[lastupdate]=getDate() " +
		"WHERE [mdid] = ? " +
		"END " +
		"ELSE " +
		"BEGIN " +
		"INSERT INTO [mdset]([mdid],[mdname],[mdcode],[isforbidden],[lastupdate]) " +
		"VALUES (?,?,?,0,getDate()) " +
		"END"
	sqlMdSetDelete = "" +
		"UPDATE [mdset] " +
		"SET [isforbidden] = 1,[lastupdate]=getDate() " +
		"WHERE [mdid] = ?"
	sqlMdSetUpdate = "" +
		"UPDATE [mdset] " +
		"SET [mdname]=?,[mdcode]=?,[lastupdate]=getDate() " +
		"WHERE [mdid] = ?"

	sqlMdCwGsRefInsert = "" +
		"INSERT INTO [mdcwgsref]([mdid],[gsid],[lastupdate]) " +
		"VALUES (?,?,getDate())"
	sqlMdCwGsRefDelete = "" +
		"DELETE FROM [mdcwgsref] " +
		"WHERE [mdid] = ?"
	sqlMdCwGsRefUpdate = "" +
		"UPDATE [mdcwgsref] " +
		"SET [gsid]=?,[lastupdate]=getDate() " +
		"WHERE [mdid]=?"
	sqlUpdateMdYyStateTransTimeTd1 = "" +
		"IF EXISTS(SELECT * FROM [mdyystatetranstime] WHERE [mdid]=? AND [mdyydate]=? AND [chanid] = ?) " +
		"BEGIN " +
		"UPDATE [mdyystatetranstime] " +
		"SET [openup] = ?,[lastupdate]=getDate() " +
		"WHERE [mdid]=? AND [mdyydate]=? AND [chanid] = ? " +
		"END " +
		"ELSE " +
		"BEGIN " +
		"INSERT INTO [mdyystatetranstime]([mdid],[mdyydate],[chanid],[openup],[openuprcv]," +
		"[closeup],[closeuprcv],[lastupdate]) " +
		"VALUES (?,?,?,?,?,?,?,getDate()) " +
		"END "
	sqlUpdateMdYyStateTransTimeTd2 = "" +
		"IF EXISTS(SELECT * FROM [mdyystatetranstime] WHERE [mdid]=? AND [mdyydate]=? AND [chanid] = ?) " +
		"BEGIN " +
		"UPDATE [mdyystatetranstime] " +
		"SET [openuprcv] = ?,[lastupdate]=getDate() " +
		"WHERE [mdid]=? AND [mdyydate]=? AND [chanid] = ? " +
		"END " +
		"ELSE " +
		"BEGIN " +
		"INSERT INTO [mdyystatetranstime]([mdid],[mdyydate],[chanid],[openup],[openuprcv]," +
		"[closeup],[closeuprcv],[lastupdate]) " +
		"VALUES (?,?,?,?,?,?,?,getDate()) " +
		"END "
	sqlUpdateMdYyStateTransTimeTd3 = "" +
		"IF EXISTS(SELECT * FROM [mdyystatetranstime] WHERE [mdid]=? AND [mdyydate]=? AND [chanid] = ?) " +
		"BEGIN " +
		"UPDATE [mdyystatetranstime] " +
		"SET [closeup] = ?,[lastupdate]=getDate() " +
		"WHERE [mdid]=? AND [mdyydate]=? AND [chanid] = ? " +
		"END " +
		"ELSE " +
		"BEGIN " +
		"INSERT INTO [mdyystatetranstime]([mdid],[mdyydate],[chanid],[openup],[openuprcv]," +
		"[closeup],[closeuprcv],[lastupdate]) " +
		"VALUES (?,?,?,?,?,?,?,getDate()) " +
		"END "
	sqlUpdateMdYyStateTransTimeTd4 = "" +
		"IF EXISTS(SELECT * FROM [mdyystatetranstime] WHERE [mdid]=? AND [mdyydate]=? AND [chanid] = ?) " +
		"BEGIN " +
		"UPDATE [mdyystatetranstime] " +
		"SET [closeuprcv] = ?,[lastupdate]=getDate() " +
		"WHERE [mdid]=? AND [mdyydate]=? AND [chanid] = ? " +
		"END " +
		"ELSE " +
		"BEGIN " +
		"INSERT INTO [mdyystatetranstime]([mdid],[mdyydate],[chanid],[openup],[openuprcv]," +
		"[closeup],[closeuprcv],[lastupdate]) " +
		"VALUES (?,?,?,?,?,?,?,getDate()) " +
		"END "
	sqlUpdateMdYyState = "" +
		"IF EXISTS(SELECT * FROM [mdyystate] WHERE [mdid]=? AND [mdyydate]=?) " +
		"BEGIN " +
		"UPDATE [mdyystate] " +
		"SET [mdyyopentime] = ?,[mdyyclosetime] = ?,[mdyysjtype] = [mdyysjtype]|?,[lastupdate]=getDate() " +
		"WHERE [mdid]=? AND [mdyydate]=? " +
		"END " +
		"ELSE " +
		"BEGIN " +
		"INSERT INTO [mdyystate]([mdid],[mdyydate],[mdyyopentime],[mdyyclosetime],[mdyysjtype],[lastupdate]) " +
		"VALUES (?,?,?,?,?,getDate()) " +
		"END "

	sqlUpdateMdYyStateRestoreTime1 = "" +
		"IF EXISTS(SELECT * FROM [mdyystaterestoretime] WHERE [mdid]=? AND [mdyydate]=? AND [chanid]=?)  " +
		"BEGIN  " +
		"UPDATE [mdyystaterestoretime]  " +
		"SET [openuprcvr] = ?,[lastupdate]=getDate()  " +
		"WHERE [mdid]=? AND [mdyydate]=? AND [chanid]=? " +
		"END  " +
		"ELSE  " +
		"BEGIN  " +
		"INSERT INTO [mdyystaterestoretime]([mdid],[mdyydate],[chanid],[openuprcvr],[closeuprcvr],[lastupdate])  " +
		"VALUES (?,?,?,?,?,getDate())  " +
		"END "
	sqlUpdateMdYyStateRestoreTime2 = "" +
		"IF EXISTS(SELECT * FROM [mdyystaterestoretime] WHERE [mdid]=? AND [mdyydate]=? AND [chanid]=?)  " +
		"BEGIN  " +
		"UPDATE [mdyystaterestoretime]  " +
		"SET [closeuprcvr] = ?,[lastupdate]=getDate()  " +
		"WHERE [mdid]=? AND [mdyydate]=? AND [chanid]=? " +
		"END  " +
		"ELSE  " +
		"BEGIN  " +
		"INSERT INTO [mdyystaterestoretime]([mdid],[mdyydate],[chanid],[openuprcvr],[closeuprcvr],[lastupdate])  " +
		"VALUES (?,?,?,?,?,getDate())  " +
		"END "
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
	//log.Debug(connStr[0])
	return comm.getDBConfigByStr(connStr[0])
}

//根据Key获取任务执行时间配置
func (c *config) GetTaskCronByKey(key object.TaskKey) (*object.TaskCron, error) {
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

//根据操作更新门店设置
func (c *config) UpdateMdSet(opr *object.OprMdSet) error {
	switch opr.OprType {
	case 1:
		return c.mdSetInsert(opr)
	case 2:
		return c.mdSetUpdate(opr)
	case 3:
		return c.mdSetDelete(opr)
	default:
		errMsg := fmt.Sprintf("md set opr type error,exp 1 or 2 or 3 got %d", opr.OprType)
		log.Error(errMsg)
		return errors.New(errMsg)
	}
}

//新增门店设置
func (c *config) mdSetInsert(opr *object.OprMdSet) error {
	comm := NewCommon()
	return comm.SetRowsBySQL(c.dbConfig, sqlMdSetInsert, opr.MdId, opr.MdId, opr.MdId, opr.MdName, opr.MdCode)
}

//更新门店设置
func (c *config) mdSetDelete(opr *object.OprMdSet) error {
	comm := NewCommon()
	return comm.SetRowsBySQL(c.dbConfig, sqlMdSetDelete, opr.MdId)
}

//删除门店设置
func (c *config) mdSetUpdate(opr *object.OprMdSet) error {
	comm := NewCommon()
	return comm.SetRowsBySQL(c.dbConfig, sqlMdSetUpdate, opr.MdName, opr.MdCode, opr.MdId)
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
	comm := NewCommon()
	return comm.SetRowsBySQL(c.dbConfig, sqlDeleteCwGsSet, opr.GsId)
}

//删除财务公司设置
func (c *config) cwGsSetUpdate(opr *object.OprCwGsSet) error {
	comm := NewCommon()
	return comm.SetRowsBySQL(c.dbConfig, sqlUpdateCwGsSet, opr.GsName, opr.GsId)
}

//根据操作更新门店财务公司关系
func (c *config) UpdateMdCwGsRef(opr *object.OprMdCwGsRef) error {
	switch opr.OprType {
	case 1:
		return c.mdCwGsRefInsert(opr)
	case 2:
		return c.mdCwGsRefUpdate(opr)
	case 3:
		return c.mdCwGsRefDelete(opr)
	default:
		errMsg := fmt.Sprintf("mdcwgsref opr type error,exp 1 or 2 or 3 got %d", opr.OprType)
		log.Error(errMsg)
		return errors.New(errMsg)
	}
}

//新增门店财务公司关系
func (c *config) mdCwGsRefInsert(opr *object.OprMdCwGsRef) error {
	comm := NewCommon()
	return comm.SetRowsBySQL(c.dbConfig, sqlMdCwGsRefInsert, opr.MdId, opr.GsId)
}

//更新门店财务公司关系
func (c *config) mdCwGsRefDelete(opr *object.OprMdCwGsRef) error {
	comm := NewCommon()
	return comm.SetRowsBySQL(c.dbConfig, sqlMdCwGsRefDelete, opr.MdId)
}

//删除门店财务公司关系
func (c *config) mdCwGsRefUpdate(opr *object.OprMdCwGsRef) error {
	comm := NewCommon()
	return comm.SetRowsBySQL(c.dbConfig, sqlMdCwGsRefUpdate, opr.GsId, opr.MdId)
}

//根据操作更新门店营业日开闭店记录传递时间
func (c *config) UpdateMdYyStateTransTimeTd(opr *object.OprMdYyStateTransTimeTd) error {
	comm := NewCommon()
	switch opr.OprType {
	case 1:
		return comm.SetRowsBySQL(
			c.dbConfig,
			sqlUpdateMdYyStateTransTimeTd1,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			opr.OprTime,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			opr.OprTime, comm.GetDefaultOprTime(), comm.GetDefaultOprTime(), comm.GetDefaultOprTime())
	case 2:
		return comm.SetRowsBySQL(
			c.dbConfig,
			sqlUpdateMdYyStateTransTimeTd2,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			opr.OprTime,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			comm.GetDefaultOprTime(), opr.OprTime, comm.GetDefaultOprTime(), comm.GetDefaultOprTime())
	case 3:
		return comm.SetRowsBySQL(
			c.dbConfig,
			sqlUpdateMdYyStateTransTimeTd3,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			opr.OprTime,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			comm.GetDefaultOprTime(), comm.GetDefaultOprTime(), opr.OprTime, comm.GetDefaultOprTime())
	case 4:
		return comm.SetRowsBySQL(
			c.dbConfig,
			sqlUpdateMdYyStateTransTimeTd4,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			opr.OprTime,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			comm.GetDefaultOprTime(), comm.GetDefaultOprTime(), comm.GetDefaultOprTime(), opr.OprTime)
	default:
		errMsg := fmt.Sprintf("MdYyStateTransTimeTd opr type error,exp 1 or 2 or 3 or 4 got %d", opr.OprType)
		log.Error(errMsg)
		return errors.New(errMsg)
	}
}

//根据操作更新门店营业日开闭店操作
func (c *config) UpdateMdYyStateRestoreTime(opr *object.OprMdYyStateRestoreTimeZb) error {
	comm := NewCommon()
	switch opr.OprType {
	case 1:
		return comm.SetRowsBySQL(
			c.dbConfig,
			sqlUpdateMdYyStateRestoreTime1,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			opr.OprTime,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			opr.MdId, opr.MdYyDate, opr.ChanId, opr.OprTime, comm.GetDefaultOprTime())
	case 2:
		return comm.SetRowsBySQL(
			c.dbConfig,
			sqlUpdateMdYyStateRestoreTime2,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			opr.OprTime,
			opr.MdId, opr.MdYyDate, opr.ChanId,
			opr.MdId, opr.MdYyDate, opr.ChanId, comm.GetDefaultOprTime(), opr.OprTime)
	default:
		errMsg := fmt.Sprintf("OprMdYyStateRestoreTimeZb opr type error,exp 1 or 2 got %d", opr.OprType)
		log.Error(errMsg)
		return errors.New(errMsg)
	}
}

//根据操作更新门店营业日开闭店操作
func (c *config) UpdateMdYyState(opr *object.OprMdYyState) error {
	comm := NewCommon()
	return comm.SetRowsBySQL(
		c.dbConfig,
		sqlUpdateMdYyState,
		opr.MdId, opr.MdYyDate,
		opr.MdYyOpenTime, opr.MdYyCloseTime, opr.MdYySjType,
		opr.MdId, opr.MdYyDate,
		opr.MdId, opr.MdYyDate, opr.MdYyOpenTime, opr.MdYyCloseTime, opr.MdYySjType)
}

//获取财务公司设置
func (c *config) GetCwGsSet() ([]*object.CwGsSet, error) {
	comm := common{}
	rows, err := comm.GetRowsBySQL(c.dbConfig, sqlGetCwGsSet)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	rList := make([]*object.CwGsSet, 0)
	var gsId int
	var gsName string
	var lastUpdate time.Time
	for rows.Next() {
		err = rows.Scan(&gsId, &gsName, &lastUpdate)
		if err != nil {
			errMsg := fmt.Sprintf("get cwgsset read data err: %s", err.Error())
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		rList = append(rList, &object.CwGsSet{
			GsId:       gsId,
			GsName:     gsName,
			LastUpdate: lastUpdate,
		})
	}
	if rows.Err() != nil {
		errMsg := fmt.Sprintf("get cwgsset read data err: %s", rows.Err().Error())
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return rList, nil
}
