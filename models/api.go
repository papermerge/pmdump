package models

import (
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
)

func MakeUserID2UIDMap(users []User) ID2UUID {
	dict := make(ID2UUID)
	for _, user := range users {
		dict[user.ID] = user.UUID
	}
	return dict
}

func MakeNodeID2UIDMap(nodes []Node) ID2UUID {
	dict := make(ID2UUID)
	for _, node := range nodes {
		dict[node.ID] = node.UUID
	}
	return dict
}

func MakePages(
	doc Document,
	docVer DocumentVersion,
	mediaRoot string,
	docPages []DocumentPageRow,
) ([]Page, error) {
	var pages []Page

	return pages, nil
}

func NewDocument(
	node Node,
	mediaRoot string,
	idsDict IDDict,
	docPages []DocumentPageRow,
) Document {
	var versions []DocumentVersion

	document := Document{
		ID:       node.ID,
		Title:    node.Title,
		UserID:   node.UserID,
		UserUUID: idsDict.UserIDs[node.UserID],
		UUID:     node.UUID,
		ParentID: node.ParentID,
	}

	if node.ParentID != nil {
		document.ParentID = node.ParentID
		parentUUID := idsDict.NodeIDs[*node.ParentID]
		document.ParentUUID = &parentUUID
	}

	originalDocPath := fmt.Sprintf("%s/docs/user_%d/document_%d/%s",
		mediaRoot,
		node.UserID,
		node.ID,
		*node.FileName,
	)
	if _, err := os.Stat(originalDocPath); err == nil {
		version := DocumentVersion{
			Number:   0,
			UUID:     uuid.New(),
			FileName: node.FileName,
		}
		pages, err := MakePages(document, version, mediaRoot, docPages)
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
		node.UserID,
		node.ID,
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
				FileName: node.FileName,
			}
			pages, err := MakePages(document, version, mediaRoot, docPages)
			if err != nil {
				fmt.Printf("Error: NewDocument: %s\n", err)
			} else {
				version.Pages = pages
			}
			versions = append(versions, version)
		}
	}

	document.Versions = versions
	return document
}

func GetDocuments(nodes []Node, mediaRoot string, idsDict IDDict, docPages []DocumentPageRow) ([]Document, error) {
	var documents []Document

	for _, node := range nodes {
		if node.Model == DocumentModelName {
			document := NewDocument(node, mediaRoot, idsDict, docPages)
			documents = append(documents, document)
		}
	}

	return documents, nil
}

func GetFolders(nodes []Node, idsDict IDDict) ([]Folder, error) {
	var folders []Folder

	for _, node := range nodes {
		if node.Model == FolderModelName {

			folder := Folder{
				ID:     node.ID,
				Title:  node.Title,
				UserID: node.UserID,
				UUID:   node.UUID,
			}

			folder.UserUUID = idsDict.UserIDs[folder.UserID]

			if node.ParentID != nil {
				parentUUID := idsDict.NodeIDs[*node.ParentID]
				folder.ParentUUID = &parentUUID
				folder.ParentID = node.ParentID
			}
			folders = append(folders, folder)
		}
	}

	return folders, nil
}

func GetFilePaths(docs []Document, mediaRoot string) ([]FilePath, error) {
	var paths []FilePath

	for _, doc := range docs {
		for _, docVer := range doc.Versions {
			var source string
			if docVer.Number == 0 {
				source = fmt.Sprintf(
					"%s/docs/user_%d/document_%d/%s",
					mediaRoot,
					doc.UserID,
					doc.ID,
					*docVer.FileName,
				)
			} else {
				source = fmt.Sprintf(
					"%s/docs/user_%d/document_%d/v%d/%s",
					mediaRoot,
					doc.UserID,
					doc.ID,
					docVer.Number,
					*docVer.FileName,
				)
			}
			uid := docVer.UUID.String()
			dest := fmt.Sprintf("docvers/%s/%s/%s/%s", uid[0:2], uid[2:4], uid, *docVer.FileName)
			path := FilePath{
				Source: source,
				Dest:   dest,
			}
			paths = append(paths, path)
		}
	}

	return paths, nil
}
