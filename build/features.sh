#!/bin/sh
set +e

NO_COLOR='\033[0m'
OK_COLOR='\033[32;01m'
ERROR_COLOR='\033[31;01m'
WARN_COLOR='\033[33;01m'
PASS="${OK_COLOR}PASS ${NO_COLOR}"
FAIL="${ERROR_COLOR}FAIL ${NO_COLOR}"

MONGO_PORT="27017"
if [ -n "${DYNAMIC_MONGO_PORT}" ]; then
    MONGO_PORT=${DYNAMIC_MONGO_PORT}
fi

if ! [ -x "$(command -v godog)" ]; then
    echo "${OK_COLOR}Installing missing godog command:${NO_COLOR}"
    # We need to be outside of the project root when installing godog
    # Otherwise we will try to verify dependencies again
    cd ..
    GO111MODULE=on go install github.com/cucumber/godog/cmd/godog@v0.12.6
    cd -
    # add $GOPATH/bin to global path to make command available w/out as path
    export PATH=${PATH}:`go env GOPATH`/bin
fi

echo "${OK_COLOR}Running features test:${NO_COLOR}"

export DATABASE_DSN="mongodb://localhost:${MONGO_PORT}/ops_gateway"
export STATS_DSN="noop://"
export PORT="3000"
export API_PORT="3001"
export DEBUG="true"
export LOG_LEVEL="debug"
export SECRET="secret key"
export BASIC_USERS="admin:admin"
export API_READONLY="false"
export PORT_SECONDARY="3100"
export API_PORT_SECONDARY="3101"
export BACKEND_UPDATE_FREQUENCY="0.5s"

"./dist/janus" start >/tmp/janus.log 2>&1 &
exit_code=$?
if [ ${exit_code} -ne 0 ]; then
    echo "${ERROR_COLOR}Failed to run primary janus instance${NO_COLOR}"
    cat janus.log

    exit ${exit_code}
fi
pid_janus=$!

echo "${OK_COLOR}Started primary instance; PID: ${pid_janus}${NO_COLOR}"

# remember primary instance ports
PORT_PRIMARY=${PORT}
API_PORT_PRIMARY=${API_PORT}

# set ports env variables to secondary to run secondary instance on another ports
export PORT=${PORT_SECONDARY}
export API_PORT=${API_PORT_SECONDARY}

"./dist/janus" start >/tmp/janus2.log 2>&1 &
exit_code=$?
if [ ${exit_code} -ne 0 ]; then
    echo "${ERROR_COLOR}Failed to run secondary janus instance${NO_COLOR}"
    cat /tmp/janus2.log

    exit ${exit_code}
fi
pid_janus2=$!

echo "${OK_COLOR}Started secondary instance; PID: ${pid_janus2}${NO_COLOR}"

# revert port values back
export PORT=${PORT_PRIMARY}
export API_PORT=${API_PORT_PRIMARY}

# make sure app started
sleep 1

godog --format=pretty --random --stop-on-failure --strict
exit_code=$?

kill ${pid_janus}
kill ${pid_janus2}

# Make sure to exit if the test failed
if [ ${exit_code} -ne 0 ]; then
    echo "${WARN_COLOR}=================================${NO_COLOR}"
    echo "${WARN_COLOR}===          PRIMARY          ===${NO_COLOR}"
    echo "${WARN_COLOR}=================================${NO_COLOR}"
    cat /tmp/janus.log

    echo "${WARN_COLOR}=================================${NO_COLOR}"
    echo "${WARN_COLOR}===         SECONDARY         ===${NO_COLOR}"
    echo "${WARN_COLOR}=================================${NO_COLOR}"
    cat /tmp/janus2.log

    echo "${WARN_COLOR}=================================${NO_COLOR}"
    echo "${WARN_COLOR}===         WIREMOCK          ===${NO_COLOR}"
    echo "${WARN_COLOR}=================================${NO_COLOR}"
    docker-compose -f assets/docker-compose.yml logs upstreams
fi

exit ${exit_code}
