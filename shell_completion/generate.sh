#!/usr/bin/env bash

./build/dbconntest
exit 0

bin='./build/dbconntest'
bin_alt='../build/dbconntest'

if [ ! -x $bin ]; then
  if [ -x $bin_alt ] ; then
    bin=$bin_alt
  else
    echo "Please run from project root or from thiencs folder"
    exit 1
  fi
fi

if [ -f shell_completion/_dbconntest ]; then
  rm -f shell_completion/_dbconntest
fi

if [ -f shell_completion/dbconntest ]; then
    rm -f shell_completion/dbconntest
fi

$bin completion bash > shell_completion/dbconntest
$bin completion zsh  > shell_completion/_dbconntest
