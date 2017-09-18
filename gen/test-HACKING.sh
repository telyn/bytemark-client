#!/bin/bash

EXPECTED=$(grep 'sha256sum:' HACKING)
EXPECTED=$(echo -e "${EXPECTED##*sha256sum: }" | tr -d '[[:space:]]')

ACTUAL=$(tree -dI 'vendor' | sha256sum)
ACTUAL=$(echo -e "${ACTUAL%% -}" | tr -d '[[:space:]]')

[ "$EXPECTED" == "$ACTUAL" ]
