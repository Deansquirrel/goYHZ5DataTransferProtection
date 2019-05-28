package worker

import (
	"context"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/global"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/repository"
	"github.com/robfig/cron"
	"strings"
)

type common struct {
	errChan chan<- error
}

func NewCommon(errChan chan<- error) *common {
	return &common{
		errChan: errChan,
	}
}

//刷新心跳时间
func (c *common) RefreshHeartBeat() {
	log.Debug("刷新心跳时间")
	rep, err := repository.NewConfig()
	if err != nil {
		c.errChan <- err
		return
	}
	err = rep.UpdateHeartBeat()
	if err != nil {
		c.errChan <- err
		return
	}
	c.errChan <- nil
}

//刷新配置
func (c *common) RefreshConfig() {
	log.Debug("配置刷新")
	rep, err := repository.NewConfig()
	if err != nil {
		c.errChan <- err
		return
	}
	idList := global.TaskList.GetIdList()
	for _, id := range idList {
		t := global.TaskList.GetObject(id)
		if t == nil {
			c.errChan <- errors.New(fmt.Sprintf("task cron is nil,key %s", id))
			return
		}
		ts := t.(*object.TaskState)
		newTs, err := rep.GetTaskCronByKey(ts.Key)
		if err != nil {
			c.errChan <- err
			return
		}
		if newTs.Cron != ts.CronStr {
			//TODO 重启服务
		}
	}
	c.errChan <- nil
}

//根据Key，启动制定服务
func (c *common) StartService(key object.TaskKey, errHandle func(err error)) {
	switch key {
	case object.TaskKeyHeartBeat:
		c.startHeartBeat(errHandle)
	case object.TaskKeyRefreshMdDataTransState:
		c.startRefreshMdDataTransState(errHandle)
	case object.TaskKeyRestoreMdYyStateTransTime:
		c.startRestoreMdYyStateTransTime(errHandle)
	case object.TaskKeyRefreshWaitRestoreDataCount:
		c.startRefreshWaitRestoreDataCount(errHandle)
	case object.TaskKeyRestoreMdYyStateRestoreTime:
		c.startRestoreMdYyStateRestoreTime(errHandle)
	case object.TaskKeyRestoreMdYyState:
		c.startRestoreMdYyState(errHandle)
	case object.TaskKeyRestoreMdSet:
		c.startRestoreMdSet(errHandle)
	case object.TaskKeyRestoreCwGsSet:
		c.startRestoreCwGsSet(errHandle)
	case object.TaskKeyRestoreMdCwGsRef:
		c.startRestoreMdCwGsRef(errHandle)
	case object.TaskKeyRefreshConfig:
		c.startRefreshConfig(errHandle)
	default:
		errMsg := fmt.Sprintf("unknow task key: %s", key)
		log.Error(errMsg)
		c.errChan <- errors.New(errMsg)
	}
}

//刷新配置
func (c *common) startRefreshConfig(errHandle func(err error)) {
	ch := make(chan error)
	comm := NewCommon(ch)
	c.startWorker(object.TaskKeyRefreshConfig, comm.RefreshConfig, ch, errHandle)
}

//Common（心跳）
func (c *common) startHeartBeat(errHandle func(err error)) {
	ch := make(chan error)
	comm := NewCommon(ch)
	c.startWorker(object.TaskKeyHeartBeat, comm.RefreshHeartBeat, ch, errHandle)
}

//通道
func (c *common) startRefreshMdDataTransState(errHandle func(err error)) {
	ch := make(chan error)
	workerTd := NewWorkerTd(ch)
	c.startWorker(object.TaskKeyRefreshMdDataTransState, workerTd.RefreshMdDataTransState, ch, errHandle)
}

//通道
func (c *common) startRestoreMdYyStateTransTime(errHandle func(err error)) {
	ch := make(chan error)
	workerTd := NewWorkerTd(ch)
	c.startWorker(object.TaskKeyRestoreMdYyStateTransTime, workerTd.RestoreMdYyStateTransTime, ch, errHandle)
}

//总部库
func (c *common) startRefreshWaitRestoreDataCount(errHandle func(err error)) {
	ch := make(chan error)
	w := NewWorkerZb(ch)
	c.startWorker(object.TaskKeyRefreshWaitRestoreDataCount, w.RefreshWaitRestoreDataCount, ch, errHandle)
}

//总部库
func (c *common) startRestoreMdYyStateRestoreTime(errHandle func(error)) {
	ch := make(chan error)
	w := NewWorkerZb(ch)
	c.startWorker(object.TaskKeyRestoreMdYyStateRestoreTime, w.RestoreMdYyStateRestoreTime, ch, errHandle)
}

//总部库
func (c *common) startRestoreMdYyState(errHandle func(err error)) {
	ch := make(chan error)
	w := NewWorkerZb(ch)
	c.startWorker(object.TaskKeyRestoreMdYyState, w.RestoreMdYyState, ch, errHandle)
}

//总部库
func (c *common) startRestoreMdSet(errHandle func(err error)) {
	ch := make(chan error)
	w := NewWorkerZb(ch)
	c.startWorker(object.TaskKeyRestoreMdSet, w.RestoreMdSet, ch, errHandle)
}

//总部库
func (c *common) startRestoreCwGsSet(errHandle func(err error)) {
	ch := make(chan error)
	w := NewWorkerZb(ch)
	c.startWorker(object.TaskKeyRestoreCwGsSet, w.RestoreCwGsSet, ch, errHandle)
}

//总部库
func (c *common) startRestoreMdCwGsRef(errHandle func(err error)) {
	ch := make(chan error)
	w := NewWorkerZb(ch)
	c.startWorker(object.TaskKeyRestoreMdCwGsRef, w.RestoreMdCwGsRef, ch, errHandle)
}

//启动工作线程
func (c *common) startWorker(key object.TaskKey, cmd func(), ch chan error, errHandle func(err error)) {
	s := &object.TaskState{
		Key:     key,
		Cron:    nil,
		CronStr: "",
		Running: false,
		Err:     nil,
	}

	global.TaskList.Register() <- goToolCommon.NewObject(string(s.Key), s)
	cronStr, err := c.getTaskCron(s.Key)
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
	cr := cron.New()
	err = cr.AddFunc(s.CronStr, cmd)
	if err != nil {
		s.Err = err
		return
	}

	s.Ctx, s.Cancel = context.WithCancel(context.Background())

	cr.Start()
	s.Running = true
	s.Cron = cr
	go func() {
		defer func() {
			close(ch)
		}()
		for {
			select {
			case err := <-ch:
				errHandle(err)
				s.Err = err
			case <-s.Ctx.Done():
				s.Cron.Stop()
				return
			case <-global.Ctx.Done():
				s.Cron.Stop()
				return
			}
		}
	}()
}

//根据key获取任务执行Cron时间公式
func (c *common) getTaskCron(key object.TaskKey) (string, error) {
	repConfig, err := repository.NewConfig()
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
