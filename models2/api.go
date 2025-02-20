package models2

import (
	"strings"
)

func (n *Node) Insert(flatNode FlatNode) {
	parts := strings.Split(flatNode.FullPath, "/") // Split breadcrumb into parts
	current := n

	for _, part := range parts {
		if part == "" {
			continue
		}
		if current.Children == nil {
			current.Children = make(map[string]*Node)
		}
		if _, exists := current.Children[part]; !exists {
			current.Children[part] = &Node{Title: part, ID: flatNode.ID, NodeType: NodeType(flatNode.Model)}
		}
		current = current.Children[part]
	}
}

func (n *Node) GetUserDocuments() []Node {
	var results []Node

	if n.NodeType == DocumentType {
		results = append(results, *n)
	}

	for _, child := range n.Children {
		docs := child.GetUserDocuments()
		results = append(results, docs...)
	}

	return results
}

func InsertDocVersionsAndPages(docs []Node, docPages []DocumentPageRow, mediaRoot string) error {
	return nil
}
