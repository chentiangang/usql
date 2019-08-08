package main

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path"
	"strings"

	"github.com/manifoldco/promptui"

	"gopkg.in/yaml.v2"
)

type Node struct {
	Name     string  `json:"name"`
	Url      string  `json:"url"`
	Children []*Node `json:"children"`
}

func LoadConfig() (node []*Node, err error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(path.Join(u.HomeDir, ".usql.yml"))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, &node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

const prev = "-parent-"

var (
	templates = &promptui.SelectTemplates{
		Help:     "Use: \x1b[91m↓(j) ↑(k) → ← \x1b[0m Search: \x1b[91m/\x1b[0m ",
		Label:    "✨ {{ . | green}}",
		Active:   "➤  {{ .Name | cyan }} {{if .Url}}{{.Url | cyan}}{{end}}",
		Inactive: "  {{.Name | faint}} ",
	}
)

func choose(parent, trees []*Node) *Node {
	prompt := promptui.Select{
		Label:     "Select Database:",
		Items:     trees,
		Templates: templates,
		Searcher: func(input string, index int) bool {
			node := trees[index]
			content := fmt.Sprintf("%s %s", node.Name, node.Url)

			// 如果输入的值中有空格
			if strings.Contains(input, " ") {
				// 则以空格为分隔符,
				for _, key := range strings.Split(input, " ") {

					// 删除里面的空格
					key = strings.TrimSpace(key)

					// 如果key不等于空
					if key != "" {
						// 如果key不存在 name , url中, 返回false
						if !strings.Contains(content, key) {
							return false
						}
					}
				}
				// 否则是存在的，返回true
				return true
			}
			// 如果输入的字符是name和name中的子串，返回true
			if strings.Contains(content, input) {
				return true
			}

			//否则返回false
			return false
		},
	}
	//
	index, _, err := prompt.Run()
	if err != nil {
		return nil
	}
	node := trees[index]
	if len(node.Children) > 0 {
		first := node.Children[0]
		if first.Name != prev {
			first = &Node{Name: prev}
			node.Children = append(node.Children[:0], append([]*Node{first}, node.Children...)...)
		}
		return choose(trees, node.Children)
	}
	if node.Name == prev {
		return choose(nil, parent)
	}
	return node
}
