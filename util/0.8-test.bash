#!/usr/bin/env bash

export DRONE_DIR="${thisdir}/.."
export PLUGIN_VERBOSE="true"
export TOKEN="{\"type\":\"service_account\",\"project_id\":\"my-cloud-proj\",\"private_key_id\":\"abcdefghijklmnopqrstuvwxyz\",\"private_key\":\"-----BEGINPRIVATEKEY-----\nStandardPrivateKeyData\n-----ENDPRIVATEKEY-----\n\",\"client_email\":\"my-svc-user@my-cloud-proj.iam.gserviceaccount.com\",\"client_id\":\"12345678901234567890\",\"auth_uri\":\"https://accounts.google.com/o/oauth2/auth\",\"token_uri\":\"https://accounts.google.com/o/oauth2/token\",\"auth_provider_x509_cert_url\":\"https://www.googleapis.com/oauth2/v1/certs\",\"client_x509_cert_url\":\"https://www.googleapis.com/robot/v1/metadata/x509/my-svc-user%40my-cloud-proj.iam.gserviceaccount.com\"}",
export PLUGIN_PROJECT="my-cloud-proj"
export PLUGIN_PREVIEW="false"
export PLUGIN_ASYNC="false"
export PLUGIN_DRYRUN="true"

read -r -d '' PLUGIN_CONFIGURATIONS <<'EOF'
[
	{
		"name": "my-deployment1",
		"group": "deployment",
		"state": "latest",
		"path": "./my-deployment1.yml",
		"automaticRollbackOnError": true,
		"createPolicy": "CREATE_OR_ACQUIRE",
		"deletePolicy": "DELETE",
		"properties": {
			"p1": "val1",
			"p2": "val2",
            "accessControl": {
                "bindings": [
                    {
                        "role": "roles/storage.objectAdmin",
                        "members": [
                          "my-acct@gmail.com",
                          "another-acct@gmail.com"
                        ]
                    },
                    {
                        "role": "roles/storage.objectViewer",
                        "members": [
                          "some-user@gmail.com",
                          "another-user@gmail.com"
                        ]
                    }
                ]
            }
		},
		"labels": {
			"l1": "val1",
			"l2": "val2"
		}
	},
	{
		"name": "my-deployment2",
		"group": "deployment",
		"state": "present",
		"path": "./my-deployment2.yml",
		"automaticRollbackOnError": true,
		"createPolicy": "CREATE_OR_ACQUIRE",
		"deletePolicy": "DELETE",
		"properties": {
			"p1": "val1",
			"p2": "val2"
		},
		"labels": {
			"l1": "val1",
			"l2": "val2"
		}
	},
	{
		"name": "my-deployment3",
		"group": "deployment",
		"state": "absent",
		"path": "./my-deployment3.yml"
	},
	{
		"name": "my-composite1",
		"group": "composite",
		"state": "present",
		"path": "./my-composite1.yml",
		"status": "SUPPORTED",
		"labels": {
			"l1": "val1",
			"l2": "val2"
		}
	},
	{
		"name": "my-composite2",
		"group": "composite",
		"state": "absent",
		"path": "./my-composite2.yml",
		"status": "DEPRECATED",
		"labels": {
			"l1": "val1",
			"l2": "val2"
		}
	}
]
EOF

export PLUGIN_CONFIGURATIONS
$@

# EOF

