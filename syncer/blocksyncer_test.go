package syncer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSigCountShouldBeNone_AfterStop(t *testing.T) {
	aBlockSyncer := &BlockSyncer{}

	aBlockSyncer.Start()
	aBlockSyncer.Stop()

	assert.Len(t, aBlockSyncer.startSig, 0, "starg sig should be none after stop")
}

func TestCheckShouldNotBeBlocked_AfterStart(t *testing.T) {
	aBlockSyncer := &BlockSyncer{}

	aBlockSyncer.Start()
	aBlockSyncer.CheckStatus()
}
