package models_app_v3_3

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/papermerge/pmdump/types"
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
			current.Children[part] = &Node{
				Title:     part,
				ID:        flatNode.ID,
				NodeType:  NodeType(flatNode.Model),
				FileName:  flatNode.FileName,
				PageCount: flatNode.PageCount,
				Version:   flatNode.Version,
			}
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

func ForEachDocument(
	db *types.DBConn,
	n any,
	op NodeOperation,
) {
	node := n.(*Node)

	if node.NodeType == DocumentType {
		op(db, n)
	}

	for _, child := range node.Children {
		ForEachDocument(db, child, op)
	}
}

func ForEachNode(
	n *Node,
	quickOper NodeQuickOperation,
) {

	quickOper(n)

	for _, child := range n.Children {
		ForEachNode(child, quickOper)
	}
}

func UpdateNodeUUID(n *Node) {
	n.ID = uuid.New()
}

func GetFilePaths(docs []Node, user_id uuid.UUID, mediaRoot string) ([]types.FilePath, error) {
	var paths []types.FilePath

	for _, doc := range docs {
		for _, docVer := range doc.Versions {
			var source string
			if docVer.Number == 0 {
				source = fmt.Sprintf(
					"%s/docs/user_%d/document_%d/%s",
					mediaRoot,
					user_id,
					doc.ID,
					*doc.FileName,
				)
			} else {
				source = fmt.Sprintf(
					"%s/docs/user_%d/document_%d/v%d/%s",
					mediaRoot,
					user_id,
					doc.ID,
					docVer.Number,
					*doc.FileName,
				)
			}
			uid := docVer.ID.String()
			dest := fmt.Sprintf("docvers/%s/%s/%s/%s", uid[0:2], uid[2:4], uid, *doc.FileName)
			path := types.FilePath{
				Source: source,
				Dest:   dest,
			}
			paths = append(paths, path)
		}
	}

	return paths, nil
}

func (users TargetUserList) Get(username string) *TargetUser {
	for i := 0; i < len(users); i++ {
		if users[i].Username == username {
			return &users[i]
		}
	}

	return nil
}
