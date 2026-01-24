set terminal png size 1200,800 font "Arial,12"
set output "insert_times.png"

set title "Single Insert Time vs Heap Size" font "Arial,16"
set xlabel "Heap Size (n)" font "Arial,14"
set ylabel "Insert Time (nanoseconds)" font "Arial,14"

set logscale x 10
set grid

set key top left

# Plot with points and lines
# Expected complexity: O(log n) for single insert
plot "insert_times.dat" using 1:2 with linespoints pointtype 7 pointsize 1.5 linewidth 2 title "Insert time", \
     "" using 1:(log($1)/log(2)*10) with lines linewidth 1 dashtype 2 title "O(log n) reference"
