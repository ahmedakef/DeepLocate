package structures

import (
	"math"
	"time"

	utils "dlocate/osutils"
)

const k uint = 4

//KDTree is a K dimentional binary tree (range search)
type KDTree struct {
	Root *node
}

//Insert adds a value (node) to the tree
func (tree *KDTree) Insert(value *utils.FileMetadata) {
	tree.Root = insertNode(tree.Root, value, 0)
}

//Search returns a list of all the valid nodes in the tree
func (tree *KDTree) Search(start *utils.FileMetadata, end *utils.FileMetadata) []*utils.FileMetadata {
	return searchRange(tree.Root, start, end, 0)
}

//SearchPartial will fill the other search paramters given a partial range to search into,
func (tree *KDTree) SearchPartial(start *utils.FileMetadata, end *utils.FileMetadata) []*utils.FileMetadata {
	if end.Size == 0 {
		end.Size = math.MaxInt64
	}
	if start.CTime.IsZero() {
		start.CTime = start.MTime
	}
	if start.ATime.IsZero() {
		start.ATime = start.CTime
	}
	if end.ATime.IsZero() {
		//add time to handle search process time
		end.ATime = time.Now().Add(time.Hour * 1)
	}
	if end.CTime.IsZero() {
		//add time to handle search process time
		end.CTime = end.ATime
	}
	if end.MTime.IsZero() {
		//add time to handle search process time
		end.MTime = end.CTime
	}
	return tree.Search(start, end)
}

type node struct {
	Value       utils.FileMetadata
	Left, Right *node
}

func newNode(val *utils.FileMetadata) *node {
	var temp = &node{*val, nil, nil}
	return temp
}

func insertNode(root *node, value *utils.FileMetadata, depth uint) *node {
	if root == nil {
		return newNode(value)
	}

	var cd uint = depth % k

	if compareNodes(value, &root.Value, cd) {
		root.Left = insertNode(root.Left, value, depth+1)
	} else {
		root.Right = insertNode(root.Right, value, depth+1)
	}

	return root
}

func searchRange(root *node, s *utils.FileMetadata, e *utils.FileMetadata, depth uint) []*utils.FileMetadata {
	results := make([]*utils.FileMetadata, 0)

	if root == nil {
		return results
	}

	var cd uint = depth % k

	if compareNodes(s, &root.Value, cd) {
		results = append(results, searchRange(root.Left, s, e, depth+1)...)
	}

	if compareNodes(&root.Value, e, cd) {
		results = append(results, searchRange(root.Right, s, e, depth+1)...)
	}

	if matchNode(&root.Value, s, e) {
		results = append(results, &root.Value)
	}

	return results
}

func matchNode(v *utils.FileMetadata, s *utils.FileMetadata, e *utils.FileMetadata) bool {
	if v.MTime.Before(s.MTime) || v.MTime.After(e.MTime) {
		return false
	}
	if v.CTime.Before(s.CTime) || v.CTime.After(e.CTime) {
		return false
	}
	if v.ATime.Before(s.ATime) || v.ATime.After(e.ATime) {
		return false
	}
	if v.Size < s.Size || v.Size > e.Size {
		return false
	}

	return true
}

// compareNodes is a function to compare two nodes' values
func compareNodes(v1 *utils.FileMetadata, v2 *utils.FileMetadata, dim uint) bool {
	switch dim {
	case 0:
		return v1.MTime.Before(v2.MTime)
	case 1:
		return v1.CTime.Before(v2.CTime)
	case 2:
		return v1.ATime.Before(v2.ATime)
	case 3:
		return v1.Size < v2.Size
	default:
		return false
	}
}
