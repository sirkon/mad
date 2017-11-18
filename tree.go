package mad

// TreeNodeData represents tree node data
type TreeNodeData struct {
	Node     *TreeNode
	Code     *code
	Comment  *comment
	Integer  *integer
	Unsigned *unsigned
	Float    *float
	String   *String
	Boolean  *boolean
}

// TreeNode represents tree node:
//   * It must has a header
//   * It probably has related data
type TreeNode struct {
	Header header
	Data   []TreeNodeData
}

// Tree consists of tree nodes, i.e. there must be a header before the content or just a valid `raw` fenced code block
type Tree struct {
	Node []TreeNode
}
