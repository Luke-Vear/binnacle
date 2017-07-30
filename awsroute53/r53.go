package r53

import (
	"log"

	"sort"

	"github.com/Luke-Vear/binnacle/dnsprovider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

// Route53DNS struct embeds structs' fields and methods binnacle/dnsprovider package.
type Route53DNS struct {
	dns.Conf
	dns.Data
}

// UPSERT method updates or inserts record data in route53 dns.
func (r *Route53DNS) UPSERT() error {

	// Create session
	sess := session.Must(session.NewSession())
	svc := route53.New(sess)

	// Build up resource record slice to upsert
	var resourceRecordSlice []*route53.ResourceRecord

	for _, ipaddr := range r.IngressIPs() {
		rr := route53.ResourceRecord{
			Value: aws.String(ipaddr),
		}
		resourceRecordSlice = append(resourceRecordSlice, &rr)
	}

	params := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(r.ZoneID()),
		ChangeBatch: &route53.ChangeBatch{
			Comment: aws.String(r.Info()),
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name:            aws.String(r.RecordName()),
						Type:            aws.String("A"),
						TTL:             aws.Int64(r.TTL()),
						ResourceRecords: resourceRecordSlice,
					},
				},
			},
		},
	}

	resp, err := svc.ChangeResourceRecordSets(params)
	if err != nil {
		return err

	}
	log.Println(resp)
	return nil
}

// GET method goes to route53 dns and gets the current data for the  records stored there.
func (r *Route53DNS) GET() error {

	// Create session
	sess := session.Must(session.NewSession())
	svc := route53.New(sess)

	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(r.ZoneID()),
		StartRecordName: aws.String(r.RecordName()),
		MaxItems:        aws.String("1"),
	}

	resp, err := svc.ListResourceRecordSets(params)
	if err != nil {
		return err
	}

	// Extract the IP addresses from the response and append them to temporary
	// slice, then set recordIPs to the temporary slice.
	if len(resp.ResourceRecordSets) == 0 {
		return nil
	}

	var rips []string
	for i := range resp.ResourceRecordSets[0].ResourceRecords {
		rips = append(rips, *resp.ResourceRecordSets[0].ResourceRecords[i].Value)
	}

	sort.Strings(rips)
	r.SetRecordIPs(rips)

	return nil
}
