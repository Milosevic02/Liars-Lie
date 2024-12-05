#!/bin/bash

# Path to the executable
EXECUTABLE="./liarslie"

echo "***Test 2***"
echo "Details:"
echo "value 9, max-value 20, num-agents 10, liar-ratio 0.3"
echo "Extend and Play"
cd ..
(
  echo "start --value 9 --max-value 20 --num-agents 10 --liar-ratio 0.3"
  sleep 5
  echo "extend --value 9 --max-value 20 --num-agents 5 --liar-ratio 0.3"
  sleep 5
  echo "play"
  sleep 1
  echo "stop"
) | $EXECUTABLE > output.log

# Check the output
if grep -q "Network value: 9" output.log; then
  echo "Test passed"
else
  echo "Test failed"
fi

rm output.log
