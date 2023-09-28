#!/bin/sh

chown -R app:app /static

su app -c "/bin/gostatic --files /static --addr :80 $*"
