package fxml

type (
	//info about the current node passing to the XWalker callback function
	XNodInfo struct {
		//slice of element names, where the first is the root node and the
		//last is the current node
		Path []string
		//zero based index of the current node in its parent's children list
		Index int
		//reverse index of the current node in its parent's children list,
		//where -1 means the last child, -2 is the one before last and so on.
		RIndex int
	}
	//action for the next iteration: continue (WRCont), skip remaining nodes
	//in the same level (WRSkip) or terminate (WRTerm).
	XWalkResult byte
	// Walk the entire XML tree in depth-first order. The first parameter
	// is a XNodInfo struct, the second parameter is the node itself.  It
	// returns XWalkResult indicating whether and how to proceed.
	//
	// Apart from more informative XNodeInfo and better flow control, the
	// most important difference between XWalker and XTraverser is that
	// XWalker works on pointer of XMLTree, i.e. it allows in place modification
	// of the tree nodes.
	XWalker func(XNodInfo, *XMLTree) XWalkResult

	XWalkerWithTree func(XNodInfo, *XMLTree) (XWalkResult, *XMLTree)
)

const (
	WRCont XWalkResult = iota //continue normally
	WRSkip                    //skip remaining nodes in the same level
	WRTerm                    //terminate the walk process
)

func (xt *XMLTree) walk(ni XNodInfo, w XWalker) XWalkResult {
	ni.Path = append(ni.Path, xt.Name.Local)
	wr := w(ni, xt)
	if wr != WRCont {
		return wr
	}
	for i, c := range xt.Children {
		ni.Index = i
		ni.RIndex = i - len(xt.Children)
		wr := c.walk(ni, w)
		xt.Children[i] = c
		switch wr {
		case WRTerm:
			return WRTerm
		case WRSkip:
			return WRCont
		}
	}
	return WRCont
}

func (xt *XMLTree) walkWithReturn(ni XNodInfo, w XWalkerWithTree, x *XMLTree) (XWalkResult, *XMLTree) {
	ni.Path = append(ni.Path, xt.Name.Local)
	wr, x := w(ni, xt)
	if wr != WRCont {
		return wr, x
	}
	for i, c := range xt.Children {
		ni.Index = i
		ni.RIndex = i - len(xt.Children)
		wr, x := c.walkWithReturn(ni, w, x)
		xt.Children[i] = c
		switch wr {
		case WRTerm:
			return WRTerm, x
		case WRSkip:
			return WRCont, x
		}
	}
	return WRCont, x
}

// walk through the XMLTree using the given walker.
func (xt *XMLTree) Walk(w XWalker) {
	xt.walk(XNodInfo{}, w)
}

func (xt *XMLTree) WalkMutate(w XWalkerWithTree) (x *XMLTree) {
	var newXML XMLTree
	_, x = xt.walkWithReturn(XNodInfo{}, w, &newXML)
	return
}
