#!/bin/bash

for i in `seq 0 7`; do

  file="examples/example0${i}.txt"

  echo -e "\nAbout to run file $file\n"
  read -p "Press enter to continue"

  go run . "$file"

done

for i in `seq 0 1`; do

  file="examples/badexample0${i}.txt"

  echo -e "\nAbout to run file $file\n"
  read -p "Press enter to continue"

  go run . "$file"

done
