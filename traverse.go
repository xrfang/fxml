package fxml

import "path"

// Walk the entire XML tree in depth-first order. The first parameter
// is a UNIX style path of the current node, the second parameter is
// the "current" node.  If it returns false, traverse will terminate.
type XTraverser func(string, XMLTree) bool

func (xt XMLTree) traverse(pfx string, v XTraverser) bool {
	p := path.Join(pfx, xt.Name.Local)
	if !v(p, xt) {
		return false
	}
	for _, c := range xt.Children {
		if !c.traverse(p, v) {
			return false
		}
	}
	return true
}

// walk through the XMLTree using the given traverser.
func (xt XMLTree) Traverse(v XTraverser) bool {
	return xt.traverse("", v)
}
