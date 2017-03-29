/*
Copyright 2017 yasukun

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package treebeard

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const jspre = "treebeard"

type (
	Node struct {
		Id       string   `json:"_"`
		ParentId string   `json:"_"`
		Name     string   `json:"name"`
		Toggled  bool     `json:"toggled,omitempty"`
		Active   bool     `json:"active"`
		Path     []string `json:"path,omitempty"`
		Isdir    bool     `json:"isdir,omitempty"`
		Children []*Node  `json:"children,omitempty"`
	}
)

func (this *Node) Size() int {
	var size int = len(this.Children)
	for _, c := range this.Children {
		size += c.Size()
	}
	return size
}

func (this *Node) Add(nodes ...*Node) bool {
	var size = this.Size()
	for _, n := range nodes {
		if n.ParentId == this.Id {
			this.Children = append(this.Children, n)
		} else {
			for _, c := range this.Children {
				if c.Add(n) {
					break
				}
			}
		}
	}
	return this.Size() == size+len(nodes)
}

// pathList ...
func pathList(path string) []string {
	return strings.Split(path, string(filepath.Separator))
}

// DirWalk ...
func DirWalk(root string) (*Node, []*Node, error) {
	parents := make(map[string]int)
	var rootNode *Node = nil
	jstree := []*Node{}
	idx := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		rel, _ := filepath.Rel(root, path)
		if _, ok := parents[rel]; !ok {
			parents[rel] = idx
		}
		if val, ok := parents[filepath.Dir(rel)]; ok {
			node_id := fmt.Sprintf("%s_%d", jspre, idx)
			parent_id := fmt.Sprintf("%s_%d", jspre, val)
			if node_id == parent_id {
				parent_id = "#"
			}
			// log.Printf("id: %s, parent: %s, text: %s, path: %s\n", node_id, parent_id, info.Name(), path)
			if idx == 0 {
				rootNode = &Node{Id: node_id, ParentId: parent_id, Name: info.Name(), Toggled: true, Active: false, Children: nil}
			} else {
				jstree = append(jstree, &Node{Id: node_id, ParentId: parent_id, Name: info.Name(), Path: pathList(path), Isdir: info.IsDir(), Active: false, Children: nil})
			}
		}
		idx += 1
		return nil
	})

	if err != nil {
		return rootNode, jstree, err
	}

	return rootNode, jstree, nil
}

// TreeBeard ...
func TreeBeard(rootDir string) (*Node, error) {
	root, data, err := DirWalk(rootDir)
	if err != nil {
		return nil, err
	}
	root.Add(data...)
	return root, nil
}
