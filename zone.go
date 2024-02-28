package hetznerdns

import (
	"fmt"
)

const zonesBasePath = "/zones"

type Zone struct {
	Id              string          `json:"id"`
	IsSecondaryDns  bool            `json:"is_secondary_dns"`
	LegacyDnsHost   string          `json:"legacy_dns_host"`
	Name            string          `json:"name"`
	NS              []string        `json:"ns"`
	Owner           string          `json:"owner"`
	Paused          bool            `json:"paused"`
	Permission      string          `json:"permission"`
	Project         string          `json:"project"`
	RecordsCount    int64           `json:"records_count"`
	Registrar       string          `json:"registrar"`
	Status          string          `json:"status"`
	TTL             int64           `json:"ttl"`
	TXTVerification TXTVerification `json:"txt_verification"`
	Verified        string          `json:"verified"`
}

type TXTVerification struct {
	Name  *string `json:"name"`
	Token *string `json:"token"`
}

type validateZoneResult struct {
	ParsedRecords int                `json:"parsed_records"`
	Error         *validateZoneError `json:"error"`
}

type validateZoneError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type zoneList struct {
	Zones []*Zone `json:"zones"`
	Meta  *meta   `json:"meta"`
}

type zone struct {
	Zone *Zone `json:"zone"`
}

type meta struct {
	Pagination *pagination `json:"pagination"`
}

type pagination struct {
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
	LastPage     int `json:"last_page"`
	TotalEntries int `json:"total_entries"`
}

type ZoneService interface {
	GetAllZones() ([]*Zone, error)
	GetAllZonesByName(name string) ([]*Zone, error)
	GetZoneById(zone_id string) (*Zone, error)
	CreateZone(request Zone) (*Zone, error)
	UpdateZone(zone_id string, request Zone) (*Zone, error)
	DeleteZone(zone_id string) (bool, error)
	ValidateZoneFile(zonefile string) error
	ExportZoneFile(zone_id string) (*string, error)
	ImportZoneFile(zone_id string, zonefile string) (*Zone, error)
}

type zoneService struct {
	client *client
}

func (service *zoneService) GetAllZones() ([]*Zone, error) {
	return service.GetAllZonesByName("")
}

func (service *zoneService) GetAllZonesByName(name string) ([]*Zone, error) {
	var zones []*Zone
	page := 1
	last_page := 1
	per_page := 100

	for page <= last_page {
		zoneList := new(zoneList)
		_, err := service.client.
			createJsonRequest(200).
			setQueryParams(
				map[string]string{
					"page":        fmt.Sprint(page),
					"per_page":    fmt.Sprint(per_page),
					"search_name": name,
				}).
			setResult(zoneList).
			execute("GET", zonesBasePath)
		if err != nil {
			return nil, err
		}
		page = page + 1
		last_page = zoneList.Meta.Pagination.LastPage
		zones = append(zones, zoneList.Zones...)
	}

	return zones, nil
}

func (service *zoneService) GetZoneById(zone_id string) (*Zone, error) {
	if err := validateNotEmpty("zone_id", zone_id); err != nil {
		return nil, err
	}

	zone := new(zone)
	_, err := service.client.
		createJsonRequest(200).
		setResult(zone).
		execute("GET", zonesBasePath+"/"+zone_id)
	if err != nil {
		return nil, err
	}
	return zone.Zone, nil
}

func (service *zoneService) CreateZone(request Zone) (*Zone, error) {
	zone := new(zone)
	_, err := service.client.
		createJsonRequest(201).
		setResult(zone).
		setBody(request).
		execute("POST", zonesBasePath)
	if err != nil {
		return nil, err
	}
	return zone.Zone, nil
}

func (service *zoneService) UpdateZone(zone_id string, request Zone) (*Zone, error) {
	if err := validateNotEmpty("zone_id", zone_id); err != nil {
		return nil, err
	}

	zone := new(zone)
	_, err := service.client.
		createJsonRequest(200).
		setResult(zone).
		setBody(request).
		execute("PUT", zonesBasePath+"/"+zone_id)
	if err != nil {
		return nil, err
	}
	return zone.Zone, nil
}

func (service *zoneService) DeleteZone(zone_id string) (bool, error) {
	if err := validateNotEmpty("zone_id", zone_id); err != nil {
		return false, err
	}

	zone := new(zone)
	_, err := service.client.
		createJsonRequest(200, 404).
		setResult(zone).
		execute("DELETE", zonesBasePath+"/"+zone_id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (service *zoneService) ValidateZoneFile(zonefile string) error {
	if err := validateNotEmpty("zonefile", zonefile); err != nil {
		return err
	}

	r := new(validateZoneResult)
	_, err := service.client.
		createTextRequest(200).
		setBody(zonefile).
		setResult(r).
		execute("POST", zonesBasePath+"/file/validate")

	if err != nil {
		return fmt.Errorf(r.Error.Message)
	}
	return err
}

func (service *zoneService) ExportZoneFile(zone_id string) (*string, error) {
	if err := validateNotEmpty("zone_id", zone_id); err != nil {
		return nil, err
	}

	zone, err := service.client.
		createTextRequest(200).
		execute("GET", zonesBasePath+"/"+zone_id+"/export")
	if err != nil {
		return nil, err
	}
	data := string(zone)
	return &data, err
}

func (service *zoneService) ImportZoneFile(zone_id string, zonefile string) (*Zone, error) {
	if err := validateNotEmpty("zone_id", zone_id); err != nil {
		return nil, err
	}
	if err := validateNotEmpty("zonefile", zonefile); err != nil {
		return nil, err
	}

	zone := new(zone)
	_, err := service.client.
		createTextRequest(200).
		setResult(zone).
		setBody(zonefile).
		execute("POST", zonesBasePath+"/"+zone_id+"/import")

	if err != nil {
		return nil, err
	}
	return zone.Zone, err
}
