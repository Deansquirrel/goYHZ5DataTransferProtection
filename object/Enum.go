package object

type TaskKey string

//任务Key
const (
	TaskKeyHeartBeat                   TaskKey = "HeartBeat"
	TaskKeyRefreshMdDataTransState     TaskKey = "RefreshMdDataTransState"
	TaskKeyRestoreMdYyStateTransTime   TaskKey = "RestoreMdYyStateTransTime"
	TaskKeyRefreshWaitRestoreDataCount TaskKey = "RefreshWaitRestoreDataCount"
	TaskKeyRestoreMdYyStateRestoreTime TaskKey = "RestoreMdYyStateRestoreTime"
	TaskKeyRestoreMdYyState            TaskKey = "RestoreMdYyState"
	TaskKeyRestoreMdSet                TaskKey = "RestoreMdSet"
	TaskKeyRestoreCwGsSet              TaskKey = "RestoreCwGsSet"
	TaskKeyRestoreMdCwGsRef            TaskKey = "RestoreMdCwGsRef"
)

//数据库连接类型key
type ConnTypeKey string

const (
	ConnTypeKeyZb ConnTypeKey = "zb"
	ConnTypeKeyTd ConnTypeKey = "td"
)

//通道库中的客户识别码
const TdUserCode = "mt_dlyh_201008"
