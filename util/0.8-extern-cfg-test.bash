#!/usr/bin/env bash

thispath="${0}"
thisdir="${thispath%/*}"

export DRONE_DIR="${thisdir}/.."
export PLUGIN_VERBOSE="true"
export TOKEN="{\"type\":\"service_account\",\"project_id\":\"my-cloud-proj\",\"private_key_id\":\"abcdefghijklmnopqrstuvwxyz\",\"private_key\":\"-----BEGINPRIVATEKEY-----\nStandardPrivateKeyData\n-----ENDPRIVATEKEY-----\n\",\"client_email\":\"my-svc-user@my-cloud-proj.iam.gserviceaccount.com\",\"client_id\":\"12345678901234567890\",\"auth_uri\":\"https://accounts.google.com/o/oauth2/auth\",\"token_uri\":\"https://accounts.google.com/o/oauth2/token\",\"auth_provider_x509_cert_url\":\"https://www.googleapis.com/oauth2/v1/certs\",\"client_x509_cert_url\":\"https://www.googleapis.com/robot/v1/metadata/x509/my-svc-user%40my-cloud-proj.iam.gserviceaccount.com\"}",
export PLUGIN_PROJECT="my-cloud-proj"
export PLUGIN_PREVIEW="false"
export PLUGIN_ASYNC="false"
export PLUGIN_DRYRUN="true"
export PLUGIN_CONFIGFILE="${thisdir}/0.8-extern-cfg-test.yml"

$@

# EOF

