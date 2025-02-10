package env

import (
	"crypto/x509"
	"os"
	"strings"
)

func TrustedCertsPoolFromEnv() (*x509.CertPool, error) {
	certs := strings.Split(TrustedCerts, ",")

	caCertPool := x509.NewCertPool()

	for _, cert := range certs {
		caCert, err := os.ReadFile(cert)
		if err != nil {
			return nil, err
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}

	return caCertPool, nil
}
