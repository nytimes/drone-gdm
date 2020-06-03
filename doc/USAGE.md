# Drone-GDM: Usage

The bulk of the input parameters are mapped directly to `gcloud` command options.
Documentation follows for the handful of parameters which are particular to `drone-gdm`.

### State and Action
The `state` can be one of `absent`, `present`, or `latest`.

| Plugin "state" | Object Exists? | Action      |
| -------------- | -------------- | ----------- |
| present        | no             | `create`    |
| present        | yes            | _no action_   |
| latest         | no             | `create`    |
| latest         | yes            | `update`    |
| absent         | no             | _no action_   |
| absent         | yes            | `delete`    |

The specific `action` selected by drone-gdm can be provided to your template
as a property, by specifying `passAction: true`. This will invoke your
configuration or template with `--properties=action:<action from table above>`.

### Variables
To circumvent data-type limitations imposed by the passing of properties via the
deployment manager `--properties` option, external configuration files (see the
[examples](./EXAMPLES.md) for more info), are processed first as [golang templates](https://golang.org/pkg/text/template/) with the following top-level interfaces available for variable interpolation:
 - `.drone` - Drone environment variables provided by the CI system during plugin invocation
 - `.plugin` - Plugin parameters passed via environment during plugin invocation
 - `.context` - Any variables defined in the `vars` section of the plugin invocation
 - `.config` - Any variables defined in the `vars` section of the configuration definition
 - `.properties` - Variables defined in the `properties` section of the configuration definition
 - `.gdm` - A dictionary containing:
   - `name` - entity name for the configuration/template/composite
   - `status` - the entity status (e.g. DEPRECATED, EXPERIMENTAL, SUPPORTED)
   - `project` - the GCP project name
   - `action` - the gcloud "action" parameter (i.e. `create`, `update`, or `delete`)



