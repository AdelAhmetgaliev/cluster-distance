set terminal png size 1500, 1000 font "Sans,22"

set xlabel "B - V, ᵐ"
set xrange [-0.4:1.8]
set xtics -0.4, 0.4, 1.8

set ylabel "Mv, ᵐ"
set yrange [12.0:-5.6]
set ytics -5.6, 2.0, 12.0

set grid
set key top right box opaque

f(x) = -0.601 + 0.72 * x

set output "magv_to_bv.png"

plot "../data/output/magv_to_bv_interpolated.dat" with line ls 5 lc rgb 'red' title "Звезды ГП", \
    "../data/output/processing_stars_magv_to_bv.dat" with points ls 5 lc rgb 'blue' title "Звезды скопления", \
    "../data/output/corrected_stars_magv_to_bv.dat" with points ls 5 lc rgb 'navy' title "Скорректированные звезды скопления"
