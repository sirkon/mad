package mad

// TreeNodeData represents tree node data
type TreeNodeData struct {
	Node     *TreeNode
	Code     *Code
	Comment  *Comment
	Integer  *Integer
	Unsigned *Unsigned
	Float    *Float
	String   *String
	Boolean  *Boolean
}

// TreeNode represents tree node:
//   * It must has a header
//   * It probably has related data
type TreeNode struct {
	Header Header
	Data   []TreeNodeData
}

// Tree consists of tree nodes, i.e. there must be a header before the content or just a valid `raw` fenced code block
type Tree struct {
	Node []TreeNode
}
