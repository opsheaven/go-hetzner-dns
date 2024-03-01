package gohetznerdns

type Record struct {
	Type   string `json:"type"`
	Id     string `json:"id"`
	ZoneId string `json:"zone_id"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	TTL    int    `json:"ttl"`
}
type record struct {
	Record *Record `json:"record"`
}
type Records struct {
	Records        []*Record `json:"records"`
	ValidRecords   []*Record `json:"valid_records"`
	InvalidRecords []*Record `json:"invalid_records"`
}

/*
type RecordService interface {
	GetAllRecords(zone_id string) ([]*Record, error)
	GetRecord(record_id string) (*Record, error)
	CreateRecord(request Record) (*Record, error)
	UpdateRecord(record_id string, request Record) (*Record, error)
	DeleteRecord(record_id string) error
	CreateRecords(request Records) (*Records, error)
	UpdateRecords(request Records) (*Records, error)
}

type recordService struct {
	client *client
}

// GetAllRecords implements RecordsService.
func (service *recordService) GetAllRecords(zone_id string) ([]*Record, error) {
	if err := validateNotEmpty("zone_id", zone_id); err != nil {
		return nil, err
	}
	records := new(Records)
	_, err := service.client.
		createJsonRequest(200).
		setQueryParams(
			map[string]string{
				"zone_id": zone_id,
			}).
		setResult(records).
		execute("GET", recordsBasePath)
	if err != nil {
		return nil, err
	}
	return records.Records, nil
}

// GetRecord implements RecordsService.
func (service *recordService) GetRecord(record_id string) (*Record, error) {
	if err := validateNotEmpty("record_id", record_id); err != nil {
		return nil, err
	}
	record := new(record)
	_, err := service.client.
		createJsonRequest(200).
		setResult(record).
		execute("GET", recordsBasePath+"/"+record_id)
	if err != nil {
		return nil, err
	}
	return record.Record, nil
}

// CreateRecord implements RecordsService.
func (service *recordService) CreateRecord(request Record) (*Record, error) {
	record := new(record)
	_, err := service.client.
		createJsonRequest(200).
		setResult(record).
		setBody(request).
		execute("POST", recordsBasePath)
	if err != nil {
		return nil, err
	}
	return record.Record, nil
}

// UpdateRecord implements RecordsService.
func (service *recordService) UpdateRecord(record_id string, request Record) (*Record, error) {
	if err := validateNotEmpty("record_id", record_id); err != nil {
		return nil, err
	}
	record := new(record)
	_, err := service.client.
		createJsonRequest(200).
		setResult(record).
		setBody(request).
		execute("PUT", recordsBasePath+"/"+record_id)
	if err != nil {
		return nil, err
	}
	return record.Record, nil
}

// GetRecord implements RecordsService.
func (service *recordService) DeleteRecord(record_id string) error {
	if err := validateNotEmpty("record_id", record_id); err != nil {
		return err
	}
	_, err := service.client.
		createTextRequest(200, 404).
		execute("DELETE", recordsBasePath+"/"+record_id)

	return err
}

// BulkCreateRecords implements RecordsService.
func (service *recordService) CreateRecords(request Records) (*Records, error) {
	response := new(Records)
	_, err := service.client.
		createJsonRequest(200).
		setResult(response).
		setBody(request).
		execute("POST", recordsBasePath+"/bulk")
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (service *recordService) UpdateRecords(request Records) (*Records, error) {
	response := new(Records)
	_, err := service.client.
		createJsonRequest(200).
		setResult(response).
		setBody(request).
		execute("PUT", recordsBasePath+"/bulk")
	if err != nil {
		return nil, err
	}
	return response, nil
}
*/
