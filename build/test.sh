#!/bin/sh
#
# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

export CGO_ENABLED=1
NO_COLOR='\033[0m'
OK_COLOR='\033[32;01m'
ERROR_COLOR='\033[31;01m'
WARN_COLOR='\033[33;01m'
PASS="${OK_COLOR}PASS ${NO_COLOR}"
FAIL="${ERROR_COLOR}FAIL ${NO_COLOR}"

TARGETS=$@

if [ "$RUN_INTEGRATION" ]; then
    echo "${OK_COLOR}Running unit and integration tests: ${NO_COLOR}"
    go test -race -cover -cpu=1,2,4 -tags=integration ${TARGETS}
else
    echo "${OK_COLOR}Running unit tests: ${NO_COLOR}"
    go test -race -cover -cpu=1,2,4 ${TARGETS}
fi

echo "${OK_COLOR}Formatting: ${NO_COLOR}"
ERRS=$(find cmd pkg -type f -name \*.go | xargs gofmt -l 2>&1 || true)
if [ -n "${ERRS}" ]; then
    echo "${ERROR_COLOR}FAIL - the following files need to be gofmt'ed: ${NO_COLOR}"
    for e in ${ERRS}; do
        echo "    $e"
    done
    exit 1
fi
echo ${PASS}

echo "${OK_COLOR}Vetting: ${NO_COLOR}"
ERRS=$(go vet ${TARGETS} 2>&1 || true)
if [ -n "${ERRS}" ]; then
    echo ${FAIL}
    echo "${ERRS}"
    exit 1
fi
echo ${PASS}

echo "${OK_COLOR}Linting: ${NO_COLOR}"
ERRS=$(golint ${TARGETS} 2>&1 || true)
if [ -n "${ERRS}" ]; then
    echo ${FAIL}
    echo "${ERRS}"
    exit 1
fi
echo ${PASS}
