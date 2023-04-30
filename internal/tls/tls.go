package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func GetTlsConfig(prefix string) (*tls.Config, error) {
	certVar := fmt.Sprintf("%s_TLS_CERTIFICATE", prefix)
	keyVar := fmt.Sprintf("%s_TLS_KEY", prefix)
	caVar := fmt.Sprintf("%s_TLS_CA_CERTIFICATE", prefix)

	certPem, ok := os.LookupEnv(certVar)
	if !ok {
		return nil, fmt.Errorf("not set: %s", certVar)
	}

	keyPem, ok := os.LookupEnv(keyVar)
	if !ok {
		return nil, fmt.Errorf("not set: %s", keyVar)
	}

	cert, err := tls.X509KeyPair([]byte(certPem), []byte(keyPem))
	if err != nil {
		return nil, fmt.Errorf("failed to parse keypair: %v", err)
	}

	caPem, ok := os.LookupEnv(caVar)
	if !ok {
		return nil, fmt.Errorf("not set: %s", caVar)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM([]byte(caPem)) {
		return nil, fmt.Errorf("failed to parse CA certificate")
	}

	return &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		RootCAs:      certPool,
	}, nil
}
