set terminal png size 1500, 1000 font "Sans,22"

set xlabel "B - V, ᵐ"
set xrange [-0.6:1.6]
set xtics -0.6, 0.4, 1.6

set ylabel "Mᵥ"
set yrange [20:-10]

set grid
set key top right box opaque

set fit nolog

set output "magv_to_bv.png"

plot "../data/main_magv_to_bv.dat" with points ls 5 lc rgb 'red' ps 2 title "Звезды ГП", \
    "../data/main_magv_to_bv_interp.dat" with line ls 5 lc rgb 'red' title "Линия звезд ГП", \
    "../data/stars_magv_to_bv.dat" with points ls 5 lc rgb 'blue' ps 2 title "Звезды кластера", \
    "../data/stars_magv_to_bv_corrected.dat" with points ls 5 lc rgb 'green' ps 2 title "Звезды кластера скорректированные"
