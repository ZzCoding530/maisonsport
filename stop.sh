#!/bin/bash

# Find the process ID (PID) of the program running on port 8080
pid=$(lsof -t -i:8080)

# Check if the PID is found
if [ -n "$pid" ]; then
  echo "Stopping program with PID $pid..."
  kill $pid
else
  echo "No program is running on port 8080."
fi
