#!/bin/bash

export WORKSPACE_ROOT=$(realpath $(dirname ${BASH_SOURCE[0]})/../)
export APPS_ROOT=$(realpath $(dirname ${BASH_SOURCE[0]})/../../apps)

# ===== WORKSPACE ======================================================================================================
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

elc start mysql proxy

# ===== BASKET =========================================================================================================
elc -c mysql mysql -uroot -proot -e 'CREATE DATABASE IF NOT EXISTS basket;'

cp -n ${APPS_ROOT}/basket/.env{.example,}
elc -c basket compose build
elc -c basket run composer install
elc -c basket run npm install
elc -c basket run php artisan migrate:fresh
elc -c basket run php artisan db:seed

elc restart basket

# ===== PRODUCT ========================================================================================================
elc -c mysql mysql -uroot -proot -e 'CREATE DATABASE IF NOT EXISTS product;'

cp -n ${APPS_ROOT}/product/.env{.example,}
elc -c product compose build
elc -c product run composer install
elc -c product run npm install
elc -c product run php artisan migrate:fresh
elc -c product run php artisan db:seed

elc restart product

# ===== API ============================================================================================================
cp -n ${APPS_ROOT}/api/.env{.example,}
elc -c api compose build
elc -c api run composer install
elc -c api run npm install

elc restart api

# ===== API HYPERF =====================================================================================================
cp -n ${APPS_ROOT}/api-hyperf/.env{.example,}
elc -c api-hyperf compose build
elc -c api-hyperf run composer install
