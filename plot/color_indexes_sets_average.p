set terminal png size 1500, 1000 font "Sans,22"

set xlabel "B - V, ᵐ"
set xrange [-0.4:2.8]
set xtics -0.4, 0.4, 2.8

set ylabel "U - B, ᵐ"
set yrange [1.6:-1.0]
set ytics -1.0, 0.2, 1.6

set grid
set key top right box opaque

set fit nolog

f(x) = -1.367204 + 0.72 * x
g(x) = 0.242136 + 0.72 * x
h(x) = -0.6303 + 0.72 * x

set output "color_indexes_sets_average.png"

plot "../data/main_color_indexes.dat" with points ls 5 lc rgb 'red' ps 2 title "Звезды ГП", \
    "../data/main_color_indexes_interp.dat" with line ls 5 lc rgb 'red' title "Линия нормальных цветов", \
    "../data/stars1_average_color_index.dat" with points ls 5 lc rgb 'green' ps 2 title "Ср. звезда кластера №1", \
    "../data/stars2_average_color_index.dat" with points ls 5 lc rgb 'purple' ps 2 title "Ср. звезда кластера №2", \
    "../data/stars3_average_color_index.dat" with points ls 5 lc rgb 'navy' ps 2 title "Ср. звезда кластера №3", \
    [-0.4:1.1532]f(x) with line ls 5 lc rgb 'green' title "Линия покраснения №1", \
    [-0.4:0.3912]g(x) with line ls 5 lc rgb 'purple' title "Линия покраснения №2", \
    [-0.4:0.3475]h(x) with line ls 5 lc rgb 'navy' title "Линия покраснения №3"
