rm all_apples.png
python all_apple.py
montage -density 100 -tile 5x5 -geometry +5+50 -border 10 *.png all_apples.png
mv all_apples.png tmp
rm *.png
mv tmp all_apples.png
