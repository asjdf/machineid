//go:build linux
// +build linux

package machineid

import "os"

const (
	// k8sCAPath is the default path for kubernetes service account certificate.
	k8sCAPath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	// dbusPath is the default path for dbus machine id.
	dbusPath = "/var/lib/dbus/machine-id"
	// dbusPathEtc is the default path for dbus machine id located in /etc.
	// Some systems (like Fedora 20) only know this path.
	// Sometimes it's the other way round.
	dbusPathEtc = "/etc/machine-id"
)

// machineID returns the uuid specified at `/var/lib/dbus/machine-id` or `/etc/machine-id`.
// If there is an error reading the files an empty string is returned.
// See https://unix.stackexchange.com/questions/144812/generate-consistent-machine-unique-id
func machineID() (string, IDType, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		// k8s cluster detected
		ca, err := os.ReadFile(k8sCAPath)
		if err == nil && len(ca) > 0 {
			return protect("k8s", string(ca)), TypeKubernetes, nil
		}
		// fallback
	}

	id, err := os.ReadFile(dbusPath)
	if err != nil {
		// try fallback path
		id, err = os.ReadFile(dbusPathEtc)
	}
	if err != nil {
		return "", TypeUnknown, err
	}
	return trim(string(id)), TypeStandalone, nil
}
