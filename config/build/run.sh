#!/bin/bash

SCRIPT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

REQUIREMENTS_FILE="roles_requirements.yml"
ROLES_PATH="roles"
PLAYBOOK_BUILD="ansible_build.yml"
PLAYBOOK_VERSION_UPDATE="ansible_update_version_file.yml"
VERSION_FILE_TMP="/tmp/version_file.txt"
#COMMIT_ID ## All ready set by codeship

exists () { (
    IFS=:
    for d in $PATH; do
      if test -x "${d}/$1"; then return 0; fi
    done
    return 1
) }

exists "ansible"
ansible_installed="$?"
set -e

### Step 1 Check ansible and install if not
if [ "${ansible_installed}" == "1" ]; then
    echo " $0 > Will install ansible via pip"
    pip install -q ansible 
    pip install -q boto
    pip install -q boto3
else
    echo " $0 > Ansilbe already installed."
fi

# Clear VERSION_FILE_TMP
rm -f ${VERSION_FILE_TMP}

# Change to script dir
cd ${SCRIPT_DIR}

### Step 2 Get external role
echo " $0 >  Create roles dir \"${SCRIPT_DIR}/${ROLES_PATH}\""
mkdir -p ${SCRIPT_DIR}/${ROLES_PATH}
rm -rf ${SCRIPT_DIR}/${ROLES_PATH}/

echo " $0 >  Get external roles ${REQUIREMENTS_FILE} installing "
ansible-galaxy install -r ${REQUIREMENTS_FILE} -p ${SCRIPT_DIR}/${ROLES_PATH}

### Step 3 Run ansible
set +e
#check if checkout is linked to tag. 
tag="$(git describe --exact-match --tags HEAD)"
found_tag="$?"
set -e
if [ "$found_tag" == "0" ]; then
    # Since this a tag we dont need to build we just need to up our release file in live envirionment
    # to point to this tag
    version_file="LATEST-LIVE"
else
    # Build and upload aritifact
    echo " $0 >  Running Ansible "
    ansible-playbook -i 127.0.0.1, -c local ${PLAYBOOK_BUILD} -vvvv
    current_branch="$(git branch --contains $COMMIT_ID | tail -1 | tr -d '[[:space:]]' | tr -d '*')"
    if [ "${current_branch}" == "master" ]; then 
        # if thats the master branch than lets point our dev environment on it
        version_file="LATEST-STAGING" # TODO: change verions file to LATEST-DEV
        echo "$0 > Setting version_file to $version_file"
    else
        echo "$0 > Not on Master branch so skipping version_file. your on ${current_branch}"
    fi
fi

## Step 4 update version file if needed 
if [ ! -z "${version_file}" ]; then
    echo " $0 >  updateing version file ${version_file} to ${COMMIT_ID}"
    echo "${COMMIT_ID}" > ${VERSION_FILE_TMP}
    ansible-playbook -i 127.0.0.1, -c local ${PLAYBOOK_VERSION_UPDATE} -vvvv -e version_file_name=${version_file}
else
    echo " $0 >  SkippVersion file is  ${version_file}"
fi
