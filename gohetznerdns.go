package gohetznerdns

// Hetzner DNS Public API interface entry interface
// Exposes DNS and Record service to manage DNS Zone and records.
// See api documentation for more information [https://dns.hetzner.com/api-docs]
type HetznerDNS interface {

	// Configures API Base URL
	SetBaseURL(baseUrl string) error

	// Configures API Token
	SetToken(token string) error

	// Returns Zone Service
	GetZoneService() ZoneService

	// Returns Record Service
	GetRecordService() RecordService
}

type hetznerDNS struct {
	client        *client
	ZoneService   ZoneService
	RecordService RecordService
}

var _ HetznerDNS = &hetznerDNS{}

// Creates new Hetzner DNS Public API Client with the given token
// see [HetznerDNS.setToken] to update token after creation
func NewClient(token string) (HetznerDNS, error) {
	cli := newClient()
	dns := &hetznerDNS{
		client:        cli,
		ZoneService:   &zoneService{client: cli},
		RecordService: &recordService{client: cli},
	}
	return dns, dns.SetToken(token)
}

func (dns *hetznerDNS) SetBaseURL(baseUrl string) error {
	return dns.client.setBaseURL(baseUrl)
}

func (dns *hetznerDNS) SetToken(token string) error {
	if err := validateNotEmpty("token", &token); err != nil {
		return err
	}
	dns.client.setToken(token)
	return nil
}

func (dns *hetznerDNS) GetZoneService() ZoneService {
	return dns.ZoneService
}

func (dns *hetznerDNS) GetRecordService() RecordService {
	return dns.RecordService
}
