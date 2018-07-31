package task

type TaskWrapper struct {
	log       LogFunc
	component *EtlTask
}

func NewTaskWrapper(log LogFunc, component *EtlTask) TaskWrapper {
	tw := TaskWrapper{}
	tw.log = log
	tw.component = component
	return tw
}

func (t TaskWrapper) Run() {
	if err := t.component.Execute(); err != nil {
		t.log(Error, "There are error when run task:"+err.Error())
	}
}
