#!/bin/bash
set -x

EXPECTED=$(grep 'sha256sum:' HACKING)
EXPECTED=$(echo -e "${EXPECTED##*sha256sum: }" | tr -d '[[:space:]]')

TREE=$(tree -dI 'vendor')
echo -e "$TREE"

ACTUAL=$(echo -e "$TREE" | sha256sum)
ACTUAL=$(echo -e "${ACTUAL%% -}" | tr -d '[[:space:]]')

[ "$EXPECTED" == "$ACTUAL" ]
