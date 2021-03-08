package jsonvector

// Old node.Get*ByPath() methods.
// Kept for backward compatibility.

// Look and get child object by given path and separator.
func (n *Node) GetObjectByPath(path string) *Node {
	return n.GetObjectPS(path, ".")
}

// Look and get child array by given path and separator.
func (n *Node) GetArrayByPath(path string) *Node {
	return n.GetArrayPS(path, ".")
}

// Look and get child bytes by given path and separator.
func (n *Node) GetBytesByPath(path string) []byte {
	return n.GetBytesPS(path, ".")
}

// Look and get child string by given path and separator.
func (n *Node) GetStringByPath(path string) string {
	return n.GetStringPS(path, ".")
}

// Look and get child bool by given path and separator.
func (n *Node) GetBoolByPath(path string) bool {
	return n.GetBoolPS(path, ".")
}

// Look and get child float by given path and separator.
func (n *Node) GetFloatByPath(path string) (float64, error) {
	return n.GetFloatPS(path, ".")
}

// Look and get child integer by given path and separator.
func (n *Node) GetIntByPath(path string) (int64, error) {
	return n.GetIntPS(path, ".")
}

// Look and get child unsigned int by given path and separator.
func (n *Node) GetUintByPath(path string) (uint64, error) {
	return n.GetUintPS(path, ".")
}
