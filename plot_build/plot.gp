
set terminal png size 1200,800 font "Arial,12"
set output "build_times.png"

set title "Heap Build Time vs Number of Elements" font "Arial,16"
set xlabel "Number of Elements (n)" font "Arial,14"
set ylabel "Time (nanoseconds)" font "Arial,14"

set logscale x 10
set logscale y 10
set grid

set key top left

# Plot with points and lines
plot "build_times.dat" using 1:2 with linespoints pointtype 7 pointsize 1.5 linewidth 2 title "Build time", \
     "" using 1:($1*10) with lines linewidth 1 dashtype 2 title "O(n) reference"
