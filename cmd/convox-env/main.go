package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/lob/rack/pkg/crypt"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "usage: convox-env <command>\n")
		os.Exit(1)
	}

	cenv, err := fetchConvoxEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: could not fetch convox env: %s\n", err)
		os.Exit(1)
	}

	env := append(os.Environ(), cenv...)

	err = exec(os.Args[1], os.Args[2:], env)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func fetchConvoxEnv() ([]string, error) {
	eu := os.Getenv("CONVOX_ENV_URL")
	if eu == "" {
		return nil, nil
	}

	u, err := url.Parse(eu)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "s3" {
		return nil, fmt.Errorf("unrecognized env url")
	}

	res, err := S3().GetObject(&s3.GetObjectInput{
		Bucket: aws.String(u.Host),
		Key:    aws.String(u.Path),
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if k := os.Getenv("CONVOX_ENV_KEY"); k != "" {
		dec, err := crypt.New().Decrypt(k, data)
		if err != nil {
			return nil, err
		}

		data = dec
	}

	env := []string{}

	sc := bufio.NewScanner(bytes.NewReader(data))

	allowed := map[string]bool{}

	if ev := os.Getenv("CONVOX_ENV_VARS"); ev != "" {
		for _, v := range strings.Split(ev, ",") {
			allowed[v] = true
		}
	}

	for sc.Scan() {
		if s := sc.Text(); s != "" {
			if allowed["*"] || allowed[strings.Split(s, "=")[0]] {
				env = append(env, s)
			}
		}
	}

	return env, nil
}

func S3() *s3.S3 {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(caCertificates))

	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}

	client := http.DefaultClient
	client.Transport = tr

	return s3.New(session.New(), &aws.Config{
		Region:           aws.String(os.Getenv("AWS_REGION")),
		HTTPClient:       client,
		S3ForcePathStyle: aws.Bool(true),
	})
}
