package hetznerdns

type HetznerDNS struct {
	client        *client
	ZoneService   ZoneService
	RecordService RecordService
}

func NewClient() *HetznerDNS {
	cli := newClient()
	return &HetznerDNS{
		client:        cli,
		ZoneService:   &zoneService{client: cli},
		RecordService: &recordService{client: cli},
	}
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
