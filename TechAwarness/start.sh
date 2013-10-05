#!/bin/bash
# ---------------------------------------------

cd Presentation
$HOME/go/bin/present -http=localhost:10001 &
unset GOPATH
godoc -http=localhost:10000 &
