package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/global"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/worker"
	"github.com/robfig/cron"
	"sort"
)

const (
	//TODO 待参数化
	webHook = "3dd601535aa95027b92e78b8a820ba62be5069293092a302d8c17ef63e095cac"
)

//启动服务内容
func StartService() error {
	log.Debug("StartService")

	ch := make(chan error)
	commWorker := worker.NewCommon(ch)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case err := <-ch:
				if err != nil {
					log.Debug(err.Error())
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	//主任务
	{
		commWorker.StartService(object.TaskKeyHeartBeat, errHandle)
		commWorker.StartService(object.TaskKeyRefreshMdDataTransState, errHandle)
		commWorker.StartService(object.TaskKeyRestoreMdYyStateTransTime, errHandle)
		commWorker.StartService(object.TaskKeyRefreshWaitRestoreDataCount, errHandle)
		commWorker.StartService(object.TaskKeyRestoreMdYyStateRestoreTime, errHandle)
		commWorker.StartService(object.TaskKeyRestoreMdYyState, errHandle)
		commWorker.StartService(object.TaskKeyRestoreMdSet, errHandle)
		commWorker.StartService(object.TaskKeyRestoreCwGsSet, errHandle)
		commWorker.StartService(object.TaskKeyRestoreMdCwGsRef, errHandle)
	}

	//辅助任务
	{
		commWorker.StartService(object.TaskKeyRefreshConfig, errHandle)
	}

	//Test 显示任务运行情况
	{
		go func() {
			c := cron.New()
			err := c.AddFunc("0 * * * * ?", showTaskState)
			if err != nil {
				log.Debug(err.Error())
			} else {
				c.Start()
			}
		}()
	}
	cancel()
	return nil
}

//显示任务状态
func showTaskState() {
	idList := global.TaskList.GetIdList()
	sort.Strings(idList)

	for _, k := range idList {
		testPrintS(global.TaskList.GetObject(k).(*object.TaskState))
	}
}

//输出TaskState
func testPrintS(ts *object.TaskState) {
	if ts == nil {
		log.Debug("ts is nil")
		return
	}
	str, err := json.Marshal(ts)
	if err != nil {
		log.Debug(fmt.Sprintf("%s - %s", ts.Key, err.Error()))
	} else {
		log.Debug(fmt.Sprintf("%s - %s", ts.Key, str))
	}
	c := ts.Cron
	if c != nil {
		if len(c.Entries()) > 0 {
			e := c.Entries()[0]
			log.Debug(fmt.Sprintf("Prev time: %s,Next time: %s",
				goToolCommon.GetDateTimeStr(e.Prev),
				goToolCommon.GetDateTimeStr(e.Next)))
		}
	}
}

func errHandle(err error) {
	if err == nil {
		return
	}
	log.Error(err.Error())
	msg := err.Error()
	sendErr := sendDingTalkTextMsg(msg)
	if sendErr != nil {
		log.Error(fmt.Sprintf("send msg,msg: %s, error: %s", msg, sendErr.Error()))
	}
}

func sendDingTalkTextMsg(msg string) error {
	dt := dingTalkRobot{
		config: &object.DingTalkRobotConfigData{
			FWebHookKey: webHook,
			FAtMobiles:  "",
			FIsAtAll:    0,
		},
	}
	return dt.SendMsg(msg)
}
