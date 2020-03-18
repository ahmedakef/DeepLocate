package structures

import (
	utils "../osutils"
)

const k uint = 5

//KDTree is a K dimentional binary tree (range search)
type KDTree struct {
	root *node
}

//Insert adds a value (node) to the tree
func (tree *KDTree) Insert(value *utils.FileMetadata) {
	tree.root = insertNode(tree.root, value, 0)
}

//Search adds a value (node) to the tree
func (tree *KDTree) Search(start *utils.FileMetadata, end *utils.FileMetadata) []*utils.FileMetadata {
	return searchRange(tree.root, start, end, 0)
}

type node struct {
	value       *utils.FileMetadata
	left, right *node
}

func newNode(val *utils.FileMetadata) *node {
	temp := &node{val, nil, nil}
	return temp
}

func insertNode(root *node, value *utils.FileMetadata, depth uint) *node {
	if root == nil {
		return newNode(value)
	}

	var cd uint = depth % k

	if compareNodes(value, root.value, cd) {
		root.left = insertNode(root.left, value, depth+1)
	} else {
		root.right = insertNode(root.right, value, depth+1)
	}

	return root
}

func searchRange(root *node, s *utils.FileMetadata, e *utils.FileMetadata, depth uint) []*utils.FileMetadata {
	results := make([]*utils.FileMetadata, 0)

	if root == nil {
		return results
	}

	var cd uint = depth % k

	if compareNodes(s, root.value, cd) {
		results = append(results, searchRange(root.left, s, e, depth+1)...)
	}

	if compareNodes(root.value, e, cd) {
		results = append(results, searchRange(root.right, s, e, depth+1)...)
	}

	return results
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
	case 4:
		return v1.Extension < v2.Extension
	default:
		return false
	}
}
