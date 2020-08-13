package task

type TaskWrapper struct {
	log       LogFunc
	component *EtlTask
}

// NewTaskWrapper return a wrapper for cron
func NewTaskWrapper(log LogFunc, component *EtlTask) TaskWrapper {
	tw := TaskWrapper{}
	tw.log = log
	tw.component = component
	return tw
}

// Run it is for cron
func (t TaskWrapper) Run() {
	if err := t.component.Execute(); err != nil {
		t.log(Error, "there are error when run task:"+err.Error())
	}
}
