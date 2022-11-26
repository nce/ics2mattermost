#!/usr/bin/env sh
git rev-parse --short HEAD | tr -d '\n'  > version
