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

func (c *TestExecuteTask_Mock_Batch) Begin() error {
	time.Sleep(time.Second * 2)
	c.IsBegin = true

	return nil
}

func (c *TestExecuteTask_Mock_Batch) Reset() error {
	c.IsReset = true

	return nil
}

type TestExecuteTaskErr_Mock_Batch struct {
	IsBegin bool
	IsReset bool
}

func (c *TestExecuteTaskErr_Mock_Batch) GetName() string {
	return "TestExecuteTaskErr_Mock_Batch"
}

func (c *TestExecuteTaskErr_Mock_Batch) Begin() error {
	time.Sleep(time.Second * 2)
	c.IsBegin = true

	return errors.New("Test err")
}

func (c *TestExecuteTaskErr_Mock_Batch) Reset() error {
	c.IsReset = true

	return nil
}

func TestExecuteTask(t *testing.T) {
	aBatch := new(TestExecuteTask_Mock_Batch)
	aTask := NewTask(aBatch)
	aTask.Execute()
	assert.True(t, aBatch.IsReset, "execute should execute reset")
	assert.True(t, aBatch.IsBegin, "execute should execute begin")
	assert.Equal(t, 2, aTask.LastExecuteCost)
	assert.NotZero(t, aTask.LastExecuteTime)
}

func TestExecuteTaskErr(t *testing.T) {
	aBatch := new(TestExecuteTaskErr_Mock_Batch)
	aTask := NewTask(aBatch)
	aTask.Execute()
	assert.True(t, aBatch.IsReset, "execute should execute reset")
	assert.True(t, aBatch.IsBegin, "execute should execute begin")
	assert.Equal(t, 2, aTask.LastExecuteCost)
	assert.Equal(t, "This task has something error when beginning, please check log", aTask.LastExecuteState)
}
