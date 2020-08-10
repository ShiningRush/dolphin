package task

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testBatch struct {
}

func (c *testBatch) GetName() string {
	return "testBatch"
}

func (c *testBatch) Begin(e *EtlTask) error {
	fmt.Println("test")

	return nil
}

func (c *testBatch) Reset(e *EtlTask) error {
	return nil
}

type testBatchB struct {
	testBatch
}

func (c *testBatchB) GetName() string {
	return "testBatchB"
}

func TestGetAllTask(t *testing.T) {
	aBatch, bBatch := new(testBatch), new(testBatchB)
	aTask := NewTask(aBatch)
	bTask := NewTask(bBatch)

	tm := NewTaskManager()
	tm.Add(aTask).Add(bTask)
	assert.Len(t, tm.GetAllTasks(), 2, "Should have 2 tasks")
}

func TestAddBatch(t *testing.T) {
	var aBatch Batch
	aBatch = new(testBatch)
	aTask := NewTask(aBatch)

	tm := NewTaskManager()
	tm.Add(aTask)
}

func TestAddRepeatBatch(t *testing.T) {
	var aBatch Batch
	aBatch = new(testBatch)
	aTask := NewTask(aBatch)

	tm := NewTaskManager()
	tm.Add(aTask)
	assert.Panics(t, func() { tm.Add(aTask) }, "Should happened panic")
}

func TestBuildManager(t *testing.T) {
	var aBatch Batch
	aBatch = new(testBatch)
	aTask := NewTask(aBatch)

	bTask := NewTask(&testBatchB{})
	bTask.Type = Plan
	bTask.PlanTime = "@every 1s"

	tm := NewTaskManager()
	tm.Add(aTask).Add(bTask).Build()
}
