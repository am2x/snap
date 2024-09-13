package main

// invert a binary tree
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// invert tree inverts a binary tree
func (root *TreeNode) ReverseTree() {
	if root == nil {
		return
	}

	//swap the left and right children
	root.Left, root.Right = root.Right, root.Left

	//recursively invert the left and right subtrees
	root.Left.ReverseTree()
	root.Right.ReverseTree()
}
