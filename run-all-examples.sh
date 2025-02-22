#!/bin/bash

# Quit if any example fails
set -e

for filename in ./examples/*.vorn; do
    ./vorn "$filename"
done
