#!/bin/bash
echo $(git rev-parse --short HEAD) > $1
