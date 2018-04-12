#!/bin/sh
set +e

NO_COLOR='\033[0m'
OK_COLOR='\033[32;01m'
ERROR_COLOR='\033[31;01m'
WARN_COLOR='\033[33;01m'
PASS="${OK_COLOR}PASS ${NO_COLOR}"
FAIL="${ERROR_COLOR}FAIL ${NO_COLOR}"

echo "${OK_COLOR}Running features test:${NO_COLOR}"

export DATABASE_DSN="mongodb://localhost:27017/ops_gateway"
export ADMIN_PASSWORD="admin"
export STATS_DSN="noop://"
export PORT="3000"
export API_PORT="3001"
export INSECURE_SKIP_VERIFY="true"
export DEBUG="true"
export LOG_LEVEL="debug"
export SECRET="secret key"
export ADMIN_USERNAME="admin"
export API_READONLY="false"
export PORT_SECONDARY="3100"
export API_PORT_SECONDARY="3101"
export STORAGE_DSN="redis://localhost:6379"
export BACKEND_UPDATE_FREQUENCY="1s"

"./dist/janus" > /tmp/janus.log 2>&1 &
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

"./dist/janus" > /tmp/janus2.log 2>&1 &
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

go test -godog -stop-on-failure
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
    docker-compose logs service1
fi

exit ${exit_code}
