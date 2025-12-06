# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a design repository for a heap data structure API intended for the Go standard library. The repository is in early stages and currently contains only design documentation with no implementation code yet.

## Key Design Documents

- **THOUGHTS.txt**: Contains detailed design considerations for the heap API, including:
  - Proposed types: `Heap[T cmp.Ordered]` and `HeapFunc[T any]`
  - Method signatures and their rationale
  - The `Item` type for tracking heap elements to enable deletion and value changes
  - Trade-offs around delayed heapification and Item allocation overhead
  - Comparison with `container/heap` from the standard library

## Development Commands

This is a Go project. Standard Go commands will be used once implementation begins:

- **Run tests**: `go test ./...`
- **Run tests with coverage**: `go test -cover ./...`
- **Build**: `go build ./...`
- **Run specific test**: `go test -run TestName`
- **Format code**: `go fmt ./...`

## Architecture Considerations

The core design challenge is providing access to heap elements for modification and deletion while keeping internal storage private. The proposed solution uses an `Item` type returned from `Insert()` that acts as a correlate to track elements as they move within the heap.

**Delayed heapification**: Initial inserts only append to a slice; the heap is actually built on the first call to `Min()` or `ExtractMin()`, or when the user explicitly calls `Build()`.

When implementing, pay careful attention to the trade-offs discussed in THOUGHTS.txt, particularly around Item allocation overhead and the delayed heapification behavior.
