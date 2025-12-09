package main

import (
	"cmp"
	"slices"
)

type UnionFind struct {
	parent []int
	size   []int
}

func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	size := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &UnionFind{parent: parent, size: size}
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) { // by size
	rootX := uf.Find(x)
	rootY := uf.Find(y)
	if rootX == rootY {
		return
	}
	if uf.size[rootX] < uf.size[rootY] {
		uf.parent[rootX] = rootY
		uf.size[rootY] += uf.size[rootX]
	} else {
		uf.parent[rootY] = rootX
		uf.size[rootX] += uf.size[rootY]
	}
}

func (uf *UnionFind) Connected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

func (uf *UnionFind) TopNSizes(n int) []int {
	groupSizes := make([]int, 0)
	for i := 0; i < len(uf.parent); i++ {
		if i == uf.parent[i] {
			groupSizes = append(groupSizes, uf.size[i])
		}
	}
	slices.SortFunc(groupSizes, func(a, b int) int {
		return cmp.Compare(b, a) // descending order
	})
	if len(groupSizes) > n {
		return groupSizes[:n]
	}
	return groupSizes
}

func (uf *UnionFind) TopNGroups(n int) [][]int {
	groups := make(map[int][]int)
	for i := 0; i < len(uf.parent); i++ {
		root := uf.Find(i)
		groups[root] = append(groups[root], i)
	}
	groupList := make([][]int, 0, len(groups))
	for _, group := range groups {
		groupList = append(groupList, group)
	}
	slices.SortFunc(groupList, func(a, b []int) int {
		return cmp.Compare(len(b), len(a)) // descending order
	})
	if len(groupList) > n {
		return groupList[:n]
	}
	return groupList
}
