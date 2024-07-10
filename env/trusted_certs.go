package env

import (
	"crypto/x509"
	"github.com/allape/goenv"
	"os"
	"strings"
)

func TrustedCertsPoolFromEnv() (*x509.CertPool, error) {
	certs := strings.Split(goenv.Getenv(TrustedCerts, ""), ",")

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
