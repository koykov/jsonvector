package jsonvector

// Old node.Get*ByPath() methods.
// Kept for backward compatibility.

// Look and get child object by given path and separator.
func (n *Node) GetObjectByPath(path, sep string) *Node {
	return n.GetObjectPS(path, sep)
}

// Look and get child array by given path and separator.
func (n *Node) GetArrayByPath(path, sep string) *Node {
	return n.GetArrayPS(path, sep)
}

// Look and get child bytes by given path and separator.
func (n *Node) GetBytesByPath(path, sep string) []byte {
	return n.GetBytesPS(path, sep)
}

// Look and get child string by given path and separator.
func (n *Node) GetStringByPath(path, sep string) string {
	return n.GetStringPS(path, sep)
}

// Look and get child bool by given path and separator.
func (n *Node) GetBoolByPath(path, sep string) bool {
	return n.GetBoolPS(path, sep)
}

// Look and get child float by given path and separator.
func (n *Node) GetFloatByPath(path, sep string) (float64, error) {
	return n.GetFloatPS(path, sep)
}

// Look and get child integer by given path and separator.
func (n *Node) GetIntByPath(path, sep string) (int64, error) {
	return n.GetIntPS(path, sep)
}

// Look and get child unsigned int by given path and separator.
func (n *Node) GetUintByPath(path, sep string) (uint64, error) {
	return n.GetUintPS(path, sep)
}
