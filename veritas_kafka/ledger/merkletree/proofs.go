package merkletree

import (
	"bytes"
	"errors"
	"hash"
	"math"
)

type SparseMerkleProof struct {
	SideNodes             [][]byte
	NonMembershipLeafData []byte
}

func (proof *SparseMerkleProof) sanityCheck(th *treeHasher) bool {
	if len(proof.SideNodes) > th.pathSize()*8 ||
		(proof.NonMembershipLeafData != nil && len(proof.NonMembershipLeafData) != len(leafPrefix)+th.pathSize()+th.hasher.Size()) {
		return false
	}
	for _, v := range proof.SideNodes {
		if len(v) != th.hasher.Size() {
			return false
		}
	}

	return true
}

type SparseCompactMerkleProof struct {
	SideNodes             [][]byte
	NonMembershipLeafData []byte
	BitMask               []byte
	NumSideNodes          int
}

func (proof *SparseCompactMerkleProof) sanityCheck(th *treeHasher) bool {
	if proof.NumSideNodes < 0 || proof.NumSideNodes > th.pathSize()*8 ||
		len(proof.BitMask) != int(math.Ceil(float64(proof.NumSideNodes)/float64(8))) ||
		(proof.NumSideNodes > 0 && len(proof.SideNodes) != proof.NumSideNodes-countSetBits(proof.BitMask)) {
		return false
	}

	return true
}

func VerifyProof(proof SparseMerkleProof, root []byte, key []byte, value []byte, hasher hash.Hash) bool {
	result, _ := verifyProofWithUpdates(proof, root, key, value, hasher)
	return result
}

func verifyProofWithUpdates(proof SparseMerkleProof, root []byte, key []byte, value []byte, hasher hash.Hash) (bool, [][][]byte) {
	th := newTreeHasher(hasher)
	path := th.path(key)

	if !proof.sanityCheck(th) {
		return false, nil
	}

	var updates [][][]byte
	var currentHash, currentData []byte
	if bytes.Equal(value, defaultValue) {
		if proof.NonMembershipLeafData == nil {
			currentHash = th.placeholder()
		} else {
			actualPath, valueHash := th.parseLeaf(proof.NonMembershipLeafData)
			if bytes.Equal(actualPath, path) {
				// This is not an unrelated leaf; non-membership proof failed.
				return false, nil
			}
			currentHash, currentData = th.digestLeaf(actualPath, valueHash)

			update := make([][]byte, 2)
			update[0], update[1] = currentHash, currentData
			updates = append(updates, update)
		}
	} else {
		valueHash := th.digest(value)
		update := make([][]byte, 2)
		update[0], update[1] = valueHash, value
		updates = append(updates, update)

		currentHash, currentData = th.digestLeaf(path, valueHash)
		update = make([][]byte, 2)
		update[0], update[1] = currentHash, currentData
		updates = append(updates, update)
	}

	for i := len(proof.SideNodes) - 1; i >= 0; i-- {
		node := make([]byte, th.pathSize())
		copy(node, proof.SideNodes[i])

		if hasBit(path, i) == right {
			currentHash, currentData = th.digestNode(node, currentHash)
		} else {
			currentHash, currentData = th.digestNode(currentHash, node)
		}

		update := make([][]byte, 2)
		update[0], update[1] = currentHash, currentData
		updates = append(updates, update)
	}

	return bytes.Equal(currentHash, root), updates
}

func VerifyCompactProof(proof SparseCompactMerkleProof, root []byte, key []byte, value []byte, hasher hash.Hash) bool {
	decompactedProof, err := DecompactProof(proof, hasher)
	if err != nil {
		return false
	}
	return VerifyProof(decompactedProof, root, key, value, hasher)
}

func CompactProof(proof SparseMerkleProof, hasher hash.Hash) (SparseCompactMerkleProof, error) {
	th := newTreeHasher(hasher)

	if !proof.sanityCheck(th) {
		return SparseCompactMerkleProof{}, errors.New("bad proof")
	}

	bitMask := emptyBytes(int(math.Ceil(float64(len(proof.SideNodes)) / float64(8))))
	var compactedSideNodes [][]byte
	for i := 0; i < len(proof.SideNodes); i++ {
		node := make([]byte, th.hasher.Size())
		copy(node, proof.SideNodes[i])
		if bytes.Equal(node, th.placeholder()) {
			setBit(bitMask, i)
		} else {
			compactedSideNodes = append(compactedSideNodes, node)
		}
	}

	return SparseCompactMerkleProof{
		SideNodes:             compactedSideNodes,
		NonMembershipLeafData: proof.NonMembershipLeafData,
		BitMask:               bitMask,
		NumSideNodes:          len(proof.SideNodes),
	}, nil
}

func DecompactProof(proof SparseCompactMerkleProof, hasher hash.Hash) (SparseMerkleProof, error) {
	th := newTreeHasher(hasher)

	if !proof.sanityCheck(th) {
		return SparseMerkleProof{}, errors.New("bad proof")
	}

	decompactedSideNodes := make([][]byte, proof.NumSideNodes)
	position := 0
	for i := 0; i < proof.NumSideNodes; i++ {
		if hasBit(proof.BitMask, i) == 1 {
			decompactedSideNodes[i] = th.placeholder()
		} else {
			decompactedSideNodes[i] = proof.SideNodes[position]
			position++
		}
	}

	return SparseMerkleProof{
		SideNodes:             decompactedSideNodes,
		NonMembershipLeafData: proof.NonMembershipLeafData,
	}, nil
}
