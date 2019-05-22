package service

import (
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/global"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
	rep "github.com/Deansquirrel/goYHZ5DataTransferProtection/repository"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/worker"
	"github.com/kataras/iris/core/errors"
	"github.com/robfig/cron"
	"strings"
)

//启动服务内容
func StartService() error {
	log.Debug("StartService")
	startHeartBeat()
	startRefreshMdDataTransState()
	startRestoreMdYyStateTransTime()

	return nil
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
	log.Debug(fmt.Sprintf("根据key获取任务执行Cron时间公式 %s", key))
	//TODO 根据key获取任务执行Cron时间公式
	repCommon := rep.NewCommon()
	taskCron, err := repCommon.GetTaskCronByKey(string(key))
	if err != nil {
		return "", err
	}
	if taskCron == nil {
		return "", errors.New("get taskCron err: return is nil")
	}
	return taskCron.Cron, nil
}
