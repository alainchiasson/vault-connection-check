package main

import (
	"log"
	"net"
	"os"

	// "context"
	// "time"
	"crypto/tls"
	// "log/slog"
)

var version string

type DnsTest struct {
	HostName          string
	ExpectedIpAddress string
}

type CNAMETest struct {
	HostName      string
	ExpectedCNAME string
}

func show_env(key string) {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("%s not set\n", key)
	} else {
		log.Printf("%s=%s\n", key, val)
	}
}

func get_proxy() {

}

func test_https_no_proxy(endpoint string) {
	// Connect https
	// Skip the server verification - we want to print the certificate
	//

	log.Printf("==> Verifying connectivity to : %s", endpoint)

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", endpoint, conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	log.Println("Connected...")
	log.Println("Printing received certificate information ...")

	certs := conn.ConnectionState().PeerCertificates
	for _, cert := range certs {
		log.Println("---------")
		log.Printf("Common Name: %s \n", cert.Issuer.CommonName)
		log.Printf("Issuer Name: %s\n", cert.Issuer)
		log.Printf("Expiry: %s \n", cert.NotAfter.Format("2006-January-02"))
	}

	log.Println(" ==< Certificate End ...")

}

func test_https_with_proxy(endpoint string, testproxy string) {
	// Connect https
	// Skip the server verification - we want to print the certificate
	//

	log.Printf("==> Verifying connectivity to : %s via proxy %s", endpoint, testproxy)

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", endpoint, conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	log.Println("Connected...")
	log.Println("Printing received certificate information ...")

	certs := conn.ConnectionState().PeerCertificates
	for _, cert := range certs {
		log.Println("---------")
		log.Printf("Common Name: %s \n", cert.Issuer.CommonName)
		log.Printf("Issuer Name: %s\n", cert.Issuer)
		log.Printf("Expiry: %s \n", cert.NotAfter.Format("2006-January-02"))
	}

	log.Println(" ==< Certificate End ...")

}

func test_dns_resolution(DnsNames []DnsTest) {

	log.Println("==> Verifying DNS name resolution... ")

	for _, DnsName := range DnsNames {
		ip, _ := net.LookupHost(DnsName.HostName)

		log.Println(DnsName.HostName, DnsName.ExpectedIpAddress, ip)
	}

}

func test_cname_resolution(CNames []CNAMETest) {

	log.Println("==> Verifying CNAME Mapping ... ")
	for _, CName := range CNames {
		record, _ := net.LookupCNAME(CName.HostName)

		log.Println(CName.HostName, CName.ExpectedCNAME, record)
	}
}

func main() {

	log.Println("standard logger")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var DnsNames = []DnsTest{
		DnsTest{
			HostName:          "home.home.chiasson.org",
			ExpectedIpAddress: "192.168.50.2",
		},
		DnsTest{
			HostName:          "oscp-01.home.chiasson.org",
			ExpectedIpAddress: "192.168.50.20",
		},
		DnsTest{
			HostName:          "api.hosc.home.chiasson.org",
			ExpectedIpAddress: "192.168.50.20",
		},
	}

	var CNames = []CNAMETest{
		CNAMETest{
			HostName:      "api.hosc.home.chiasson.org",
			ExpectedCNAME: "oscp-01.home.chiasson.org.",
		},
	}

	test_dns_resolution(DnsNames)

	test_cname_resolution(CNames)

	test_https_no_proxy("www.google.com:443")

	show_env("https_proxy")
	show_env("HTTPS_PROXY")
	show_env("no_proxy")
	show_env("NO_PROXY")

	test_https_with_proxy("www.google.com:443", "proxy:80")
	test_https_with_proxy("www.google.com:443", "proxy:80")

	// 	request := gorequest.New()
	// resp, body, errs:= request.Proxy("http://proxy:999").Get("http://example.com").End()
	// resp2, body2, errs2 := request.Proxy("http://proxy2:999").Get("http://example2.com").End()

}
