#!/bin/bash

# Get the name of the current script
current_script=$(basename "$0")

# Run all scripts in the current directory except the current script
for script in *.sh; do
  if [ "$script" != "$current_script" ]; then
    bash "$script"
  fi
done