#!/bin/ash
echo $(git rev-parse HEAD) > $1
