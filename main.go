package main

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/Luke-Vear/binnacle/awsroute53"
	"github.com/Luke-Vear/binnacle/dnsprovider"
	"github.com/Luke-Vear/binnacle/gclouddns"
)

// Main goroutine loops indefinitely, polling kubernetes API server for ingress IPs.
// It then checks if ingress IPs are different than are in DNS records of DNS
// provider and if DNS records are different, updates them with ingress IPs.
func main() {

	// Constructs dns.Provider based on config read from kubernetes config map.
	dnsConf, err := getKubeConfigMap()
	if err != nil {
		log.Panic(err)
	}
	dp := newProvider(dnsConf)

	for {
		// Check kubernetes API server for ingress IPs.
		iips, err := getKubeIngressIPs(dp)
		if err != nil {
			log.Println(err)
			time.Sleep(10 * time.Second)
			continue
		}
		dp.SetIngressIPs(iips)

		// Check record cache/fetch external DNS records.
		if dp.RecordIPs() == nil {
			log.Println("DNS providers records field nil, fetching from provider.")
			updateCachedRecords(dp)
		}

		log.Println("DNS Record IPs:", dp.RecordIPs())
		log.Println("Ingress IPs:", dp.IngressIPs())

		// If provider DNS records do not match ingress records, update DNS.
		if !reflect.DeepEqual(dp.IngressIPs(), dp.RecordIPs()) {
			err := updateProviderRecords(dp)
			if err != nil {
				log.Println("updateProviderRecords issue,", err)
			}

			err = updateCachedRecords(dp)
			if err != nil {
				log.Println("updateCachedRecords issue,", err)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

// Returns a provider based on passed in dns.Conf.provider field.
func newProvider(c dns.Conf) dns.Provider {
	switch c.Provider() {
	case "r53":
		return &r53.Route53DNS{
			Conf: c,
		}
	case "gc":
		return &gc.GCloudDNS{
			Conf: c,
		}
	}
	return nil
}

// Updates the locally cached dns.Data.recordIPs from the external DNS provider.
func updateCachedRecords(dp dns.Provider) error {

	if err := dp.GET(); err != nil {
		return err
	}
	log.Println("Updated local DNS record cache.")
	return nil
}

// Updates the IP addresses on the records stored in the external DNS provider to
// match the IP addresses on the dns.Data.ingressIPs field.
func updateProviderRecords(dp dns.Provider) error {

	dp.SetInfo(fmt.Sprintf("Changing from %v to %v.", dp.RecordIPs(), dp.IngressIPs()))
	if err := dp.UPSERT(); err != nil {
		return err
	}
	log.Println("External DNS provider records updated to:", dp.IngressIPs())
	return nil
}
