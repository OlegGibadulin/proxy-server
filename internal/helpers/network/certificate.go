package network

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func GenerateCertificate(url *url.URL) (tls.Certificate, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return tls.Certificate{}, err
	}
	scriptsDir := filepath.Join(pwd, "scripts")
	certsDir := filepath.Join(scriptsDir, "certs")

	genPath := filepath.Join(scriptsDir, "gen_cert.sh")
	certKeyPath := filepath.Join(scriptsDir, "cert.key")

	certName := fmt.Sprintf("%s.crt", url.Scheme)
	certPath := filepath.Join(certsDir, certName)

	_, errStat := os.Stat(certPath)
	if os.IsNotExist(errStat) {
		genCmd := exec.Command(genPath, url.Scheme, strconv.Itoa(rand.Int()))
		if _, err := genCmd.CombinedOutput(); err != nil {
			return tls.Certificate{}, err
		}
	}

	cert, err := tls.LoadX509KeyPair(certPath, certKeyPath)
	if err != nil {
		return tls.Certificate{}, err
	}
	return cert, nil
}
