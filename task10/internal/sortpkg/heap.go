package sortpkg

import "github.com/pozedorum/WB_project_2/task10/pkg/options"

type HeapItem struct {
	line  string // Текущая строка
	key   string
	index int // Индекс сканера (0..len(scanners)-1)
}

type MinHeap struct {
	items []HeapItem
	fs    options.FlagStruct
}

func (h MinHeap) Len() int {
	return len(h.items)
}

func (h MinHeap) Less(i, j int) bool {
	if *h.fs.RFlag {
		return h.items[i].key > h.items[j].key
	}

	return h.items[i].key < h.items[j].key
}

func (h MinHeap) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *MinHeap) Push(x any) {
	h.items = append(h.items, x.(HeapItem))
}

func (h *MinHeap) Pop() any {
	old := h.items
	n := len(old)
	x := old[n-1]
	h.items = old[:n-1]
	return x
}
