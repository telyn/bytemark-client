#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CHANGELOG=$DIR/../doc/changelog

mv $CHANGELOG $CHANGELOG.bak
head -n 1 $CHANGELOG.bak > $CHANGELOG
echo "" >> $CHANGELOG
git log --oneline master..develop | sed -e 's/^[a-f0-9]* /  * /' >> $CHANGELOG
echo "" >> $CHANGELOG
echo " -- `git config --get user.name` <`git config --get user.email`> `gdate -R`" >> $CHANGELOG
echo "" >> $CHANGELOG
cat $CHANGELOG.bak >> $CHANGELOG
vim $CHANGELOG
if [ "`head -n 1 $CHANGELOG`" == "`head -n 1 $CHANGELOG.bak`" ]; then
    mv $CHANGELOG.bak $CHANGELOG
else
    rm $CHANGELOG.bak
fi
