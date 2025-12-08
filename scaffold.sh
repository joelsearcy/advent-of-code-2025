#!/bin/bash

# Scaffold script for Advent of Code 2025 in Go
# Usage: ./scaffold.sh DAY_NUMBER

# zero pad day number to two digits
DAY=$1
DAY_DIR=$(printf "%02d" $DAY)
DAY_FILE="day$DAY.go"
SAMPLE_FILE="sample.txt"
INPUT_FILE="input.txt"

# Create directory for the day
mkdir -p "$DAY_DIR"

# Copy template.go to the new directory as solution.go
cp ./template.go "$DAY_DIR/$DAY_FILE"

cd "$DAY_DIR"

# Initialize a new Go module
go mod init advent-of-code-2025/$DAY_DIR

touch $SAMPLE_FILE $INPUT_FILE

echo "Scaffolded Advent of Code 2025 Day $DAY in directory $DAY_DIR"

# Open the new directory in VS Code and load day$DAY.go, sample.txt, and input.txt
code -r $DAY_FILE $SAMPLE_FILE $INPUT_FILE
