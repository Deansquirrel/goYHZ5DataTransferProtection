package global

import (
	"context"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
)

const (
	//PreVersion = "0.0.4 Build20190529"
	//TestVersion = "0.0.0 Build20190101"
	Version = "0.0.0 Build20190101"
)

var Ctx context.Context
var Cancel func()

//程序启动参数
var Args *object.ProgramArgs

//系统参数
var SysConfig *object.SystemConfig

//TaskList
var TaskList goToolCommon.IObjectManager

var TaskKeyList []object.TaskKey

func init() {
	TaskKeyList = make([]object.TaskKey, 0)

	//主任务
	TaskKeyList = append(TaskKeyList, object.TaskKeyHeartBeat)
	TaskKeyList = append(TaskKeyList, object.TaskKeyRefreshMdDataTransState)
	TaskKeyList = append(TaskKeyList, object.TaskKeyRestoreMdYyStateTransTime)
	TaskKeyList = append(TaskKeyList, object.TaskKeyRefreshWaitRestoreDataCount)
	TaskKeyList = append(TaskKeyList, object.TaskKeyRestoreMdYyStateRestoreTime)
	TaskKeyList = append(TaskKeyList, object.TaskKeyRestoreMdYyState)
	TaskKeyList = append(TaskKeyList, object.TaskKeyRestoreMdSet)
	TaskKeyList = append(TaskKeyList, object.TaskKeyRestoreCwGsSet)
	TaskKeyList = append(TaskKeyList, object.TaskKeyRestoreMdCwGsRef)

	//辅助任务
	TaskKeyList = append(TaskKeyList, object.TaskKeyRefreshConfig)
}
