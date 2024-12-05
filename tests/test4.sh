#!/bin/bash

# Path to the executable
EXECUTABLE="./liarslie"

echo "***Test 4***"
echo "Details:"
echo "value 7, max-value 20, num-agents 10, liar-ratio 0.3"
echo "Start, Play Expert, Kill agent 2, Extend, and Play"
cd ..
(
  echo "start --value 7 --max-value 20 --num-agents 10 --liar-ratio 0.3"
  sleep 5
  echo "playexpert --num-agents 5 --liar-ratio 0.3"
  sleep 1
  echo "kill --id 2"
  sleep 1
  echo "extend --value 7 --max-value 20 --num-agents 5 --liar-ratio 0.3"
  sleep 5
  echo "play"
  sleep 1
  echo "stop"
) | $EXECUTABLE > output.log

# Check the output
if grep -q "Network value: 7" output.log; then
  echo "Test passed"
else
  echo "Test failed"
fi

rm output.log
