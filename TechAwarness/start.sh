#!/bin/bash
# ---------------------------------------------

cd Presentation
$HOME/go/bin/present -http=localhost:10001 &
godoc -http=localhost:10003 &
unset GOPATH
godoc -http=localhost:10000 &
