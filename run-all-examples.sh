#!/bin/bash

for filename in ./examples/*.vorn; do
    ./vorn "$filename"
done
