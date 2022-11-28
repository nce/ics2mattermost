#!/usr/bin/env sh
if [ -z $VERSION ]; then
  git rev-parse --short HEAD | tr -d '\n'  > version
else
  echo -n $VERSION > version
fi
