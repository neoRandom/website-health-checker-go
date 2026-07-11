package util

import "crypto/tls"

func IsTLSDeprecated(tlsVersion uint16) (vName string, depr bool) {
	switch tlsVersion {
	case tls.VersionTLS13:
		vName = "TLS 1.3"
		depr = false
	case tls.VersionTLS12:
		vName = "TLS 1.2"
		depr = false
	case tls.VersionTLS11:
		vName = "TLS 1.1 (Deprecated)"
		depr = true
	case tls.VersionTLS10:
		vName = "TLS 1.0 (Deprecated)"
		depr = true
	default:
		vName = "Unknown Protocol"
		depr = true
	}

	return
}
