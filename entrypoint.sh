#!/bin/sh

chown -R app:app /static

exec su app -c "/bin/gostatic --files /static --addr :80 $*"
