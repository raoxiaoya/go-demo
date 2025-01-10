package iterator

import (
	"fmt"
	"iter"
	"strings"
)

// https://www.liwenzhou.com/posts/Go/iter/

func Run() {
	ShowMembers2()
}

type Member struct {
	members []string
}

func NewMember() *Member {
	m := &Member{}
	m.members = make([]string, 0)
	return m
}

func (m *Member) Add(name string) {
	m.members = append(m.members, name)
}

func (m *Member) Iterator() func(yield func(string) bool) {
	return func(yield func(string) bool) {
		for _, v := range m.members {
			if !yield(strings.ToUpper(v)) {
				break
			}
		}
	}
}

func ShowMembers() {
	m := NewMember()
	m.Add("zhangsan")
	m.Add("lisi")
	m.Add("wangwu")

	for v := range m.Iterator() {
		if v == "LISI" {
			break
		}
		fmt.Println("name:", v)
	}
}

func ShowMembers1() {
	m := NewMember()
	m.Add("zhangsan")
	m.Add("lisi")
	m.Add("wangwu")

	next, stop := iter.Pull(m.Iterator())
	defer stop()
	for {
		v, ok := next()
		if !ok {
			break
		}
		if v == "LISI" {
			stop()
			break
		}
		fmt.Println("name:", v)
	}
}

func ShowMembers2() {
	m := NewMember()
	m.Add("zhangsan")
	m.Add("lisi")
	m.Add("wangwu")

	m.Iterator()(func(v string) bool {
		fmt.Println("name:", v)
		return true
	})
}
