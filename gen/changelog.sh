#!/bin/bash
set -x
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CHANGELOG=$DIR/../doc/changelog
VERSIONGO=$DIR/../lib/version.go

DATE=date
if which gdate > /dev/null; then
    DATE=gdate
fi

mv $CHANGELOG $CHANGELOG.bak
head -n 1 $CHANGELOG.bak > $CHANGELOG
echo "" >> $CHANGELOG
git log --oneline master..develop | sed -e 's/^[a-f0-9]* /  * /' >> $CHANGELOG
echo "" >> $CHANGELOG
echo " -- `git config --get user.name` <`git config --get user.email`>  `$DATE -R`" >> $CHANGELOG
echo "" >> $CHANGELOG
cat $CHANGELOG.bak >> $CHANGELOG
vim $CHANGELOG

VERSION=$(head -n 1 $CHANGELOG | grep -o '(.*)' | grep -oE '[^()]+')

cat > $VERSIONGO <<VERS 
package lib

const (
        // Version is the version number the client believes it is
	Version = "$VERSION"
)
VERS

