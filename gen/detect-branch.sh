#!/bin/bash
ORIGIN_RELEASE_REFS=$(git for-each-ref --format "%(refname)" refs/remotes/origin/release\*)
ORIGIN_REFS=$(git for-each-ref --format "%(refname)" refs/remotes/origin)
LOCAL_RELEASE_REFS=$(git for-each-ref --format "%(refname)" refs/heads/release\*)
LOCAL_REFS=$(git for-each-ref --format "%(refname)" refs/heads)

for i in refs/heads/master $LOCAL_RELEASE_REFS refs/heads/develop $LOCAL_REFS \
    refs/remotes/origin/master $ORIGIN_RELEASE_REFS refs/remotes/origin/develop $ORIGIN_REFS
    do
	[ "`git rev-parse $i`" == "`git rev-parse HEAD`" ] && git rev-parse --abbrev-ref $i | awk -F'/' '{ if(length($2) != 0) { print $2 } else { print $1 } }' && exit 0
done
echo HEAD
