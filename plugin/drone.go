package plugin

import (
	"os"
	"strings"
)

type MapStrStr map[string]string

func DroneVars() MapStrStr {
	return varsFromEnvPrefix("DRONE_")
}

func PluginVars() MapStrStr {
	return varsFromEnvPrefix("PLUGIN_")
}

func varsFromEnvPrefix(prefix string) MapStrStr {
	envVars := make(MapStrStr)
	for _, e := range os.Environ() {
		comps := strings.SplitN(e, "=", 2)
		if (len(comps) > 1) && strings.HasPrefix(comps[0], prefix) {
			// TODO: Move this func to another translation unit + add filter
			vname := strings.ToLower(strings.TrimPrefix(comps[0], prefix))
			envVars[vname] = comps[1]
		}
	}
	return envVars
}
