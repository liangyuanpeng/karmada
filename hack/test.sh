#!/usr/bin/env bash

if cosign > /dev/null 2>&1; then
  echo "cosign exist"
else
  echo "cosign is not exist"
fi