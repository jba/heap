// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Grafana/Pyroscope minheap implementation adapted for benchmarking.
// Original: https://github.com/grafana/pyroscope/blob/3c106704479ba88f2a43280c62d9840af140fa22/pkg/util/minheap/minheap.go

package heap

// grafanaPush adds an element to the heap.
func grafanaPush(h []int, x int) []int {
	h = append(h, x)
	grafanaUp(h, len(h)-1)
	return h
}

// grafanaPop removes the minimum element from the heap.
func grafanaPop(h []int) []int {
	n := len(h) - 1
	h[0], h[n] = h[n], h[0]
	grafanaDown(h, 0, n)
	n = len(h)
	h = h[0 : n-1]
	return h
}

// grafanaChangeMin replaces the minimum and restores heap invariant.
func grafanaChangeMin(h []int, x int) {
	h[0] = x
	grafanaDown(h, 0, len(h))
}

// grafanaMin returns the minimum element.
func grafanaMin(h []int) int {
	return h[0]
}

func grafanaUp(h []int, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || (h[j] >= h[i]) {
			break
		}
		h[i], h[j] = h[j], h[i]
		j = i
	}
}

func grafanaDown(h []int, i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h[j2] < h[j1] {
			j = j2 // = 2*i + 2  // right child
		}
		if h[j] >= h[i] {
			break
		}
		h[i], h[j] = h[j], h[i]
		i = j
	}
	return i > i0
}

// grafanaHeapify builds a heap from an unsorted slice.
func grafanaHeapify(h []int) {
	n := len(h)
	for i := n/2 - 1; i >= 0; i-- {
		grafanaDown(h, i, n)
	}
}
