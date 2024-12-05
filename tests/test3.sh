#!/bin/bash

# Path to the executable
EXECUTABLE="./liarslie"

echo "***Test 3***"
echo "Details:"
echo "value 9, max-value 20, num-agents 10, liar-ratio 0.3"
echo "Extend, Kill, Play Expert, and Stop"
cd ..
(
  echo "start --value 9 --max-value 20 --num-agents 10 --liar-ratio 0.3"
  sleep 5
  echo "extend --value 9 --max-value 20 --num-agents 5 --liar-ratio 0.3"
  sleep 5
  echo "kill --id 3"
  sleep 1
  echo "playexpert --num-agents 5 --liar-ratio 0.3"
  sleep 1
  echo "stop"
) | $EXECUTABLE > output.log

# Check the output
if grep -q "Network value (expert mode): 9" output.log; then
  echo "Test passed"
else
  echo "Test failed"
fi

rm output.log
