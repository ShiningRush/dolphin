package task

type LogLevel string

const (
	Debug LogLevel = "Debug"
	Info  LogLevel = "Info"
	Error LogLevel = "Error"
)

// TaskType task type
type TaskType uint

// task type
const (
	// OneShot this task will run just one time
	OneShot TaskType = iota
	// Plan when you set this value, you also need to indicate its plan time
	Plan
)

// ToString transform to enum name
func (tt *TaskType) ToString() string {
	switch *tt {
	case OneShot:
		return "OneShot"
	case Plan:
		return "Plan"
	}
	return ""
}

// TaskState task state
type TaskState uint

// task state
const (
	// Init task default state
	Init TaskState = iota
	// Executing as it is
	Executing

	// Running as it is
	Running

	// Commplete as it is
	Completed

	// Stopped as it is
	Stopped
)

// ToString transform to enum name
func (ts *TaskState) ToString() string {
	switch *ts {
	case Init:
		return "Init"
	case Executing:
		return "Executing"
	case Running:
		return "Running"
	case Completed:
		return "Completed"
	case Stopped:
		return "Stopped"
	}

	return ""
}
