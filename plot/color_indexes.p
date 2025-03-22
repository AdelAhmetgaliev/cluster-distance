set terminal png size 1500, 1000 font "Sans,22"

set xlabel "B - V, ᵐ"
set xrange [-0.4:1.8]
set xtics -0.4, 0.4, 1.8

set ylabel "U - B, ᵐ"
set yrange [1.6:-1.2]
set ytics -1.2, 0.2, 1.6

set grid
set key bottom right box opaque

f(x) = -0.601 + 0.72 * x

set output "color_indexes.png"

plot "../data/output/color_indexes_interpolated.dat" with line ls 5 lc rgb 'red' title "Линия нормальных цветов", \
    [:0.3096]f(x) with line ls 5 lc rgb 'web-green' title "Линия покраснения", \
    "../data/output/processing_stars_color_indexes.dat" with points ls 5 lc rgb 'blue' title "Звезды скопления", \
    "../data/output/corrected_stars_color_indexes.dat" with points ls 5 lc rgb 'navy' title "Звезды скопления скоректированные", \
    "../data/output/average_color_index.dat" with points ls 7 lc rgb 'dark-red' ps 3 title "Средняя звезда скопления", \
    "../data/output/corrected_color_index.dat" with points ls 7 lc rgb 'dark-cyan' ps 3 title "Скорректированная средняя звезда скопления"
