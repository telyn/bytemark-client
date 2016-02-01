#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [ -z "$MAJORVERSION" ]; then
    MAJORVERSION=0
fi
if [ -z "$MINORVERSION" ]; then
    MINORVERSION=0
fi
BUILD_DATE=`date +%Y-%m-%d\ %H:%M`
if [ -z "$BUILD_NUMBER" ]; then
    BUILD_NUMBER=0
fi
GIT_COMMIT=`git rev-parse HEAD`
GIT_BRANCH=`$DIR/detect-branch.sh`

echo "package lib" > version.go
echo "const (" >> version.go
echo "  majorversion = $MAJORVERSION" >> version.go
echo "  minorversion = $MINORVERSION" >> version.go
echo "  buildnumber = $BUILD_NUMBER" >> version.go
echo "  gitcommit = \"$GIT_COMMIT\"" >> version.go
echo "  gitbranch = \"$GIT_BRANCH\"" >> version.go
echo "  builddate = \"$BUILD_DATE\"" >> version.go
echo ")" >> version.go
