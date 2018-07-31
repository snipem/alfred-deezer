#!/bin/bash
# Download deezer cli for opening files in the browser
wget https://raw.githubusercontent.com/snipem/deezer-cli/master/deezer -O workflow/deezer && 
chmod +x workflow/deezer

# Build alfred deezer api implementation
go build -v -o workflow/alfred-deezer
