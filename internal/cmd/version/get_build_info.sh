#!/bin/bash

set -e

current_branch=$(git rev-parse --abbrev-ref HEAD)

if [ $current_branch = 'main' ]; then
  # If the current branch is main, write the current commit hash, tag or branch name, and build date to the files
  echo "main branch detected, writing build info to files..."
  git rev-parse HEAD > commit.txt
  version=$(git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD)
  echo "$version" &> version.txt
  git show -z -s --format=%ci > build_date.txt
else
  # If the current branch is not main, fetch the latest values from the main branch and write those to the files
  echo "non-main branch detected, fetching build info from main branch..."
  git checkout main -- ./commit.txt ./version.txt ./build_date.txt
fi