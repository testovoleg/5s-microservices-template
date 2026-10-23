package httpUtils

import "strings"

// FilenameFromContentDisposition parses filename*= / filename= from a Content-Disposition header.
func FilenameFromContentDisposition(cd string) string {
	if cd == "" {
		return ""
	}

	filename := ""

	if idx := strings.Index(cd, "filename*="); idx != -1 {
		start := idx + len("filename*=")
		end := strings.Index(cd[start:], "\"")
		if end != -1 {
			filename = cd[start : start+end]
		}
	}

	if filename == "" {
		if idx := strings.Index(cd, "filename="); idx != -1 {
			start := idx + len("filename=")
			if start < len(cd) && cd[start] == '"' {
				start++
				end := strings.Index(cd[start:], "\"")
				if end != -1 {
					filename = cd[start : start+end]
				}
			} else {
				end := strings.IndexAny(cd[start:], " ;")
				if end != -1 {
					filename = cd[start : start+end]
				} else {
					filename = cd[start:]
				}
			}
		}
	}

	return filename
}
