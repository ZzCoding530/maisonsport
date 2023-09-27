#!/bin/bash

# Run the Go program in the background and save console output to console.log
go run main.go > console.log 2>&1 &

echo "Go program is running in the background. Console output is saved to console.log."

