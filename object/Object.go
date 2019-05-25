package object

import (
	"github.com/robfig/cron"
	"time"
)

type TaskState struct {
	Key     TaskKey
	Cron    *cron.Cron
	CronStr string
	Running bool
	Err     error
}

//门店数据传递信道状态
type MdDataTransState struct {
	BatchNo      string    //操作批次号
	MdCode       string    //门店通道码
	ChanId       int       //信道ID
	DataRowCount int       //暂存数据总量
	LastUp       time.Time //门店最后上传时间
	LastUpRcv    time.Time //总部最后接收时间
	RecordTime   time.Time //记录时间
}

//门店营业日开闭店记录传递时间（通道库）
type MdYyStateTransTimeTd struct {
	MdId       int       //门店ID
	MdYyDate   string    //门店营业日
	OpenUp     time.Time //开店上传时间
	OpenUpRcv  time.Time //开店接收时间
	CloseUp    time.Time //闭店上传时间
	CloseUpRcv time.Time //闭店接收时间
	LastUpdate time.Time //最后更新时间
}

type MdYyStateTransTimeTdTrans struct {
	OprSn    int       //操作顺序号
	MdId     int       //门店ID
	MdYyDate string    //门店营业日
	OprType  int       //操作类型，1-开店上传、2-开店接收、3-闭店上传、4-闭店接收
	OprTime  time.Time //操作时间
}

//门店营业日开闭店记录恢复时间（总部库）
type MdYyStateTransTimeZb struct {
	MdId         int       //门店ID
	MdYyDate     string    //门店营业日
	OpenRestore  time.Time //开店恢复时间
	CloseRestore time.Time //闭店恢复时间
	LastUpdate   time.Time //最后更新时间
}

//门店营业日开闭店记录（总部库）
type MdYyState struct {
	MdId          int       //门店ID
	MdYyDate      string    //门店营业日
	MdYyOpenTime  time.Time //开店时间
	MdYyCloseTime time.Time //闭店时间
	MdYySjType    int       //数据完成标识
	LastUpdate    time.Time //最后更新时间
}

//待恢复信道数据总量（总部库）
type WaitRestoreDataCount struct {
	BatchNo      string    //操作批次号
	TfRcv1       int       //待恢复数据总量1
	GetDataTime1 time.Time //最早取走时间1
	TfRcv2       int       //待恢复数据总量2
	GetDataTime2 time.Time //最早取走时间2
	TfRcv3       int       //待恢复数据总量3
	GetDataTime3 time.Time //最早取走时间3
	TfRcv4       int       //待恢复数据总量4
	GetDataTime4 time.Time //最早取走时间4
	RecordDate   time.Time //记录时间
}

//门店设置
type MdSet struct {
	MdId        int       //门店ID
	MdName      string    //门店名称
	MdCode      string    //门店通道码
	IsForbidden int       //是否禁用（0-在用，1-禁用）
	LastUpdate  time.Time //最后更新时间
}

//门店设置操作记录（触发器触发）
type OprMdSet struct {
	OprSn   int       //操作SN
	MdId    int       //门店ID
	MdName  string    //门店名称
	MdCode  string    //门店通道码
	OprType int       //操作类型，1-新增（启用）、2-修改、3-删除
	OprTime time.Time //操作时间
}

//财务公司设置
type CwGsSet struct {
	GsId       int       //公司ID
	GsName     string    //公司名称
	LastUpdate time.Time //最后更新时间
}

//财务公司设置操作记录（触发器触发）
type OprCwGsSet struct {
	OprSn   int       //操作SN
	GsId    int       //公司ID
	GsName  string    //公司名称
	OprType int       //操作类型，1-新增（启用）、2-修改、3-删除
	OprTime time.Time //操作时间
}

//门店财务公司关系设置
type MdCwGsRef struct {
	MdId       int       //门店ID
	GsId       int       //公司ID
	LastUpdate time.Time //最后更新时间
}

//门店财务公司关系设置操作记录（触发器触发）
type OprMdCwGsRef struct {
	OprSn   int       //操作SN
	MdId    int       //门店ID
	GsId    int       //公司ID
	OprType int       //操作类型，1-新增、2-修改、3-删除
	OprTime time.Time //操作时间
}

//连接配置
type ConnInfo struct {
	TdConn string //通道库连接配置
	Z5Conn string //Z5总部库连接配置
}

//任务执行时间配置表
type TaskCron struct {
	TaskKey         TaskKey //任务标识
	TaskDescription string  //任务描述
	Cron            string  //任务执行cron公式
}

//固化项目表
type ZlFixedList struct {
	FlId    string //固化项ID
	FlSign  string //固化项名称
	FlName  string //固化项中文
	FlIndex string //固化项索引码
	FlSn    string //固化项顺序
}
