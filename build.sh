#!/bin/bash

set -e

mkdir -p bin

for dir in cmd/*; do
    if [ -d "$dir" ]; then
        tool=$(basename "$dir")
        echo "Building $tool..."
        go build -o "bin/${tool}.exe" "./cmd/${tool}"
    fi
done

echo "Done!"
