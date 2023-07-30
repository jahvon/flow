#!/usr/bin/env bash

echo -n $(git rev-parse HEAD) > commit.txt

echo -n $(git describe --tags --exact-match 2>/dev/null \
  || git symbolic-ref -q --short HEAD \
	|| git rev-parse --short HEAD) &> version.txt

echo -n $(git show -z -s --format=%ci) > build_date.txt
