#!/usr/bin/env bash

env_dir=$1

if [[ -z "$env_dir" ]]; then
  echo "usage: $0 <environment-directory>"
  exit 1
fi

if [ ! -d "$env_dir" ]; then
  echo "unknown environment"
  exit 1
fi

cd ./bootstrap && ./build.sh
cd ../$env_dir && ./publish.sh
