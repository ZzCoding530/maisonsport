#!/bin/bash

# Check if the program is running on port 8080
if lsof -i :8080 | grep LISTEN; then
  echo "Program is running on port 8080."
else
  echo "Program is not running on port 8080."
fi
