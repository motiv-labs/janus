#!/usr/bin/env bash

NO_COLOR='\033[0m'
OK_COLOR='\033[32;01m'
ERROR_COLOR='\033[31;01m'
WARN_COLOR='\033[33;01m'

AUTH_MOCK_PORT="9088"
if [[ -n "${DYNAMIC_AUTH_PORT}" ]]; then
    AUTH_MOCK_PORT=${DYNAMIC_AUTH_PORT}
fi

UPSTREAMS_MOCK_PORT="9089"
if [[ -n "${DYNAMIC_UPSTREAMS_PORT}" ]]; then
    UPSTREAMS_MOCK_PORT=${DYNAMIC_UPSTREAMS_PORT}
fi

echo "${OK_COLOR}Uploading auth mock fixture:${NO_COLOR}"
curl -X DELETE --silent --show-error --fail --output /dev/null "http://localhost:${AUTH_MOCK_PORT}/__admin/mappings"
for fixture in assets/stubs/auth-service/*; do
    curl --silent --show-error --fail --output /dev/null "http://localhost:${AUTH_MOCK_PORT}/__admin/mappings" -d "@${fixture}" --header "Content-Type: application/json"
    echo "Added fixture: ${fixture}"
done

echo "${OK_COLOR}Uploading upstreams mock fixture:${NO_COLOR}"
curl -X DELETE --silent --show-error --fail --output /dev/null "http://localhost:${UPSTREAMS_MOCK_PORT}/__admin/mappings"
for fixture in assets/stubs/upstreams/*; do
    curl --silent --show-error --fail --output /dev/null "http://localhost:${UPSTREAMS_MOCK_PORT}/__admin/mappings" -d "@${fixture}" --header "Content-Type: application/json"
    echo "Added fixture: ${fixture}"
done
