#!/bin/bash


function toucher() {
    dir=$1
    for f in $dir/*; do
        echo "touching " $f
        touch $f
    done
}

path=../data/Aging_Study_1

toucher $path/nov_28_2018
toucher $path/dec_3_2018
toucher $path/dec_4_2018
toucher $path/dec_5_2018
toucher $path/dec_6_2018
toucher $path/dec_7_2018
toucher $path/dec_10_2018
