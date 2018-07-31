package task

import (
	"errors"

	"github.com/robfig/cron"
)

type TaskManager struct {
	etlTasks  map[string]*EtlTask
	scheduler *cron.Cron
	log       LogFunc
}

type LogFunc func(ll LogLevel, errMsg string)

func NewTaskManager() *TaskManager {
	tm := new(TaskManager)
	tm.etlTasks = make(map[string]*EtlTask)
	tm.scheduler = cron.New()
	return tm
}

func (t *TaskManager) SetLog(log LogFunc) {
	t.log = log
}

func (t *TaskManager) GetAllTasks() map[string]*EtlTask {
	return t.etlTasks
}

func (t *TaskManager) Add(aTask *EtlTask) *TaskManager {
	if _, ok := t.etlTasks[aTask.Batch.GetName()]; !ok {
		t.etlTasks[aTask.Batch.GetName()] = aTask
	} else {
		panic("there have existed a same task, name:" + aTask.Batch.GetName())
	}

	return t
}

func (t *TaskManager) Reset(taskName string) error {
	if tarTsk, ok := t.etlTasks[taskName]; ok {
		return tarTsk.Batch.Reset()
	}
	return nil
}

func (t *TaskManager) Start(taskName string) error {
	if tarTsk, ok := t.etlTasks[taskName]; ok {
		return tarTsk.Start()
	}

	return nil
}

func (t *TaskManager) Stop(taskName string) error {
	if tarTsk, ok := t.etlTasks[taskName]; ok {
		return tarTsk.Stop()
	}
	return nil
}

func (t *TaskManager) Execute(taskName string) error {
	if tarTsk, ok := t.etlTasks[taskName]; ok {
		if tarTsk.State == Completed || tarTsk.State == Stopped || tarTsk.State == Executing {
			return errors.New("Can not execute a unready task(running or init)")
		}
		go tarTsk.Execute()
	}
	return nil
}

func (t *TaskManager) Build() error {
	for _, v := range t.etlTasks {
		if v.Type == OneShot {
			if v.PlanTime != "" {
				err := t.scheduler.AddJob(v.PlanTime, NewTaskWrapper(t.log, v))
				if err != nil {
					return err
				}
			} else {
				go v.Execute()
			}

		}

		if v.Type == Plan {
			err := t.scheduler.AddJob(v.PlanTime, NewTaskWrapper(t.log, v))
			if err != nil {
				return err
			}
		}
	}

	t.scheduler.Start()
	return nil
}
