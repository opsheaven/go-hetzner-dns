package gohetznerdns

type HetznerDNS struct {
	client        *client
	ZoneService   ZoneService
	RecordService RecordService
}

func NewClient(token string) (*HetznerDNS, error) {
	cli := newClient()
	dns := &HetznerDNS{
		client:        cli,
		ZoneService:   &zoneService{client: cli},
		RecordService: &recordService{client: cli},
	}
	return dns, dns.SetToken(token)
}

func (dns *HetznerDNS) SetBaseURL(baseUrl string) error {
	return dns.client.setBaseURL(baseUrl)
}

func (dns *HetznerDNS) SetToken(token string) error {
	if err := validateNotEmpty("token", token); err != nil {
		return err
	}
	dns.client.setToken(token)
	return nil
}
