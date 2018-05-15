#!/bin/bash

base_path=$(pwd)
for D in $(find ./** -type d | grep -Ev 'vendor|.git|mocks') ; do
    cd $D
    mockery -all -case=underscore &
    cd $base_path
    echo $!
    pids[${i}]=$!
done

for pid in ${pids[*]}; do
    wait $pid
done
echo DONE
