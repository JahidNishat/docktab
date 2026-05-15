package docker

import "strings"

func shortID(id string, n int) string {
	if len(id) <= n {
		return id
	}
	return id[:n]
}

func firstContainerName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return strings.TrimPrefix(names[0], "/")
}

func splitRepoTag(ref string) (string, string) {
	if ref == "" || ref == "<none>:<none>" {
		return "<none>", "<none>"
	}

	lastSlash := strings.LastIndex(ref, "/")
	lastColon := strings.LastIndex(ref, ":")

	if lastColon > lastSlash {
		return ref[:lastColon], ref[lastColon+1:]
	}

	return ref, "latest"
}
