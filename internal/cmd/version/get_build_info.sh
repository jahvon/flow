#!/usr/bin/env bash

current_branch=$(git symbolic-ref -q --short HEAD)
current_dir=$(pwd)

if [ "$current_branch" == "main" ]; then
  # If the current branch is main, write the current commit hash, tag or branch name, and build date to the files
  echo -n $(git rev-parse HEAD) > commit.txt
  echo -n $(git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD) &> version.txt
  echo -n $(git show -z -s --format=%ci) > build_date.txt
else
  # If the current branch is not main, fetch the latest values from the main branch and write those to the files
  echo -n $(git show main:$current_dir/commit.txt) > commit.txt
  echo -n $(git show main:$current_dir/version.txt) > version.txt
  echo -n $(git show main:$current_dir/build_date.txt) > build_date.txt
fi