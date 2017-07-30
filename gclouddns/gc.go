package gc

import (
	"github.com/Luke-Vear/binnacle/dnsprovider"
)

type GCloudDNS struct {
	dns.Conf
	dns.Data
}

// UPSERT method updates or inserts record data in gcloud dns.
func (g *GCloudDNS) UPSERT() error {
	return nil
}

// GET method goes to gcloud dns and gets the current data for the  records stored there.
func (g *GCloudDNS) GET() error {
	return nil
}
