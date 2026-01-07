package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/jba/heap"
)

func main() {
	// Generate sizes from 1 to 1 million on a log scale for better visualization
	sizes := []int{
		1, 2, 5, 10, 20, 50, 100, 200, 500,
		1000, 2000, 5000, 10000, 20000, 50000,
		100000, 200000, 500000, 1000000,
	}

	// Create data file
	f, err := os.Create("insert_times.dat")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating data file: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(f, "# heap_size avg_insert_time_ns")

	for _, n := range sizes {
		// Build a heap of size n, then measure single insert time
		const trials = 1000 // Many trials for accurate single-insert timing
		var totalTime time.Duration

		for t := 0; t < trials; t++ {
			// Build heap of size n
			h := heap.NewOrdered[int]()
			for i := 0; i < n; i++ {
				h.Insert(rand.Int())
			}
			// Time a single insert
			val := rand.Int()
			start := time.Now()
			h.Insert(val)
			elapsed := time.Since(start)
			totalTime += elapsed
		}

		avgTime := float64(totalTime.Nanoseconds()) / float64(trials)
		fmt.Printf("heap_size=%d: avg_insert=%.2f ns\n", n, avgTime)
		fmt.Fprintf(f, "%d %.2f\n", n, avgTime)
	}

	if err := f.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "error closing data file: %v\n", err)
		os.Exit(1)
	}

	// Run gnuplot
	cmd := exec.Command("gnuplot", "plot.gp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running gnuplot: %v\n", err)
		fmt.Fprintln(os.Stderr, "Data saved to insert_times.dat - you can plot it manually")
		os.Exit(1)
	}

	fmt.Println("\nPlot saved to insert_times.png")
}
