#!/bin/sh

set -ex

cd "$(dirname "$0")"

export PATH=$(pwd):$PATH

./test_57/test.sh
./test_8/test.sh
