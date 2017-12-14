#!/bin/bash

cat << EOF
*******************************************************************************
*** ERROR: gcr.io/cloud-builders/golang-project:wheezy is deprecated, use   ***
*** gcr.io/cloud-builders/go builder instead.                               ***
*******************************************************************************
EOF

exit 1
