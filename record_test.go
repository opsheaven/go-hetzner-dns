package gohetznerdns

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/assert"
)

func TestGetAllRecordsWithEmptyZoneId(t *testing.T) {
	zone_id := "            "
	service := &recordService{}
	_, err := service.GetAllRecords(&zone_id)
	assert.Error(t, err, "901 : zone_id is empty")
}

func TestGetAllRecords(t *testing.T) {
	zone_id := "zone_id"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		response := `
		{
			"records":[
				{
				"type":"A",
				"id":"123",
				"zone_id":"Zone_id",
				"name":"Record",
				"value":"Value",
				"ttl":3600
				}
			]
		}
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	records, err := recordService.GetAllRecords(&zone_id)

	assert.NilError(t, err)
	assert.Equal(t, len(records), 1)
	assert.Equal(t, *records[0].Id, "123")
	assert.Equal(t, *records[0].Name, "Record")
	assert.Equal(t, *records[0].TTL, 3600)
	assert.Equal(t, *records[0].Type, "A")
	assert.Equal(t, *records[0].Value, "Value")
	assert.Equal(t, *records[0].ZoneId, "Zone_id")

}

func TestGetAllRecordsError(t *testing.T) {
	zone_id := "zone_id"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		response := `
		{
			"records":[
				{
				"type":"A",
				"id":"123",
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	_, err := recordService.GetAllRecords(&zone_id)

	assert.Error(t, err, "unexpected end of JSON input")

}

func TestGetRecordWithEmptyRecordId(t *testing.T) {
	record_id := "         "
	service := &recordService{}
	_, err := service.GetRecord(&record_id)
	assert.Error(t, err, "901 : record_id is empty")
}

func TestGetRecord(t *testing.T) {
	record_id := "123"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records/123", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		response := `
		{
			"record":{
				"type":"A",
				"id":"123",
				"zone_id":"Zone_id",
				"name":"Record",
				"value":"Value",
				"ttl":3600
			}
		}
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	record, err := recordService.GetRecord(&record_id)

	assert.NilError(t, err)
	assert.Assert(t, record != nil)
	assert.Equal(t, *record.Id, "123")
	assert.Equal(t, *record.Name, "Record")
	assert.Equal(t, *record.TTL, 3600)
	assert.Equal(t, *record.Type, "A")
	assert.Equal(t, *record.Value, "Value")
	assert.Equal(t, *record.ZoneId, "Zone_id")
}

func TestGetRecordError(t *testing.T) {
	record_id := "123"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records/123", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		response := `
		{
			"records":[
				{
				"type":"A",
				"id":"123",
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	_, err := recordService.GetRecord(&record_id)

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestCreateRecord(t *testing.T) {
	recordId := "123"
	zoneType := "A"
	zoneId := "domain"
	zoneName := "www"
	zoneValue := "192.168.1.1"
	ttl := 3600
	record := &Record{
		Type:   &zoneType,
		ZoneId: &zoneId,
		Name:   &zoneName,
		Value:  &zoneValue,
		TTL:    &ttl,
	}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		expected, _ := json.Marshal(record)
		assert.Equal(t, string(body), string(expected))
		result := &Record{
			Id:     &recordId,
			Type:   &zoneType,
			ZoneId: &zoneId,
			Name:   &zoneName,
			Value:  &zoneValue,
			TTL:    &ttl,
		}
		recordResponse := &RecordResponse{Record: result}
		resultBody, _ := json.Marshal(recordResponse)
		fmt.Fprint(w, string(resultBody))
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	actual, err := recordService.CreateRecord(record)

	assert.NilError(t, err)
	assert.Assert(t, record != nil)
	assert.Equal(t, *actual.Id, recordId)
	assert.Equal(t, *actual.Name, zoneName)
	assert.Equal(t, *actual.TTL, ttl)
	assert.Equal(t, *actual.Type, zoneType)
	assert.Equal(t, *actual.Value, zoneValue)
	assert.Equal(t, *actual.ZoneId, zoneId)
}

func TestCreateRecordError(t *testing.T) {
	zoneType := "A"
	zoneId := "domain"
	zoneName := "www"
	zoneValue := "192.168.1.1"
	ttl := 3600
	record := &Record{
		Type:   &zoneType,
		ZoneId: &zoneId,
		Name:   &zoneName,
		Value:  &zoneValue,
		TTL:    &ttl,
	}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records", func(w http.ResponseWriter, r *http.Request) {
		response := `
		{
			"record":{
				"type":"A",
				"id":"123",
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	_, err := recordService.CreateRecord(record)

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestUpdateRecordWithEmptyRecordId(t *testing.T) {
	zoneType := "A"
	zoneId := "domain"
	zoneName := "www"
	zoneValue := "192.168.1.1"
	ttl := 3600
	record := &Record{
		Type:   &zoneType,
		ZoneId: &zoneId,
		Name:   &zoneName,
		Value:  &zoneValue,
		TTL:    &ttl,
	}
	service := &recordService{}
	_, err := service.UpdateRecord(record)
	assert.Error(t, err, "900 : record_id is nil")
}

func TestUpdateRecord(t *testing.T) {
	recordId := "123"
	recordType := "A"
	zoneId := "domain"
	recordName := "www"
	recordValue := "192.168.1.1"
	ttl := 3600
	record := &Record{
		Id:     &recordId,
		Type:   &recordType,
		ZoneId: &zoneId,
		Name:   &recordName,
		Value:  &recordValue,
		TTL:    &ttl,
	}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records/"+recordId, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		expected, _ := json.Marshal(record)
		assert.Equal(t, string(body), string(expected))
		recordResponse := &RecordResponse{Record: record}
		result, _ := json.Marshal(recordResponse)
		fmt.Fprint(w, string(result))
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	actual, err := recordService.UpdateRecord(record)

	assert.NilError(t, err)
	assert.Assert(t, record != nil)
	assert.Equal(t, *record.Id, *actual.Id)
	assert.Equal(t, *record.Name, *actual.Name)
	assert.Equal(t, *record.TTL, *actual.TTL)
	assert.Equal(t, *record.Type, *actual.Type)
	assert.Equal(t, *record.Value, *actual.Value)
	assert.Equal(t, *record.ZoneId, *actual.ZoneId)
}

func TestUpdateRecordWithError(t *testing.T) {
	recordId := "123"
	recordType := "A"
	zoneId := "domain"
	recordName := "www"
	recordValue := "192.168.1.1"
	ttl := 3600
	record := &Record{
		Id:     &recordId,
		Type:   &recordType,
		ZoneId: &zoneId,
		Name:   &recordName,
		Value:  &recordValue,
		TTL:    &ttl,
	}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records/"+recordId, func(w http.ResponseWriter, r *http.Request) {
		response := `
		{
			"record":{
				"type":"A",
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	record, err := recordService.UpdateRecord(record)

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestDeleteRecordWithEmptyRecordId(t *testing.T) {
	recordId := "         "
	service := &recordService{}
	err := service.DeleteRecord(&recordId)
	assert.Error(t, err, "901 : record_id is empty")
}

func TestDeleteRecord(t *testing.T) {
	recordId := "123"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records/"+recordId, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	err := recordService.DeleteRecord(&recordId)

	assert.NilError(t, err)
}
