package gohetznerdns

import (
	"fmt"
)

type ZoneService interface {
	GetAllZones() ([]*Zone, error)
	GetAllZonesByName(name *string) ([]*Zone, error)
	GetZoneById(zoneId *string) (*Zone, error)
	CreateZone(request *ZoneRequest) (*Zone, error)
	UpdateZone(zoneId *string, request *ZoneRequest) (*Zone, error)
	DeleteZone(zoneId *string) error
	ValidateZoneFile(zoneFile *string) error
	ExportZoneFile(zoneId *string) (*string, error)
	ImportZoneFile(zoneId, zoneFile *string) (*Zone, error)
}

type zoneService struct {
	client *client
}

var _ ZoneService = &zoneService{}

func (service *zoneService) GetAllZones() ([]*Zone, error) {
	return service.GetAllZonesByName(nil)
}

func (service *zoneService) GetAllZonesByName(name *string) ([]*Zone, error) {
	var zones []*Zone
	page := 1
	last_page := 1
	per_page := 100

	for page <= last_page {
		params := map[string]string{
			"page":     fmt.Sprint(page),
			"per_page": fmt.Sprint(per_page),
		}
		if name != nil {
			params["search_name"] = *name
		}
		zoneList := new(ZoneList)
		_, err := service.client.
			createJsonRequest(200).
			setQueryParams(params).
			setResult(zoneList).
			execute("GET", zonesBasePath)
		if err != nil {
			return nil, err
		}
		page = page + 1
		last_page = *zoneList.Meta.Pagination.LastPage
		zones = append(zones, zoneList.Zones...)
	}

	return zones, nil
}

func (service *zoneService) GetZoneById(zoneId *string) (*Zone, error) {
	if err := validateNotEmpty("zoneId", zoneId); err != nil {
		return nil, err
	}

	zone := new(ZoneResponse)
	_, err := service.client.
		createJsonRequest(200).
		setResult(zone).
		execute("GET", zonesBasePath+"/"+*zoneId)
	if err != nil {
		return nil, err
	}
	if zone.Error != nil {
		return nil, zone.Error.Error()
	}
	return zone.Zone, nil
}

func (service *zoneService) CreateZone(request *ZoneRequest) (*Zone, error) {
	zone := new(ZoneResponse)
	_, err := service.client.
		createJsonRequest(201).
		setResult(zone).
		setBody(request).
		execute("POST", zonesBasePath)
	if err != nil {
		return nil, err
	}
	if zone.Error != nil {
		return nil, zone.Error.Error()
	}
	return zone.Zone, nil
}

func (service *zoneService) UpdateZone(zoneId *string, request *ZoneRequest) (*Zone, error) {
	if err := validateNotEmpty("zoneId", zoneId); err != nil {
		return nil, err
	}
	zone := new(ZoneResponse)
	_, err := service.client.
		createJsonRequest(200).
		setResult(zone).
		setBody(request).
		execute("PUT", zonesBasePath+"/"+*zoneId)
	if err != nil {
		return nil, err
	}
	if zone.Error != nil {
		return nil, zone.Error.Error()
	}
	return zone.Zone, nil
}

func (service *zoneService) DeleteZone(zoneId *string) error {
	if err := validateNotEmpty("zoneId", zoneId); err != nil {
		return err
	}

	_, err := service.client.
		createJsonRequest(200, 404).
		execute("DELETE", zonesBasePath+"/"+*zoneId)
	return err
}

func (service *zoneService) ValidateZoneFile(zoneFile *string) error {
	if err := validateNotEmpty("zoneFile", zoneFile); err != nil {
		return err
	}

	zone := &ZoneResponse{}
	_, err := service.client.
		createTextRequest(200).
		setBody(*zoneFile).
		setResult(zone).
		execute("POST", zonesBasePath+"/file/validate")

	if err != nil {
		return err
	}

	if zone.Error != nil {
		return zone.Error.Error()
	}
	return err
}

func (service *zoneService) ExportZoneFile(zoneId *string) (*string, error) {
	if err := validateNotEmpty("zoneId", zoneId); err != nil {
		return nil, err
	}

	zone, err := service.client.
		createTextRequest(200).
		execute("GET", zonesBasePath+"/"+*zoneId+"/export")
	if err != nil {
		return nil, err
	}
	data := string(zone)
	return &data, err
}

func (service *zoneService) ImportZoneFile(zoneId, zoneFile *string) (*Zone, error) {
	if err := validateNotEmpty("zoneId", zoneId); err != nil {
		return nil, err
	}
	if err := validateNotEmpty("zoneFile", zoneFile); err != nil {
		return nil, err
	}

	zone := new(ZoneResponse)
	_, err := service.client.
		createTextRequest(200).
		setResult(zone).
		setBody(*zoneFile).
		execute("POST", zonesBasePath+"/"+*zoneId+"/import")

	if err != nil {
		return nil, err
	}
	if zone.Error != nil {
		return nil, zone.Error.Error()
	}
	return zone.Zone, err
}
