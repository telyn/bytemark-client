#!/bin/bash

MAJORVERSION=$(grep majorversion lib/version.go | grep -oP '\d+')
MINORVERSION=$(grep minorversion lib/version.go | grep -oP '\d+')
BUILD_NUMBER=$(grep buildnumber lib/version.go  | grep -oP '\d+')
BRANCH=$(grep gitbranch lib/version.go | grep -oP '".*"')
BRANCH=${BRANCH#'"'}
BRANCH=${BRANCH%'"'}

VERSION=$MAJORVERSION.$MINORVERSION.$BUILD_NUMBER~$BRANCH
if $BRANCH = "master"; then
    VERSION=$MAJORVERSION.$MINORVERSION.$BUILD_NUMBER
fi

echo $VERSION
