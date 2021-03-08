package jsonvector

// Old vector.Get*ByPath() methods.
// Kept for backward compatibility.

// Look and get object by given path and separator.
func (vec *Vector) GetObjectByPath(path, sep string) *Node {
	return vec.GetObjectPS(path, sep)
}

// Look and get array by given path and separator.
func (vec *Vector) GetArrayByPath(path, sep string) *Node {
	return vec.GetArrayPS(path, sep)
}

// Look and get bytes by given path and separator.
func (vec *Vector) GetBytesByPath(path, sep string) []byte {
	return vec.GetBytesPS(path, sep)
}

// Look and get string by given path and separator.
func (vec *Vector) GetStringByPath(path, sep string) string {
	return vec.GetStringPS(path, sep)
}

// Look and get bool by given path and separator.
func (vec *Vector) GetBoolByPath(path, sep string) bool {
	return vec.GetBoolPS(path, sep)
}

// Look and get float by given path and separator.
func (vec *Vector) GetFloatByPath(path, sep string) (float64, error) {
	return vec.GetFloatPS(path, sep)
}

// Look and get integer by given path and separator.
func (vec *Vector) GetIntByPath(path, sep string) (int64, error) {
	return vec.GetIntPS(path, sep)
}

// Look and get unsigned integer by given path and separator.
func (vec *Vector) GetUintByPath(path, sep string) (uint64, error) {
	return vec.GetUintPS(path, sep)
}
