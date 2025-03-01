package models_app_v3_3

import (
	"fmt"
	"os"
	"strconv"
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
	n *Node,
	user_id int,
	docPages []DocumentPageRow,
	mediaRoot string,
	op NodeOperation,
) {
	if n.NodeType == DocumentType {
		op(n, user_id, docPages, mediaRoot)
	}

	for _, child := range n.Children {
		ForEachDocument(child, user_id, docPages, mediaRoot, op)
	}
}

func InsertDocVersionsAndPages(
	n *Node,
	user_id int,
	docPages []DocumentPageRow,
	mediaRoot string,
) {
	var versions []DocumentVersion

	originalDocPath := fmt.Sprintf("%s/docs/user_%d/document_%d/%s",
		mediaRoot,
		user_id,
		n.ID,
		*n.FileName,
	)

	if _, err := os.Stat(originalDocPath); err == nil {
		version := DocumentVersion{
			Number:   0,
			UUID:     uuid.New(),
			FileName: *n.FileName,
		}
		pages, err := MakePages(n, user_id, version, mediaRoot, docPages)
		if err != nil {
			fmt.Printf("Error: NewDocument: %s\n", err)
		} else {
			version.Pages = pages
		}
		versions = append(versions, version)
	}

	path := fmt.Sprintf(
		"%s/docs/user_%d/document_%d/",
		mediaRoot,
		user_id,
		n.ID,
	)
	entries, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("Error reading directory:", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			strVersionNumber := entry.Name()[1:]
			versionNumber, err := strconv.Atoi(strVersionNumber)
			if err != nil {
				fmt.Printf("Error: %v", err)
				continue
			}
			version := DocumentVersion{
				Number:   versionNumber,
				UUID:     uuid.New(),
				FileName: *n.FileName,
			}
			pages, err := MakePages(n, user_id, version, mediaRoot, docPages)
			if err != nil {
				fmt.Printf("Error: NewDocument: %s\n", err)
			} else {
				version.Pages = pages
			}
			versions = append(versions, version)
		}
	}

	n.Versions = versions
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
	n.NodeUUID = uuid.New()
}

func MakePages(
	n *Node,
	user_id int,
	docVer DocumentVersion,
	mediaRoot string,
	docPages []DocumentPageRow,
) ([]Page, error) {
	var pages []Page

	for _, entry := range docPages {
		pages = append(pages, Page{
			Number: entry.PageNumber,
			Text:   entry.Text,
			UUID:   uuid.New(),
		})
	}

	var pagesPath string
	// found out pages from filesystem
	if docVer.Number == 0 {
		pagesPath = fmt.Sprintf("%s/results/user_%d/document_%d/pages/",
			mediaRoot,
			user_id,
			n.ID,
		)
	} else {
		pagesPath = fmt.Sprintf("%s/results/user_%d/document_%d/v%d/pages/",
			mediaRoot,
			user_id,
			n.ID,
			docVer.Number,
		)
	}
	pageFiles, err := os.ReadDir(pagesPath)

	if err != nil {
		fmt.Println("Error reading directory:", err)
	}

	for _, pageFile := range pageFiles {
		if !pageFile.IsDir() {
			// cut '.txt' part
			fullName := pageFile.Name()
			name := fullName[:len(fullName)-4]
			// cut 'page_' part
			name = name[5:]
			pageNumber, err := strconv.Atoi(name)
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			fullFilePath := pagesPath + pageFile.Name()
			data, err := os.ReadFile(fullFilePath)

			if err != nil {
				fmt.Printf("Error: %v", err)
			}

			pages = append(pages, Page{
				Number: pageNumber,
				UUID:   uuid.New(),
				Text:   string(data),
			})
		}
	}

	return pages, nil
}

func GetFilePaths(docs []Node, user_id int, mediaRoot string) ([]types.FilePath, error) {
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
			uid := docVer.UUID.String()
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
