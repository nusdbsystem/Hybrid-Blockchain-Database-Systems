// Package smt implements a Sparse Merkle tree.
package merkletree

import (
	"bytes"
	"hash"
)

const (
	right = 1
)

var defaultValue []byte

type keyAlreadyEmptyError struct{}

func (e *keyAlreadyEmptyError) Error() string {
	return "key already empty"
}

type SparseMerkleTree struct {
	th   treeHasher
	kvs  KVStore
	root []byte
}

func NewSparseMerkleTree(kvs KVStore, hasher hash.Hash) *SparseMerkleTree {
	smt := SparseMerkleTree{
		th:  *newTreeHasher(hasher),
		kvs: kvs,
	}

	smt.SetRoot(smt.th.placeholder())

	return &smt
}

func ImportSparseMerkleTree(kvs KVStore, hasher hash.Hash, root []byte) *SparseMerkleTree {
	smt := SparseMerkleTree{
		th:   *newTreeHasher(hasher),
		kvs:  kvs,
		root: root,
	}
	return &smt
}

func (smt *SparseMerkleTree) Root() []byte {
	return smt.root
}

func (smt *SparseMerkleTree) SetRoot(root []byte) {
	smt.root = root
}

func (smt *SparseMerkleTree) depth() int {
	return smt.th.pathSize() * 8
}

func (smt *SparseMerkleTree) Get(key []byte) ([]byte, error) {
	value, err := smt.GetForRoot(key, smt.Root())
	return value, err
}

