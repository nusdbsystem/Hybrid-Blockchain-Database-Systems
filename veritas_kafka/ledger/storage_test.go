package ledger

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateDigestComputation(t *testing.T) {
	dir, err := ioutil.TempDir("", "prefix1")
	if err != nil {
		log.Fatal(err)
	}
	ledger1, err := NewLedger(dir, true)
	if err != nil {
		log.Fatal(err)
	}
	ledger1.Append([]byte("A"), []byte("A"))
	ledger1.Append([]byte("B"), []byte("B"))
	ledger1.Append([]byte("C"), []byte("C"))
	ledger1Digest := ledger1.GetRootDigest()
	ledger1.Close()

	// Create another instance to test whether previous appended states are loaded properly from the persistence and reach the identical root digest.
	ledger2, err := NewLedger(dir, true)
	if err != nil {
		log.Fatal(err)
	}
	ledger2Digest := ledger2.GetRootDigest()
	assert.Equal(t, ledger1Digest, ledger2Digest)

	ledger2.AppendBlk([]byte("block111"))
	ledger2.AppendBlk([]byte("block222"))
	ledger2.AppendBlk([]byte("block333"))
	ledger2AfterAppendingBlkDigest := ledger2.GetRootDigest()
	// Verify the block appending does not interfere the state digest computation.
	assert.Equal(t, ledger2AfterAppendingBlkDigest, ledger2Digest)
	ledger2.Close()

	// Verify the non-interference after reloading db.
	ledger3, err := NewLedger(dir, true)
	if err != nil {
		log.Fatal(err)
	}
	ledger3Digest := ledger3.GetRootDigest()
	assert.Equal(t, ledger3Digest, ledger2AfterAppendingBlkDigest)
	ledger3.Close()
	defer os.RemoveAll(dir)
}
