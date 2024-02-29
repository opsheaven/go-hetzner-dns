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
	service := &recordService{}
	_, err := service.GetAllRecords("    ")
	assert.Error(t, err, "zone_id is invalid because cannot be empty")
}

func TestGetAllRecords(t *testing.T) {
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
	records, err := recordService.GetAllRecords("zone_id")

	assert.NilError(t, err)
	assert.Equal(t, len(records), 1)
	assert.Equal(t, records[0].Id, "123")
	assert.Equal(t, records[0].Name, "Record")
	assert.Equal(t, records[0].TTL, 3600)
	assert.Equal(t, records[0].Type, "A")
	assert.Equal(t, records[0].Value, "Value")
	assert.Equal(t, records[0].ZoneId, "Zone_id")

}

func TestGetAllRecordsError(t *testing.T) {
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
	_, err := recordService.GetAllRecords("zone_id")

	assert.Error(t, err, "unexpected end of JSON input")

}

func TestGetRecordWithEmptyRecordId(t *testing.T) {
	service := &recordService{}
	_, err := service.GetRecord("    ")
	assert.Error(t, err, "record_id is invalid because cannot be empty")
}

func TestGetRecord(t *testing.T) {
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
	record, err := recordService.GetRecord("123")

	assert.NilError(t, err)
	assert.Assert(t, record != nil)
	assert.Equal(t, record.Id, "123")
	assert.Equal(t, record.Name, "Record")
	assert.Equal(t, record.TTL, 3600)
	assert.Equal(t, record.Type, "A")
	assert.Equal(t, record.Value, "Value")
	assert.Equal(t, record.ZoneId, "Zone_id")
}

