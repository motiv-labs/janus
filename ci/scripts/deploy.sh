#!/bin/sh
set -e
base_path="$(pwd)"
# Requires $GROUP_NAME to be defined
# Requires $TAR_FILE to be defined
# Optional $INPUT defaults to release
# Optional ${DEPLOYMENT_ENVIRONMENT} default to staging.ini
RC_VERSION=$(cat ${VERSION_FILE})
INPUT=${INPUT-release}
CI_INVENTORY=${DEPLOYMENT_ENVIRONMENT-staging.ini}
mkdir -p /root/.ssh/ && chmod 0600 /root/.ssh
set +x
echo "${PRIVATE_KEY}" > "${ANSIBLE_PRIVATE_KEY_FILE}" && chmod 0600 "${ANSIBLE_PRIVATE_KEY_FILE}"
set -x
# Check we have right params
[ ! -f "${base_path}/${TAR_FILE}" ] && echo "tar file not found '${base_path}/${TAR_FILE}'" && exit 1
zcat ${TAR_FILE} | tar -xf -
cd "automation-artifcat/plays"
export ANSIBLE_FORCE_COLOR=true
ansible-playbook -i ../${DEPLOYMENT_ENVIRONMENT} ${GROUP_NAME}.yml -u policy -vvvv -t deployment -e deployment_force=true