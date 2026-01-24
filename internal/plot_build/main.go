package main

import (
	"cmp"
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
	f, err := os.Create("build_times.dat")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating data file: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(f, "# n time_ns time_per_element_ns")

	for _, n := range sizes {
		// Run multiple trials and take the average
		const trials = 5
		var totalTime time.Duration

		for t := 0; t < trials; t++ {
			// Create random data
			data := make([]int, n)
			for i := range data {
				data[i] = rand.Int()
			}

			// Time the InsertSlice operation (which includes heapify)
			h := heap.New(cmp.Compare[int])
			start := time.Now()
			h.InsertSlice(data)
			elapsed := time.Since(start)
			totalTime += elapsed
		}

		avgTime := totalTime / trials
		timePerElement := float64(avgTime.Nanoseconds()) / float64(n)
		fmt.Printf("n=%d: avg=%v (%.2f ns/element)\n", n, avgTime, timePerElement)
		fmt.Fprintf(f, "%d %d %.2f\n", n, avgTime.Nanoseconds(), timePerElement)
	}

	if err := f.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "error closing data file: %v\n", err)
		os.Exit(1)
	}

	// Create gnuplot script
	gnuplotScript := `
set terminal png size 1200,800 font "Arial,12"
set output "build_times.png"

set title "Heap InsertSlice Time vs Number of Elements" font "Arial,16"
set xlabel "Number of Elements (n)" font "Arial,14"
set ylabel "Time (nanoseconds)" font "Arial,14"

set logscale x 10
set logscale y 10
set grid

set key top left

# Plot with points and lines
plot "build_times.dat" using 1:2 with linespoints pointtype 7 pointsize 1.5 linewidth 2 title "InsertSlice time", \
     "" using 1:($1*10) with lines linewidth 1 dashtype 2 title "O(n) reference"
`
	if err := os.WriteFile("plot.gp", []byte(gnuplotScript), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing gnuplot script: %v\n", err)
		os.Exit(1)
	}

	// Run gnuplot
	cmd := exec.Command("gnuplot", "plot.gp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running gnuplot: %v\n", err)
		fmt.Fprintln(os.Stderr, "Data saved to build_times.dat - you can plot it manually")
		os.Exit(1)
	}

	fmt.Println("\nPlot saved to build_times.png")
}
