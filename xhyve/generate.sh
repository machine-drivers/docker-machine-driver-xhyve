#!/bin/bash

git apply upstream.patch  || echo "Was the patch already applied?"
for f in $(egrep '^\s*src/' upstream/Makefile  | tr '\' ' '); do ln -sf "upstream/$f" "$(basename $f)"; done
