package task

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	cfgPath string
	cfgName string
)

// SetTaskConfig set config dir for config file, default search for config, cfgName is no need extension
func SetTaskConfig(dir string, name string) {
	cfgPath = dir
	cfgName = strings.ToLower(name)
}

// EtlTask indicate a task operate etl
type EtlTask struct {
	Batch            Batch
	Syncers          []Syncer
	Type             TaskType
	State            TaskState
	PlanTime         string // this time use cron format time
	LastExecuteTime  time.Time
	LastExecuteState string
	LastExecuteCost  int
}

// Batch interface
type Batch interface {
	GetName() string
	Begin() error
	Reset() error
}

// Syncer interface uset to sync data
type Syncer interface {
	Start()
	Stop()
}

// NewTask create a new etltask with batch
func NewTask(aBatch Batch, syncers ...Syncer) *EtlTask {
	t := new(EtlTask)
	t.Syncers = syncers
	t.Type = OneShot
	t.State = Init
	t.Batch = aBatch
	if cfgPath != "" && cfgName != "" {
		setTaskWithConfig(t)
	}

	return t
}

func setTaskWithConfig(task *EtlTask) {
	v := viper.New()
	v.SetConfigName(cfgName)
	v.AddConfigPath(cfgPath)
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error when read config file: %s \n", err))
	}

	cfgType := v.GetString("taskConfig." + task.Batch.GetName() + ".type")
	planTime := v.GetString("taskConfig." + task.Batch.GetName() + ".planTime")
	if cfgType != "" {
		switch cfgType {
		case "Plan":
			task.Type = Plan
		case "OneShot":
			task.Type = OneShot
		}
	}

	if planTime != "" {
		task.PlanTime = planTime
	}
}

// Execute a task
func (e *EtlTask) Execute() error {
	if e.State == Completed || e.State == Stopped || e.State == Executing {
		return nil
	}

	e.StopSyncers()
	executeTime := time.Now()
	e.LastExecuteState = ""
	e.State = Executing
	var allError error
	if err := e.Batch.Reset(); err != nil {
		e.LastExecuteState = "This task has something error when resetting, please check log"
		errMsg := "We get a error when reset batch, batch name:" + e.Batch.GetName() + ", errors:" + err.Error()
		allError = errors.New(errMsg)
	}

	if err := e.Batch.Begin(); err != nil {
		e.LastExecuteState = "This task has something error when beginning, please check log"
		errMsg := "We get a error when begin batch, batch name:" + e.Batch.GetName() + ", errors:" + err.Error()
		allError = errors.New(errMsg)
	}

	e.LastExecuteCost = int(time.Since(executeTime) / time.Second)
	e.LastExecuteTime = executeTime
	e.StartSyncers()

	if e.Type == OneShot {
		e.State = Completed
	} else {
		e.State = Running
	}

	return allError
}

// Start a task
func (e *EtlTask) Start() error {
	if e.State != Stopped {
		return errors.New("you can just start a stopped task")
	}

	e.StartSyncers()
	e.State = Running

	return nil
}

// Stop a task
func (e *EtlTask) Stop() error {
	if e.State != Running {
		return errors.New("you can just stop a running task")
	}

	e.StopSyncers()
	e.State = Stopped
	return nil
}

// Reset a task
func (e *EtlTask) Reset() error {
	if e.State == Executing {
		return errors.New("you can not reset a executing task")
	}

	e.Batch.Reset()

	return nil
}

// StopSyncers all syncers
func (e *EtlTask) StopSyncers() {
	for _, aSyncer := range e.Syncers {
		aSyncer.Stop()
	}
}

// StartSyncers all syncers
func (e *EtlTask) StartSyncers() {
	for _, aSyncer := range e.Syncers {
		aSyncer.Start()
	}
}
