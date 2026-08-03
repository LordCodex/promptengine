package security

import (
	"regexp"
	"strings"
)

var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(api[_-]?key|access[_-]?token|auth[_-]?token|secret|password|passwd|private[_-]?key)\s*[:=]\s*['"]?([^\s'"]+)`),
	regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----[\s\S]*?-----END [A-Z ]*PRIVATE KEY-----`),
	regexp.MustCompile(`(?i)(bearer\s+)[a-z0-9._\-]{20,}`),
}

var sensitivePathMarkers = []string{
	".env",
	".npmrc",
	".pypirc",
	".netrc",
	"id_rsa",
	"id_dsa",
	"id_ecdsa",
	"id_ed25519",
	"credentials",
	"secrets",
}

func RedactSecrets(content string) (string, bool) {
	redacted := content
	changed := false
	for _, pattern := range secretPatterns {
		next := pattern.ReplaceAllStringFunc(redacted, func(match string) string {
			changed = true
			lower := strings.ToLower(match)
			if strings.HasPrefix(lower, "bearer ") {
				return match[:7] + "[REDACTED]"
			}
			if strings.Contains(match, "PRIVATE KEY-----") {
				return "[REDACTED PRIVATE KEY]"
			}
			parts := strings.FieldsFunc(match, func(r rune) bool { return r == ':' || r == '=' })
			if len(parts) > 0 {
				sep := "="
				if strings.Contains(match, ":") && !strings.Contains(match, "=") {
					sep = ":"
				}
				return strings.TrimSpace(parts[0]) + sep + "[REDACTED]"
			}
			return "[REDACTED]"
		})
		redacted = next
	}
	return redacted, changed
}

func ContainsSecret(content string) bool {
	for _, pattern := range secretPatterns {
		if pattern.MatchString(content) {
			return true
		}
	}
	return false
}

func IsSensitivePath(path string) bool {
	lower := strings.ToLower(strings.ReplaceAll(path, "\\", "/"))
	base := lower
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		base = base[idx+1:]
	}
	for _, marker := range sensitivePathMarkers {
		if base == marker || strings.Contains(lower, "/"+marker+"/") || strings.Contains(base, marker) {
			return true
		}
	}
	return false
}
