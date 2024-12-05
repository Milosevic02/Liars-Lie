#!/bin/bash

# Path to the executable
EXECUTABLE="./liarslie"

echo "***Test 1***"
echo "Details:"
echo "value 7, max-value 20, num-agents 10, liar-ratio 0.7"
echo "Killing agents: 2, 5"
echo "Play"
cd ..
(
  echo "start --value 7 --max-value 20 --num-agents 10 --liar-ratio 0.7"
  sleep 5
  echo "kill --id 2"
  sleep 1
  echo "kill --id 5"
  sleep 1
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