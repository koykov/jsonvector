// Getter API of JSON node.
package jsonvector

// // Look and get child object by given keys.
// func (n *Node) GetObject(keys ...string) *Node {
// 	c := n.Get(keys...)
// 	if c == nil || c.Type() != TypeObj {
// 		return nil
// 	}
// 	return c.Object()
// }
//
// // Look and get child array by given keys.
// func (n *Node) GetArray(keys ...string) *Node {
// 	c := n.Get(keys...)
// 	if c == nil || c.Type() != TypeArr {
// 		return nil
// 	}
// 	return c.Array()
// }
//
// // Look and get child bytes by given keys.
// func (n *Node) GetBytes(keys ...string) []byte {
// 	c := n.Get(keys...)
// 	if c == nil || c.Type() != TypeStr {
// 		return nil
// 	}
// 	return c.Bytes()
// }
//
// // Look and get child string by given keys.
// func (n *Node) GetString(keys ...string) string {
// 	c := n.Get(keys...)
// 	if c == nil || c.Type() != TypeStr {
// 		return ""
// 	}
// 	return c.String()
// }
//
// // Look and get child bool by given keys.
// func (n *Node) GetBool(keys ...string) bool {
// 	c := n.Get(keys...)
// 	if c == nil || c.Type() != TypeBool {
// 		return false
// 	}
// 	return c.Bool()
// }
//
// // Look and get child float by given keys.
// func (n *Node) GetFloat(keys ...string) (float64, error) {
// 	c := n.Get(keys...)
// 	if c == nil {
// 		return 0, ErrNotFound
// 	}
// 	if c.Type() != TypeNum {
// 		return 0, ErrIncompatType
// 	}
// 	return c.Float()
// }
//
// // Look and get child integer by given keys.
// func (n *Node) GetInt(keys ...string) (int64, error) {
// 	c := n.Get(keys...)
// 	if c == nil {
// 		return 0, ErrNotFound
// 	}
// 	if c.Type() != TypeNum {
// 		return 0, ErrIncompatType
// 	}
// 	return c.Int()
// }
//
// // Look and get child unsigned integer by given keys.
// func (n *Node) GetUint(keys ...string) (uint64, error) {
// 	c := n.Get(keys...)
// 	if c == nil {
// 		return 0, ErrNotFound
// 	}
// 	if c.Type() != TypeNum {
// 		return 0, ErrIncompatType
// 	}
// 	return c.Uint()
// }
//
// // Look and get child object by given path and separator.
// func (n *Node) GetObjectPS(path, sep string) *Node {
// 	vec := n.vec()
// 	if vec == nil {
// 		return nil
// 	}
// 	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
// 	c := n.Get(vec.ss...)
// 	if c == nil || c.Type() != TypeObj {
// 		return nil
// 	}
// 	return c.Object()
// }
//
// // Look and get child array by given path and separator.
// func (n *Node) GetArrayPS(path, sep string) *Node {
// 	vec := n.vec()
// 	if vec == nil {
// 		return nil
// 	}
// 	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
// 	c := n.Get(vec.ss...)
// 	if c == nil || c.Type() != TypeArr {
// 		return nil
// 	}
// 	return c.Array()
// }
//
// // Look and get child bytes by given path and separator.
// func (n *Node) GetBytesPS(path, sep string) []byte {
// 	vec := n.vec()
// 	if vec == nil {
// 		return nil
// 	}
// 	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
// 	c := n.Get(vec.ss...)
// 	if c == nil || c.Type() != TypeStr {
// 		return nil
// 	}
// 	return c.Bytes()
// }
//
// // Look and get child string by given path and separator.
// func (n *Node) GetStringPS(path, sep string) string {
// 	vec := n.vec()
// 	if vec == nil {
// 		return ""
// 	}
// 	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
// 	c := n.Get(vec.ss...)
// 	if c == nil || c.Type() != TypeStr {
// 		return ""
// 	}
// 	return c.String()
// }
//
// // Look and get child bool by given path and separator.
// func (n *Node) GetBoolPS(path, sep string) bool {
// 	vec := n.vec()
// 	if vec == nil {
// 		return false
// 	}
// 	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
// 	c := n.Get(vec.ss...)
// 	if c == nil || c.Type() != TypeBool {
// 		return false
// 	}
// 	return c.Bool()
// }
//
// // Look and get child float by given path and separator.
// func (n *Node) GetFloatPS(path, sep string) (float64, error) {
// 	vec := n.vec()
// 	if vec == nil {
// 		return 0, ErrInternal
// 	}
// 	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
// 	c := n.Get(vec.ss...)
// 	if c == nil {
// 		return 0, ErrNotFound
// 	}
// 	if c.Type() != TypeNum {
// 		return 0, ErrIncompatType
// 	}
// 	return c.Float()
// }
//
// // Look and get child integer by given path and separator.
// func (n *Node) GetIntPS(path, sep string) (int64, error) {
// 	vec := n.vec()
// 	if vec == nil {
// 		return 0, ErrInternal
// 	}
// 	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
// 	c := n.Get(vec.ss...)
// 	if c == nil {
// 		return 0, ErrNotFound
// 	}
// 	if c.Type() != TypeNum {
// 		return 0, ErrIncompatType
// 	}
// 	return c.Int()
// }
//
// // Look and get child unsigned int by given path and separator.
// func (n *Node) GetUintPS(path, sep string) (uint64, error) {
// 	vec := n.vec()
// 	if vec == nil {
// 		return 0, ErrInternal
// 	}
// 	vec.ss = bytealg.AppendSplitStr(vec.ss[:0], path, sep, -1)
// 	c := n.Get(vec.ss...)
// 	if c == nil {
// 		return 0, ErrNotFound
// 	}
// 	if c.Type() != TypeNum {
// 		return 0, ErrIncompatType
// 	}
// 	return c.Uint()
// }
