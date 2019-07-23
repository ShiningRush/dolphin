package dolphin

import (
	"github.com/shiningrush/dolphin/task"
)

var tm *task.TaskManager

func init() {
	tm = task.NewTaskManager()
}

func SetLog(log task.LogFunc) {
	tm.SetLog(log)
}

func SetRepo(repo task.TaskStatusRepository) {
	tm.SetRepo(repo)
}

// GetAllTasks in manager
func GetAllTasks() map[string]*task.EtlTask {
	return tm.GetAllTasks()
}

// Add a task to manager
func Add(aTask *task.EtlTask) *task.TaskManager {
	return tm.Add(aTask)
}

// Reset a task
func Reset(taskName string) error {
	return tm.Reset(taskName)
}

// Start a task
func Start(taskName string) error {
	return tm.Start(taskName)
}

// Stop a task
func Stop(taskName string) error {
	return tm.Stop(taskName)
}

// Execute a task
func Execute(taskName string) error {
	return tm.Execute(taskName)
}
