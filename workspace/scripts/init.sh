#!/bin/bash

export WORKSPACE_ROOT=$(realpath $(dirname ${BASH_SOURCE[0]})/../)
export APPS_ROOT=$(realpath $(dirname ${BASH_SOURCE[0]})/../../apps)

read -r -d '' WS_DIRS << EOF
${WORKSPACE_ROOT}/home/.config
${WORKSPACE_ROOT}/home/.npm
${WORKSPACE_ROOT}/home/composer
${WORKSPACE_ROOT}/home/composer_cache
${WORKSPACE_ROOT}/services/mysql/data
EOF

for WS_DIR in ${WS_DIRS}; do
    mkdir -p ${WS_DIR}
    chmod 777 ${WS_DIR}
done

cp -n ${WORKSPACE_ROOT}/env{.example,}.yaml

sed -E -e "s#APPS_ROOT:.*#APPS_ROOT: ${APPS_ROOT}#" -i ${WORKSPACE_ROOT}/env.yaml