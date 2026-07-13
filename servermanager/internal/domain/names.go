package domain

const (
	RuntimeRunc = "runc"
	RuntimeKata = "kata"
)

// ValidRuntime reports whether r is on the server-side runtime allowlist —
// the value ends up in HostConfig.Runtime and must never be free text.
func ValidRuntime(r string) bool {
	return r == RuntimeRunc || r == RuntimeKata
}

// Container and network names are derived from the UUIDs the api sends,
// never from user input.

func AppContainerName(appID string) string { return "stt-app-" + appID }

func DBContainerName(dbID string) string { return "stt-db-" + dbID }

func AppNetworkName(appID string) string { return "stt-net-" + appID }

// DBNetworkName deliberately does NOT share the app networks' "stt-net-"
// prefix: removing an app sweeps every attached "stt-net-*" network, and the
// database network must survive its linked app. The database usecases own the
// "stt-dbnet-*" lifecycle exclusively.
func DBNetworkName(dbID string) string { return "stt-dbnet-" + dbID }

// DBVolumeName is the named volume holding the database's data directory. It
// survives container replacement and is removed only by an explicit database
// delete.
func DBVolumeName(dbID string) string { return "stt-db-data-" + dbID }

// AppImageRepo is the image repository holding every image the manager built
// for one app — image GC lists by this prefix.
func AppImageRepo(appID string) string { return "stt-app-" + appID }

// AppImageRef tags an image with the deployment that built it. The first 8
// characters of a UUID are hex, so the tag is always valid.
func AppImageRef(appID, deploymentID string) string {
	tag := deploymentID
	if len(tag) > 8 {
		tag = tag[:8]
	}
	return AppImageRepo(appID) + ":" + tag
}
