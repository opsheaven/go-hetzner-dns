package gohetznerdns

// Client interfaces for the Hetzner DNS Public API Records endpoint
// See api documentation for more information [https://dns.hetzner.com/api-docs#tag/Records]
type RecordService interface {

	// Returns all records associated with user. [https://dns.hetzner.com/api-docs#operation/GetRecords]
	GetAllRecords(zone_id *string) ([]*Record, error)

	//Returns information about a single record. [https://dns.hetzner.com/api-docs#operation/GetRecord]
	GetRecord(record_id *string) (*Record, error)

	//Creates a new record. [https://dns.hetzner.com/api-docs#operation/CreateRecord]
	CreateRecord(request *Record) (*Record, error)

	//Updates a record. [https://dns.hetzner.com/api-docs#operation/UpdateRecord]
	UpdateRecord(request *Record) (*Record, error)

	//Deletes a record. [https://dns.hetzner.com/api-docs#operation/DeleteRecord]
	DeleteRecord(record_id *string) error
}

type recordService struct {
	client *client
}

func (service *recordService) GetAllRecords(zone_id *string) ([]*Record, error) {
	if err := validateNotEmpty("zone_id", zone_id); err != nil {
		return nil, err
	}
	records := new(Records)
	_, err := service.client.
		createJsonRequest(200).
		setQueryParams(
			map[string]string{
				"zone_id": *zone_id,
			}).
		setResult(records).
		execute("GET", recordsBasePath)
	if err != nil {
		return nil, err
	}
	return records.Records, nil
}
func (service *recordService) GetRecord(record_id *string) (*Record, error) {
	if err := validateNotEmpty("record_id", record_id); err != nil {
		return nil, err
	}
	record := new(RecordResponse)
	_, err := service.client.
		createJsonRequest(200).
		setResult(record).
		execute("GET", recordsBasePath+"/"+*record_id)
	if err != nil {
		return nil, err
	}
	return record.Record, nil
}

func (service *recordService) CreateRecord(request *Record) (*Record, error) {
	record := new(RecordResponse)
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

func (service *recordService) UpdateRecord(request *Record) (*Record, error) {
	if err := validateNotEmpty("record_id", request.Id); err != nil {
		return nil, err
	}
	record := new(RecordResponse)
	_, err := service.client.
		createJsonRequest(200).
		setResult(record).
		setBody(request).
		execute("PUT", recordsBasePath+"/"+*request.Id)
	if err != nil {
		return nil, err
	}
	return record.Record, nil
}

func (service *recordService) DeleteRecord(record_id *string) error {
	if err := validateNotEmpty("record_id", record_id); err != nil {
		return err
	}
	_, err := service.client.
		createTextRequest(200, 404).
		execute("DELETE", recordsBasePath+"/"+*record_id)

	return err
}
