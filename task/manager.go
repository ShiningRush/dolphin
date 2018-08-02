package task

import (
	"errors"

	"github.com/robfig/cron"
)

type TaskManager struct {
	etlTasks          map[string]*EtlTask
	scheduler         *cron.Cron
	log               LogFunc
	taskStatusChanged chan *EtlTask
	taskStatusRepo    TaskStatusRepository
}

type LogFunc func(ll LogLevel, errMsg string)

type TaskStatusRepository interface {
	GetAll() ([]TaskStatus, error)
	InsertOrUpdate(ts TaskStatus) error
	RemoveLegacy(newTs map[string]*EtlTask) error
}

func NewTaskManager() *TaskManager {
	tm := new(TaskManager)
	tm.etlTasks = make(map[string]*EtlTask)
	tm.scheduler = cron.New()
	tm.taskStatusChanged = make(chan *EtlTask, 100)
	return tm
}

// SetRepo if you need persist task status
func (t *TaskManager) SetRepo(repo TaskStatusRepository) {
	t.taskStatusRepo = repo
}

// SetLog if you want log err messages
func (t *TaskManager) SetLog(log LogFunc) {
	t.log = log
}

// GetAllTasks this method is commonly used by dashserver
func (t *TaskManager) GetAllTasks() map[string]*EtlTask {
	return t.etlTasks
}

// Add add a name unique task
func (t *TaskManager) Add(aTask *EtlTask) *TaskManager {
	if _, ok := t.etlTasks[aTask.Batch.GetName()]; !ok {
		aTask.StatusChanged = t.taskStatusChanged
		t.etlTasks[aTask.Batch.GetName()] = aTask
	} else {
		panic("there have existed a same task, name:" + aTask.Batch.GetName())
	}

	return t
}

// Reset a task
func (t *TaskManager) Reset(taskName string) error {
	if tarTsk, ok := t.etlTasks[taskName]; ok {
		return tarTsk.Batch.Reset()
	}
	return nil
}

// Start a task
func (t *TaskManager) Start(taskName string) error {
	if tarTsk, ok := t.etlTasks[taskName]; ok {
		return tarTsk.Start()
	}

	return nil
}

// Stop a task
func (t *TaskManager) Stop(taskName string) error {
	if tarTsk, ok := t.etlTasks[taskName]; ok {
		return tarTsk.Stop()
	}
	return nil
}

// Execute a task
func (t *TaskManager) Execute(taskName string) error {
	if tarTsk, ok := t.etlTasks[taskName]; ok {
		if tarTsk.State == Completed || tarTsk.State == Stopped || tarTsk.State == Executing {
			return errors.New("Can not execute a unready task(running or init)")
		}
		if err := tarTsk.Execute(); err != nil {
			t.log(Error, err.Error())
		}
	}
	return nil
}

// Build taskmanager, it will do oneshot  task first and add plan task to cron
// then if you specified repo, will recover task status and remove unnecessary task status
func (t *TaskManager) Build() error {
	if err := t.dispatchTask(); err != nil {
		return err
	}

	if t.taskStatusRepo != nil {
		if err := t.initTaskStatus(); err != nil {
			return errors.New("init task status error:" + err.Error())
		}

		if err := t.taskStatusRepo.RemoveLegacy(t.etlTasks); err != nil {
			return errors.New("Remove legacy tasks error:" + err.Error())
		}
		go t.handleTaskStatusChanged()
	}

	t.scheduler.Start()
	return nil
}

func (t *TaskManager) dispatchTask() error {
	for _, v := range t.etlTasks {
		if v.Type == OneShot {
			if v.PlanTime != "" {
				err := t.scheduler.AddJob(v.PlanTime, NewTaskWrapper(t.log, v))
				if err != nil {
					return err
				}
			} else {
				if err := v.Execute(); err != nil {
					t.log(Error, err.Error())
				}
			}

		}

		if v.Type == Plan {
			err := t.scheduler.AddJob(v.PlanTime, NewTaskWrapper(t.log, v))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *TaskManager) initTaskStatus() error {
	allTaskStatus, err := t.taskStatusRepo.GetAll()
	if err != nil {
		return err
	}

	for k, task := range t.etlTasks {
		for _, ts := range allTaskStatus {
			if k == ts.TaskName {
				task.Status = ts.Status
			}
		}
	}

	return nil
}

func (t *TaskManager) handleTaskStatusChanged() {
	for v := range t.taskStatusChanged {
		t.taskStatusRepo.InsertOrUpdate(v.GetTaskStatus())
	}
}
