Examples: drone-gdm
===================

### Example with Inline Configurations
```Yaml
deploy:
  gdm:
    # Indicate where to acquire the image:
    image: nytimes/drone-gdm:2.0.0

    # Provided JSON auth token (from drone secrets):
    gcloudPath: /bin/gcloud   # path to gcloud executable
    verbose: false            # (optional)
    dryRun: false             # (optional)
    token: >
      $$GOOGLE_JSON_CREDENTIALS
    project: my-gcp-project   # Da--project
    preview: false            # --preview
    async: false              # --async
    vars:
    - myCtxVar: ctxVal1
    - myOtherCtxVar: ctxVal2

    configurations:
    # A "type provider" definition:
    - name:  my-provider
      group: typeprovider
      state: present
      descriptorURL: https://cloudtasks.googleapis.com/$discovery/rest?version=v2beta3
      apiOptions: ./api-options-definition.yaml # path to api options YAML

    # Standard YAML deployment:
    - name: my-deployment
      group: deployment
      state: latest
      description: A basic GDM deployment yaml file which creates some resources
      path: ./my-deployment.yaml
      autoRollbackOnError: false
      createPolicy: CREATE_OR_ACQUIRE # Optional: CREATE_OR_ACQUIRE or CREATE
      deletePolicy: DELETE # Optional: DELETE or ABANDON
      passAction: false # if true, pass action as property, e.g. "action:update"

    # Standard jinja (template) deployment:
    - name:  my-other-deployment
      group: deployment
      state: present
      path: ./my-other-deployment.jinja
      description: A GDM Deployment
      properties:    # mapped to gcloud '--properties=...'
        myvar: myval # can be referenced in jinja as: {{ properties.myvar }}
      labels:        # mapped to '--labels' or '--update-labels', as appropriate
        mylabel: labelval

    # Deploying a composite type:
    - name:  my-composite
      version: beta  # gcloud version to use
      group: composite
      state: present
      path: ./my-composite.jinja
      description: A GDM "Composite Type"
      labels: # mapped to '--labels' or '--update-labels', as appropriate
        mylabel: labelval
      status: SUPPORTED # Optional: SUPPORTED, DEPRECATED, or EXPERIMENTAL
      passAction: false

```

### Example with Inline and External Configurations
```Yaml
deploy:
  gdm:
    # Indicate where to acquire the image:
    image: nytimes/drone-gdm:2.0.0

    # Provided JSON auth token (from drone secrets):
    gcloudPath: /bin/gcloud   # path to gcloud executable
    verbose: false            # (optional)
    dryRun: false             # (optional)
    token: >
      $$GOOGLE_JSON_CREDENTIALS
    project: my-gcp-project   # Da--project
    preview: false            # --preview
    async: false              # --async

    vars:
       prefix: test1
    configFile: my-configurations.yml
    configurations:
    - name:  my-deployment
      group: deployment
      state: present
      path: ./my-deployment.yaml
      description: A GDM Deployment
      properties:    # mapped to gcloud '--properties=...'
        myvar: myval # can be referenced in jinja as: {{ properties.myvar }}
      labels:        # mapped to '--labels' or '--update-labels', as appropriate
        mylabel: labelval
      autoRollbackOnError: false
      createPolicy: CREATE_OR_ACQUIRE # Optional: CREATE_OR_ACQUIRE or CREATE
      deletePolicy: DELETE # Optional: DELETE or ABANDON
      passAction: false # if true, pass action as property, e.g. "action:update"
```

##### my-configurations.yml
``` Yaml
# Parsed as a golang template with variables populated from "vars" above.
- name:  {{.prefix}}-composite
  version: beta  # gcloud version to use
  group: composite
  state: present
  path: ./my-composite.jinja
  description: A GDM "Composite Type"
  labels: # mapped to '--labels' or '--update-labels', as appropriate
    mylabel: labelval
  status: SUPPORTED # Optional: SUPPORTED, DEPRECATED, or EXPERIMENTAL
  passAction: false

```

