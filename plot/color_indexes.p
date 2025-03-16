set terminal png size 1500, 1000 font "Sans,22"

set xlabel "B - V, ᵐ"
set xrange [-0.4:2.4]
set xtics -0.4, 0.4, 2.4

set ylabel "U - B, ᵐ"
set yrange [1.6:-1.0]
set ytics -1.0, 0.2, 1.6

set grid
set key top right box opaque

set fit nolog

set output "color_indexes.png"

plot "../data/main_color_indexes.dat" with points ls 5 lc rgb 'red' ps 2 title "Звезды ГП", \
    "../data/stars_color_indexes.dat" with points ls 5 lc rgb 'blue' ps 2 title "Звезды кластера"
