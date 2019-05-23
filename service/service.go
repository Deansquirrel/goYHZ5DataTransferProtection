package service

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/global"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
	rep "github.com/Deansquirrel/goYHZ5DataTransferProtection/repository"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/worker"
	"github.com/kataras/iris/core/errors"
	"github.com/robfig/cron"
	"sort"
	"strings"
)

//启动服务内容
func StartService() error {
	log.Debug("StartService")
	startHeartBeat()
	startRefreshMdDataTransState()
	startRestoreMdYyStateTransTime()
	startRefreshWaitRestoreDataCount()
	startRestoreMdYyStateRestoreTime()
	startRestoreMdYyState()
	startRestoreMdSet()
	startRestoreCwGsSet()
	startRestoreMdCwGsRef()

	//Test
	{
		go func() {
			c := cron.New()
			err := c.AddFunc("5/30 * * * * ?", showTaskState)
			if err != nil {
				log.Debug(err.Error())
			} else {
				c.Start()
			}
		}()
	}

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

//Common（心跳）
func startHeartBeat() {
	ch := make(chan error)
	comm := worker.NewCommon(ch)
	startWorker(object.TaskKeyHeartBeat, comm.RefreshHeartBeat, ch)
}

//通道
func startRefreshMdDataTransState() {
	ch := make(chan error)
	workerTd := worker.NewWorkerTd(ch)
	startWorker(object.TaskKeyRefreshMdDataTransState, workerTd.RefreshMdDataTransState, ch)
}

//通道
func startRestoreMdYyStateTransTime() {
	ch := make(chan error)
	workerTd := worker.NewWorkerTd(ch)
	startWorker(object.TaskKeyRestoreMdYyStateTransTime, workerTd.RestoreMdYyStateTransTime, ch)
}

//总部库
func startRefreshWaitRestoreDataCount() {
	ch := make(chan error)
	w := worker.NewWorkerZb(ch)
	startWorker(object.TaskKeyRefreshWaitRestoreDataCount, w.RefreshWaitRestoreDataCount, ch)
}

//总部库
func startRestoreMdYyStateRestoreTime() {
	ch := make(chan error)
	w := worker.NewWorkerZb(ch)
	startWorker(object.TaskKeyRestoreMdYyStateRestoreTime, w.RestoreMdYyStateRestoreTime, ch)
}

//总部库
func startRestoreMdYyState() {
	ch := make(chan error)
	w := worker.NewWorkerZb(ch)
	startWorker(object.TaskKeyRestoreMdYyState, w.RestoreMdYyState, ch)
}

//总部库
func startRestoreMdSet() {
	ch := make(chan error)
	w := worker.NewWorkerZb(ch)
	startWorker(object.TaskKeyRestoreMdSet, w.RestoreMdSet, ch)
}

//总部库
func startRestoreCwGsSet() {
	ch := make(chan error)
	w := worker.NewWorkerZb(ch)
	startWorker(object.TaskKeyRestoreCwGsSet, w.RestoreCwGsSet, ch)
}

//总部库
func startRestoreMdCwGsRef() {
	ch := make(chan error)
	w := worker.NewWorkerZb(ch)
	startWorker(object.TaskKeyRestoreMdCwGsRef, w.RestoreMdCwGsRef, ch)
}

//启动工作线程
func startWorker(key object.TaskKey, cmd func(), ch chan error) {
	s := &object.TaskState{
		Key:     key,
		Cron:    nil,
		CronStr: "",
		Running: false,
		Err:     nil,
	}
	global.TaskList.Register() <- goToolCommon.NewObject(string(s.Key), s)
	cronStr, err := getTaskCron(s.Key)
	if err != nil {
		s.CronStr = ""
		s.Err = err
		return
	}
	cronStr = strings.Trim(cronStr, " ")
	if cronStr == "" {
		s.Err = errors.New("cron is empty")
		return
	}
	s.CronStr = cronStr
	c := cron.New()
	err = c.AddFunc(s.CronStr, cmd)
	if err != nil {
		s.Err = err
		return
	}
	c.Start()
	s.Running = true
	s.Cron = c
	go func() {
		defer func() {
			close(ch)
		}()
		for {
			select {
			case err := <-ch:
				s.Err = err
			case <-global.Ctx.Done():
				return
			}
		}
	}()
}

//根据key获取任务执行Cron时间公式
func getTaskCron(key object.TaskKey) (string, error) {
	repConfig, err := rep.NewConfig()
	if err != nil {
		return "", err
	}
	taskCron, err := repConfig.GetTaskCronByKey(key)
	if err != nil {
		return "", err
	}
	if taskCron == nil {
		return "", errors.New("get taskCron err: return is nil")
	}
	return taskCron.Cron, nil
}
