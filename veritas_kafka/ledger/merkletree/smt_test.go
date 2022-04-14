package merkletree

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicSparseMerkleTree(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "test-basic")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	store, err := NewBadgerStore(filepath.Join(tmpDir, "bs"))
	require.NoError(t, err)
	tree := NewSparseMerkleTree(store, sha256.New())

	_, err = tree.Update([]byte("foo"), []byte("bar"))
	require.NoError(t, err)

	proof, _ := tree.Prove([]byte("foo"))
	root := tree.Root()

	if VerifyProof(proof, root, []byte("foo"), []byte("bar"), sha256.New()) {
		fmt.Println("Proof verification succeeded.")
	} else {
		fmt.Println("Proof verification failed.")
	}
}
