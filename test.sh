#!/bin/bash

for i in `seq 0 7`; do

  file="examples/example0${i}.txt"
  echo -e "\n\nrunning file $file"

  go run . "$file"
  read -p "Press enter to continue"

done

for i in `seq 0 1`; do

  file="examples/badexample0${i}.txt"
  echo "\n\nrunning file $file"

  go run . "$file"
  read -p "Press enter to continue"

done
