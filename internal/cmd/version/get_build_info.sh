#!/bin/bash

set -ex

current_branch=$(git symbolic-ref -q --short HEAD)

if [ $current_branch = 'main' ]; then
  # If the current branch is main, write the current commit hash, tag or branch name, and build date to the files
  git rev-parse HEAD > commit.txt
  version=$(git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD || git rev-parse --short HEAD)
  echo "$version" &> version.txt
  git show -z -s --format=%ci > build_date.txt
else
  # If the current branch is not main, fetch the latest values from the main branch and write those to the files
  git show main:./commit.txt > commit.txt
  git show main:./version.txt > version.txt
  git show main:./build_date.txt > build_date.txt
fi