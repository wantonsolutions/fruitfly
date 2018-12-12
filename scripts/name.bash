#!/bin/bash


i=1
for f in *.jpg; do
    mv $f $i.jpg
    let "i++"
done
