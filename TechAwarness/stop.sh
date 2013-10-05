#!/bin/bash
# ---------------------------------------------

cd Presentation
kill $(pidof godoc) $(pidof present)
