#!/bin/bash

# We don't want to check error return values by our self here.
set -e

tmpfile=$(mktemp /tmp/findy-scan.XXXXXX)

go build -o "$tmpfile" tools/playground/playground.go

# subscript does the scanning and cleanup
./lichen.sh "$tmpfile"

