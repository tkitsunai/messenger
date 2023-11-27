package core

import (
	"sort"
)

type SubjectNode struct {
	branches map[string]*SubjectNode
	leaves   map[string]struct{}
}

func newSubjectNode() *SubjectNode {
	return &SubjectNode{
		branches: map[string]*SubjectNode{},
		leaves:   map[string]struct{}{},
	}
}

func (n *SubjectNode) Add(subjects []string) {
	if len(subjects) == 0 {
		return
	}
	sort.Strings(subjects)
	n.add(subjects)
}

func (n *SubjectNode) Remove(subjects []string) {
	if len(subjects) == 0 {
		return
	}
	sort.Strings(subjects)
	n.remove(subjects)
}

func (n *SubjectNode) Match(subjects []string) bool {
	sort.Strings(subjects)
	return n.match(subjects)
}

func (n *SubjectNode) add(subjects []string) {
	if len(subjects) == 1 {
		n.leaves[subjects[0]] = struct{}{}
		return
	}
	branch, ok := n.branches[subjects[0]]
	if !ok {
		branch = newSubjectNode()
		n.branches[subjects[0]] = branch
	}
	branch.add(subjects[1:])
}

func (n *SubjectNode) remove(subjects []string) {
	if len(subjects) == 1 {
		delete(n.leaves, subjects[0])
		return
	}

	if branch, ok := n.branches[subjects[0]]; ok {
		branch.remove(subjects[1:])
	}
}

func (n *SubjectNode) match(subjects []string) bool {
	if len(subjects) == 0 {
		return false
	}

	for _, key := range subjects {
		if _, ok := n.leaves[key]; ok {
			return true
		}
	}
	branch, ok := n.branches[subjects[0]]
	if !ok {
		return n.match(subjects[1:])
	}
	return branch.match(subjects[1:])
}
