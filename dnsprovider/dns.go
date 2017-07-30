package dns

// Provider interface allows choice of external DNS providers.
type Provider interface {
	UPSERT() error
	GET() error

	Provider() string
	SetProvider(s string)

	IngressName() string
	SetIngressName(s string)

	NameSpace() string
	SetNameSpace(s string)

	RecordName() string
	SetRecordName(s string)

	ZoneID() string
	SetZoneID(s string)

	TTL() int64
	SetTTL(i int64)

	Info() string
	SetInfo(s string)

	IngressIPs() []string
	SetIngressIPs(ss []string)

	RecordIPs() []string
	SetRecordIPs(ss []string)
}

// Conf struct contains information on the desired state of the DNS record.
// The struct also identifies where to find the record.
type Conf struct {
	provider    string // The name of the DNS provider.
	ingressName string // The name of the ingress in kuberenetes.
	nameSpace   string // The kubernetes namespace.
	recordName  string // The name of the DNS record that needs to be kept current.
	zoneID      string // Zone identifier for DNS provider interaction.
	ttl         int64  // Time to live of the DNS record.
}

func NewConf(p, i, n, r, z string, t int64) Conf {
	return Conf{
		provider:    p,
		ingressName: i,
		nameSpace:   n,
		recordName:  r,
		zoneID:      z,
		ttl:         t,
	}
}

func (c *Conf) Provider() string {
	return c.provider
}

func (c *Conf) SetProvider(s string) {
	c.provider = s
}

func (c *Conf) IngressName() string {
	return c.ingressName
}

func (c *Conf) SetIngressName(s string) {
	c.ingressName = s
}

func (c *Conf) NameSpace() string {
	return c.nameSpace
}

func (c *Conf) SetNameSpace(s string) {
	c.nameSpace = s
}

func (c *Conf) RecordName() string {
	return c.recordName
}

func (c *Conf) SetRecordName(s string) {
	c.recordName = s
}

func (c *Conf) ZoneID() string {
	return c.zoneID
}

func (c *Conf) SetZoneID(s string) {
	c.zoneID = s
}

func (c *Conf) TTL() int64 {
	return c.ttl
}

func (c *Conf) SetTTL(i int64) {
	c.ttl = i
}

// Data struct contains almost live information about the records.
type Data struct {
	info       string   // Info for attaching to dns change.
	ingressIPs []string // Poll k8s api for current ingress apis.
	recordIPs  []string // Check local cache or request from dns provider current records.
}

func (d *Data) Info() string {
	return d.info
}

func (d *Data) SetInfo(s string) {
	d.info = s
}

func (d *Data) IngressIPs() []string {
	return d.ingressIPs
}

func (d *Data) SetIngressIPs(ss []string) {
	d.ingressIPs = ss
}

func (d *Data) RecordIPs() []string {
	return d.recordIPs
}

func (d *Data) SetRecordIPs(ss []string) {
	d.recordIPs = ss
}