func (smt *SparseMerkleTree) GetForRoot(key []byte, root []byte) ([]byte, error) {
	if bytes.Equal(root, smt.th.placeholder()) {
		// The tree is empty, return the default value.
		return defaultValue, nil
	}

	path := smt.th.path(key)
	currentHash := root
	for i := 0; i < smt.depth(); i++ {
		currentData, err := smt.kvs.Get(currentHash)
		if err != nil {
			return nil, err
		} else if smt.th.isLeaf(currentData) {
			// We've reached the end. Is this the actual leaf?
			p, valueHash := smt.th.parseLeaf(currentData)
			if !bytes.Equal(path, p) {
				// Nope. Therefore the key is actually empty.
				return defaultValue, nil
			}
			// Otherwise, yes. Return the value.
			value, err := smt.kvs.Get(valueHash)
			if err != nil {
				return nil, err
			}
			return value, nil
		}

		leftNode, rightNode := smt.th.parseNode(currentData)
		if hasBit(path, i) == right {
			currentHash = rightNode
		} else {
			currentHash = leftNode
		}

		if bytes.Equal(currentHash, smt.th.placeholder()) {
			// We've hit a placeholder value; this is the end.
			return defaultValue, nil
		}
	}
	currentData, err := smt.kvs.Get(currentHash)
	if err != nil {
		return nil, err
	}
	_, valueHash := smt.th.parseLeaf(currentData)
	value, err := smt.kvs.Get(valueHash)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (smt *SparseMerkleTree) Has(key []byte) (bool, error) {
	val, err := smt.Get(key)
	return !bytes.Equal(defaultValue, val), err
}

func (smt *SparseMerkleTree) HasForRoot(key, root []byte) (bool, error) {
	val, err := smt.GetForRoot(key, root)
	return !bytes.Equal(defaultValue, val), err
}

func (smt *SparseMerkleTree) Update(key []byte, value []byte) ([]byte, error) {
	newRoot, err := smt.UpdateForRoot(key, value, smt.Root())
	if err == nil {
		smt.SetRoot(newRoot)
	}
	return newRoot, err
}

func (smt *SparseMerkleTree) Delete(key []byte) ([]byte, error) {
	return smt.Update(key, defaultValue)
}

func (smt *SparseMerkleTree) UpdateForRoot(key []byte, value []byte, root []byte) ([]byte, error) {
	path := smt.th.path(key)
	sideNodes, oldLeafHash, oldLeafData, err := smt.sideNodesForRoot(path, root)
	if err != nil {
		return nil, err
	}

	var newRoot []byte
	if bytes.Equal(value, defaultValue) {
		newRoot, err = smt.deleteWithSideNodes(path, sideNodes, oldLeafHash, oldLeafData)
		if _, ok := err.(*keyAlreadyEmptyError); ok {
			// This key is already empty; return the old root.
			return root, nil
		}
	} else {
		newRoot, err = smt.updateWithSideNodes(path, value, sideNodes, oldLeafHash, oldLeafData)
	}
	return newRoot, err
}

func (smt *SparseMerkleTree) DeleteForRoot(key, root []byte) ([]byte, error) {
	return smt.UpdateForRoot(key, defaultValue, root)
}

func (smt *SparseMerkleTree) deleteWithSideNodes(path []byte, sideNodes [][]byte, oldLeafHash []byte, oldLeafData []byte) ([]byte, error) {
	if bytes.Equal(oldLeafHash, smt.th.placeholder()) {
		// This key is already empty as it is a placeholder; return an error.
		return nil, &keyAlreadyEmptyError{}
	} else if actualPath, _ := smt.th.parseLeaf(oldLeafData); !bytes.Equal(path, actualPath) {
		return nil, &keyAlreadyEmptyError{}
	}

	var currentHash, currentData []byte
	nonPlaceholderReached := false
	for i := smt.depth() - 1; i >= 0; i-- {
		if sideNodes[i] == nil {
			continue
		}

		sideNode := make([]byte, smt.th.pathSize())
		copy(sideNode, sideNodes[i])

		if currentData == nil {
			sideNodeValue, err := smt.kvs.Get(sideNode)
			if err != nil {
				return nil, err
			}

			if smt.th.isLeaf(sideNodeValue) {
				currentHash = sideNode
				currentData = sideNode
				continue
			} else {
				currentData = smt.th.placeholder()
				nonPlaceholderReached = true
			}
		}

		if !nonPlaceholderReached && bytes.Equal(sideNode, smt.th.placeholder()) {
			continue
		} else if !nonPlaceholderReached {
			nonPlaceholderReached = true
		}

		if hasBit(path, i) == right {
			currentHash, currentData = smt.th.digestNode(sideNode, currentData)
		} else {
			currentHash, currentData = smt.th.digestNode(currentData, sideNode)
		}
		err := smt.kvs.Set(currentHash, currentData)
		if err != nil {
			return nil, err
		}
		currentData = currentHash
	}

	if currentHash == nil {
		currentHash = smt.th.placeholder()
	}
	return currentHash, nil
}

func (smt *SparseMerkleTree) updateWithSideNodes(path []byte, value []byte, sideNodes [][]byte, oldLeafHash []byte, oldLeafData []byte) ([]byte, error) {
	valueHash := smt.th.digest(value)
	if err := smt.kvs.Set(valueHash, value); err != nil {
		return nil, err
	}

	currentHash, currentData := smt.th.digestLeaf(path, valueHash)
	if err := smt.kvs.Set(currentHash, currentData); err != nil {
		return nil, err
	}
	currentData = currentHash
	var commonPrefixCount int
	if bytes.Equal(oldLeafHash, smt.th.placeholder()) {
		commonPrefixCount = smt.depth()
	} else {
		actualPath, _ := smt.th.parseLeaf(oldLeafData)
		commonPrefixCount = countCommonPrefix(path, actualPath)
	}
	if commonPrefixCount != smt.depth() {
		if hasBit(path, commonPrefixCount) == right {
			currentHash, currentData = smt.th.digestNode(oldLeafHash, currentData)
		} else {
			currentHash, currentData = smt.th.digestNode(currentData, oldLeafHash)
		}

		err := smt.kvs.Set(currentHash, currentData)
		if err != nil {
			return nil, err
		}

		currentData = currentHash
	}

	for i := smt.depth() - 1; i >= 0; i-- {
		sideNode := make([]byte, smt.th.pathSize())

		if sideNodes[i] == nil {
			if commonPrefixCount != smt.depth() && commonPrefixCount > i {
				copy(sideNode, smt.th.placeholder())
			} else {
				continue
			}
		} else {
			copy(sideNode, sideNodes[i])
		}

		if hasBit(path, i) == right {
			currentHash, currentData = smt.th.digestNode(sideNode, currentData)
		} else {
			currentHash, currentData = smt.th.digestNode(currentData, sideNode)
		}
		err := smt.kvs.Set(currentHash, currentData)
		if err != nil {
			return nil, err
		}
		currentData = currentHash
	}

	return currentHash, nil
}

func (smt *SparseMerkleTree) sideNodesForRoot(path []byte, root []byte) ([][]byte, []byte, []byte, error) {
	sideNodes := make([][]byte, smt.depth())

	if bytes.Equal(root, smt.th.placeholder()) {
		// If the root is a placeholder, there are no sidenodes to return.
		// Let the "actual path" be the input path.
		return sideNodes, smt.th.placeholder(), nil, nil
	}

	currentData, err := smt.kvs.Get(root)
	if err != nil {
		return nil, nil, nil, err
	} else if smt.th.isLeaf(currentData) {
		// If the root is a leaf, there are also no sidenodes to return.
		return sideNodes, root, currentData, nil
	}

	var nodeHash []byte
	for i := 0; i < smt.depth(); i++ {
		leftNode, rightNode := smt.th.parseNode(currentData)

		if hasBit(path, i) == right {
			sideNodes[i] = leftNode
			nodeHash = rightNode
		} else {
			sideNodes[i] = rightNode
			nodeHash = leftNode
		}

		if bytes.Equal(nodeHash, smt.th.placeholder()) {
			// If the node is a placeholder, we've reached the end.
			return sideNodes, nodeHash, nil, nil
		}

		currentData, err = smt.kvs.Get(nodeHash)
		if err != nil {
			return nil, nil, nil, err
		} else if smt.th.isLeaf(currentData) {
			break
		}
	}

	return sideNodes, nodeHash, currentData, err
}

func (smt *SparseMerkleTree) Prove(key []byte) (SparseMerkleProof, error) {
	proof, err := smt.ProveForRoot(key, smt.Root())
	return proof, err
}

func (smt *SparseMerkleTree) ProveForRoot(key []byte, root []byte) (SparseMerkleProof, error) {
	path := smt.th.path(key)
	sideNodes, leafHash, leafData, err := smt.sideNodesForRoot(path, root)
	if err != nil {
		return SparseMerkleProof{}, err
	}

	var nonEmptySideNodes [][]byte
	for _, v := range sideNodes {
		if v != nil {
			nonEmptySideNodes = append(nonEmptySideNodes, v)
		}
	}

	var nonMembershipLeafData []byte
	if !bytes.Equal(leafHash, smt.th.placeholder()) {
		nonMembershipLeafData = leafData
	}

	proof := SparseMerkleProof{
		SideNodes:             nonEmptySideNodes,
		NonMembershipLeafData: nonMembershipLeafData,
	}

	return proof, err
}

func (smt *SparseMerkleTree) ProveCompact(key []byte) (SparseCompactMerkleProof, error) {
	proof, err := smt.ProveCompactForRoot(key, smt.Root())
	return proof, err
}

func (smt *SparseMerkleTree) ProveCompactForRoot(key []byte, root []byte) (SparseCompactMerkleProof, error) {
	proof, err := smt.ProveForRoot(key, root)
	if err != nil {
		return SparseCompactMerkleProof{}, err
	}
	compactedProof, err := CompactProof(proof, smt.th.hasher)
	return compactedProof, err
}
