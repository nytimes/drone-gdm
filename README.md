drone-gdm
=========

![Travis CI Build Badge](https://travis-ci.org/nytimes/drone-gdm.svg?branch=master)

A simple drone plugin which wraps [Google Deployment Manager](https://cloud.google.com/deployment-manager/docs/).

Features:
 * _Create or acquire semantics_ - Set the desired `state` and the plugin
   determines whether to create, update, or delete.
 * Support for [GDM Beta Composite Types](https://cloud.google.com/deployment-manager/docs/configuration/templates/create-composite-types)

Usage
-----
The bulk of the input parameters are mapped directly to `gcloud` command options.
Documentation follows for the handful of parameters which are particular to `drone-gdm`.

#### State and Action
The `state` can be one of `absent`, `present`, or `latest`.
Plugin "state" | Object Exists? | Action
-------------- | -------------- | -----------
present        | no             | `create`
present        | yes            | _no action_
latest         | no             | `create`
latest         | yes            | `update`
absent         | no             | _no action_
absent         | yes            | `delete`

The specific `action` selected by drone-gdm can be provided to your template
as a property, by specifying `passAction: true`. This will invoke your
configuration or template with `--properties=action:<action from table above>`.

### Example
```Yaml
deploy:
  gdm:
    # Indicate where to acquire the image:
    image: nytimes/drone-gdm

    # Provided JSON auth token (from drone secrets):
    gcloudPath: /bin/gcloud   # path to gcloud executable
    verbose: false            # (optional)
    dryRun: false             # (optional)
    token: >
      $$GOOGLE_JSON_CREDENTIALS
    project: my-gcp-project   # Da--project
    preview: false            # --preview
    async: false              # --async

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
