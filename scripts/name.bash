#!/bin/bash

#this scripts renames every jpg in a folder suppiled on the command line to a
#numeric id


function rename() {
    cd $1
    i=1
    for f in *.jpg; do
        echo "Moving $1$f -> $1$i.jpg"
        mv $f $i.jpg
        let "i++"
    done
}

echo "Renaming $1"
rename $1
