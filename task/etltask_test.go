package task

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestExecuteTask_Mock_Batch struct {
	IsBegin bool
	IsReset bool
}

func (c *TestExecuteTask_Mock_Batch) GetName() string {
	return "TestExecuteTask_Mock_Batch"
}

func (c *TestExecuteTask_Mock_Batch) Begin(e *EtlTask) error {
	time.Sleep(time.Second * 2)
	c.IsBegin = true

	return nil
}

func (c *TestExecuteTask_Mock_Batch) Reset(e *EtlTask) error {
	c.IsReset = true

	return nil
}

func TestExecuteTask(t *testing.T) {
	aBatch := new(TestExecuteTask_Mock_Batch)
	aTask := NewTask(aBatch)
	aTask.ResetBeforeBegin = true
	aTask.Execute()
	assert.True(t, aBatch.IsReset, "execute should execute reset")
	assert.True(t, aBatch.IsBegin, "execute should execute begin")
	assert.Equal(t, 2, aTask.LastExecuteCost)
	assert.NotZero(t, aTask.LastExecuteTime)
}

type TestExecuteTaskErr_Mock_Batch struct {
	IsBegin bool
	IsReset bool
}

func (c *TestExecuteTaskErr_Mock_Batch) GetName() string {
	return "TestExecuteTaskErr_Mock_Batch"
}

func (c *TestExecuteTaskErr_Mock_Batch) Begin(e *EtlTask) error {
	time.Sleep(time.Microsecond * 500)
	c.IsBegin = true

	return errors.New("Test err")
}

func (c *TestExecuteTaskErr_Mock_Batch) Reset(e *EtlTask) error {
	c.IsReset = true

	return nil
}

func TestExecuteTaskErr(t *testing.T) {
	aBatch := new(TestExecuteTaskErr_Mock_Batch)
	aTask := NewTask(aBatch)
	aTask.Execute()
	assert.False(t, aBatch.IsReset, "execute should not execute reset")
	assert.True(t, aBatch.IsBegin, "execute should execute begin")
	assert.Equal(t, 0, aTask.LastExecuteCost)
	assert.Equal(t, "This task has something error when beginning, please check log", aTask.LastExecuteState)
}

type TestStatusChanged_Mock_Batch struct {
	IsChanged bool
}

func (c *TestStatusChanged_Mock_Batch) GetName() string {
	return "TestStatusChanged_Mock_Batch"
}

func (c *TestStatusChanged_Mock_Batch) Begin(e *EtlTask) error {
	time.Sleep(time.Microsecond * 500)

	return nil
}

func (c *TestStatusChanged_Mock_Batch) Reset(e *EtlTask) error {
	return nil
}

func TestStartTask_TriggerChanged(t *testing.T) {
	aBatch := new(TestStatusChanged_Mock_Batch)
	aTask := NewTask(aBatch)
	aTask.State = Stopped

	c := make(chan *EtlTask)
	aTask.StatusChanged = c
	go func() {
		for {
			select {
			case _, ok := <-c:
				aBatch.IsChanged = true
				if !ok {
					panic("you should not close startsig channel")
				}
				break
			default:
			}
		}
	}()

	err := aTask.Start()
	assert.NoError(t, err)

	time.Sleep(time.Microsecond * 10)
	assert.Equal(t, true, aBatch.IsChanged)
}
