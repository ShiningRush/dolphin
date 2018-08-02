package syncer

// BlockSyncer you should compose this struct insdead of using it
type BlockSyncer struct {
	isStopped bool
	startSig  chan bool
}

// Start syncer so that can continue job
func (t *BlockSyncer) Start() {
	t.ifNotExsistChanCreateIt()
	t.isStopped = false
	t.startSig <- true
}

// Stop syncer to avoid missing changes
func (t *BlockSyncer) Stop() {
	t.ifNotExsistChanCreateIt()
	select {
	case _, ok := <-t.startSig:
		if !ok {
			panic("you should not close startsig channel")
		}
	default:
	}
	t.isStopped = true
}

// CheckStatus will block until status is change to running
func (t *BlockSyncer) CheckStatus() {
	if !t.isStopped {
		return
	}
	<-t.startSig
}

func (t *BlockSyncer) ifNotExsistChanCreateIt() {
	if t.startSig == nil {
		t.startSig = make(chan bool, 1)
	}
}