func TestGetRecordError(t *testing.T) {
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
	_, err := recordService.GetRecord("123")

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestCreateRecord(t *testing.T) {
	record := &Record{
		Type:   "A",
		ZoneId: "Zone_Id",
		Name:   "www",
		Value:  "192.168.1.1",
		TTL:    3600,
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
		response := `
		{
			"record":{
				"type":"A",
				"id":"123",
				"zone_id":"Zone_Id",
				"name":"www",
				"value":"192.168.1.1",
				"ttl":3600
			}
		}
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	record, err := recordService.CreateRecord(*record)

	assert.NilError(t, err)
	assert.Assert(t, record != nil)
	assert.Equal(t, record.Id, "123")
	assert.Equal(t, record.Name, "www")
	assert.Equal(t, record.TTL, 3600)
	assert.Equal(t, record.Type, "A")
	assert.Equal(t, record.Value, "192.168.1.1")
	assert.Equal(t, record.ZoneId, "Zone_Id")
}

func TestCreateRecordError(t *testing.T) {
	record := &Record{
		Type:   "A",
		ZoneId: "Zone_Id",
		Name:   "www",
		Value:  "192.168.1.1",
		TTL:    3600,
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
	_, err := recordService.CreateRecord(*record)

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestUpdateRecordWithEmptyRecordId(t *testing.T) {
	record := &Record{
		Type:   "A",
		ZoneId: "Zone_Id",
		Name:   "www",
		Value:  "192.168.1.1",
		TTL:    3600,
	}
	service := &recordService{}
	_, err := service.UpdateRecord("    ", *record)
	assert.Error(t, err, "record_id is invalid because cannot be empty")
}

func TestUpdateRecord(t *testing.T) {
	record := &Record{
		Type:   "A",
		ZoneId: "Zone_Id",
		Name:   "www",
		Value:  "192.168.1.1",
		TTL:    3600,
	}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records/123", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		expected, _ := json.Marshal(record)
		assert.Equal(t, string(body), string(expected))
		response := `
		{
			"record":{
				"type":"A",
				"id":"123",
				"zone_id":"Zone_Id",
				"name":"www",
				"value":"192.168.1.1",
				"ttl":3600
			}
		}
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	record, err := recordService.UpdateRecord("123", *record)

	assert.NilError(t, err)
	assert.Assert(t, record != nil)
	assert.Equal(t, record.Id, "123")
	assert.Equal(t, record.Name, "www")
	assert.Equal(t, record.TTL, 3600)
	assert.Equal(t, record.Type, "A")
	assert.Equal(t, record.Value, "192.168.1.1")
	assert.Equal(t, record.ZoneId, "Zone_Id")
}

func TestUpdateRecordWithError(t *testing.T) {
	record := &Record{
		Type:   "A",
		ZoneId: "Zone_Id",
		Name:   "www",
		Value:  "192.168.1.1",
		TTL:    3600,
	}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records/123", func(w http.ResponseWriter, r *http.Request) {
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
	record, err := recordService.UpdateRecord("123", *record)

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestDeleteRecordWithEmptyRecordId(t *testing.T) {
	service := &recordService{}
	err := service.DeleteRecord("    ")
	assert.Error(t, err, "record_id is invalid because cannot be empty")
}

func TestDeleteRecord(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records/123", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	err := recordService.DeleteRecord("123")

	assert.NilError(t, err)
}

func TestCreateRecords(t *testing.T) {
	records := &Records{
		Records: []*Record{
			{
				Type:   "A",
				ZoneId: "Zone_Id",
				Name:   "www",
				Value:  "192.168.1.1",
				TTL:    3600,
			},
			{
				Type:   "A",
				ZoneId: "Zone_Id",
				Name:   "www",
				Value:  "192.168.1.2",
				TTL:    3600,
			},
		},
	}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records/bulk", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		expected, _ := json.Marshal(records)
		assert.Equal(t, string(body), string(expected))
		response := `
		{
			"records":[
				{
					"type":"A",
					"id":"1",
					"zone_id":"Zone_Id",
					"name":"www",
					"value":"192.168.1.1",
					"ttl":3600
				},
				{
					"type":"A",
					"id":"2",
					"zone_id":"Zone_Id",
					"name":"www",
					"value":"192.168.1.2",
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
	records, err := recordService.CreateRecords(*records)

	assert.NilError(t, err)
	assert.Assert(t, records != nil)
	assert.Equal(t, len(records.Records), 2)
	assert.Equal(t, records.Records[0].Id, "1")
	assert.Equal(t, records.Records[0].Name, "www")
	assert.Equal(t, records.Records[0].TTL, 3600)
	assert.Equal(t, records.Records[0].Type, "A")
	assert.Equal(t, records.Records[0].Value, "192.168.1.1")
	assert.Equal(t, records.Records[0].ZoneId, "Zone_Id")
	assert.Equal(t, records.Records[1].Id, "2")
	assert.Equal(t, records.Records[1].Name, "www")
	assert.Equal(t, records.Records[1].TTL, 3600)
	assert.Equal(t, records.Records[1].Type, "A")
	assert.Equal(t, records.Records[1].Value, "192.168.1.2")
	assert.Equal(t, records.Records[1].ZoneId, "Zone_Id")
}

func TestCreateRecordsError(t *testing.T) {
	records := &Records{
		Records: []*Record{
			{
				Type:   "A",
				ZoneId: "Zone_Id",
				Name:   "www",
				Value:  "192.168.1.1",
				TTL:    3600,
			},
			{
				Type:   "A",
				ZoneId: "Zone_Id",
				Name:   "www",
				Value:  "192.168.1.2",
				TTL:    3600,
			},
		},
	}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records/bulk", func(w http.ResponseWriter, r *http.Request) {
		response := `
		{
			"records":[
				{
					"type":"A",
					"id":"1",
					"zone_id":"Zone_Id",
					"name":"www",
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	records, err := recordService.CreateRecords(*records)

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestUpdateRecords(t *testing.T) {
	records := &Records{
		Records: []*Record{
			{
				Type:   "A",
				ZoneId: "Zone_Id",
				Name:   "www",
				Value:  "192.168.1.1",
				TTL:    3600,
				Id:     "1",
			},
			{
				Type:   "A",
				ZoneId: "Zone_Id",
				Name:   "www",
				Value:  "192.168.1.2",
				TTL:    3600,
				Id:     "2",
			},
		},
	}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records/bulk", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		expected, _ := json.Marshal(records)
		assert.Equal(t, string(body), string(expected))
		fmt.Fprint(w, string(body))
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	records, err := recordService.UpdateRecords(*records)

	assert.NilError(t, err)
	assert.Assert(t, records != nil)
	assert.Equal(t, len(records.Records), 2)
	assert.Equal(t, records.Records[0].Id, "1")
	assert.Equal(t, records.Records[0].Name, "www")
	assert.Equal(t, records.Records[0].TTL, 3600)
	assert.Equal(t, records.Records[0].Type, "A")
	assert.Equal(t, records.Records[0].Value, "192.168.1.1")
	assert.Equal(t, records.Records[0].ZoneId, "Zone_Id")
	assert.Equal(t, records.Records[1].Id, "2")
	assert.Equal(t, records.Records[1].Name, "www")
	assert.Equal(t, records.Records[1].TTL, 3600)
	assert.Equal(t, records.Records[1].Type, "A")
	assert.Equal(t, records.Records[1].Value, "192.168.1.2")
	assert.Equal(t, records.Records[1].ZoneId, "Zone_Id")
}

func TestUpdateRecordsError(t *testing.T) {
	records := &Records{
		Records: []*Record{
			{
				Type:   "A",
				ZoneId: "Zone_Id",
				Name:   "www",
				Value:  "192.168.1.1",
				TTL:    3600,
				Id:     "1",
			},
			{
				Type:   "A",
				ZoneId: "Zone_Id",
				Name:   "www",
				Value:  "192.168.1.2",
				TTL:    3600,
				Id:     "2",
			},
		},
	}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/records/bulk", func(w http.ResponseWriter, r *http.Request) {
		response := `
		{
			"records":[
				{
					"type":"A",
					"id":"1",
					"zone_id":"Zone_Id",
					"name":"www",
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	recordService := &recordService{client: client}
	records, err := recordService.UpdateRecords(*records)

	assert.Error(t, err, "unexpected end of JSON input")
}
