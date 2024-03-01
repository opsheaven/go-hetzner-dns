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

func TestGetAllZones(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		page := r.URL.Query().Get("page")
		if page == "1" {
			response := `
			{
				"zones":[
					{
					"id":"1",
					"name":"a"
					}
				],
				"meta":{
					"pagination": {
						"page":1,
						"per_page":1,
						"last_page":2,
						"total_entries":2
					}
				}
			}
			`
			fmt.Fprint(w, response)
		} else if page == "2" {
			response := `
			{
				"zones":[
					{
					"id":"2",
					"name":"b"
					}
				],
				"meta":{
					"pagination": {
						"page":2,
						"per_page":1,
						"last_page":2,
						"total_entries":2
					}
				}
			}
			`
			fmt.Fprint(w, response)
		}

	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	zones, err := zoneService.GetAllZones()

	assert.NilError(t, err)
	assert.Equal(t, len(zones), 2)
	assert.Equal(t, *zones[0].Id, "1")
	assert.Equal(t, *zones[0].Name, "a")
	assert.Equal(t, *zones[1].Id, "2")
	assert.Equal(t, *zones[1].Name, "b")
}

func TestGetAllZonesError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		response := `
			{
				"zones":[
					{
					"id":"1",
					"name":"a"
			`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	_, err := zoneService.GetAllZones()

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestGetZoneWithEmptyZoneId(t *testing.T) {
	token := "            "
	service := &zoneService{}
	_, err := service.GetZoneById(&token)
	assert.Error(t, err, "901 : zoneId is empty")
}

func TestGetZone(t *testing.T) {
	id := "1"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones/1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		response := `
		{
			"zone":{
				"id":"1",
				"name":"test.com"
			}
		}
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	zone, err := zoneService.GetZoneById(&id)

	assert.NilError(t, err)
	assert.Assert(t, zone != nil)
	assert.Equal(t, *zone.Id, "1")
	assert.Equal(t, *zone.Name, "test.com")
}

func TestGetZoneError(t *testing.T) {
	id := "1"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones/1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		response := `
		{
			"zone":{
				"name":"test.com",
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	_, err := zoneService.GetZoneById(&id)

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestCreateZone(t *testing.T) {
	name := "test.com"
	ttl := 3600
	zoneRequest := &ZoneRequest{Name: &name, TTL: &ttl}

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		expected, _ := json.Marshal(zoneRequest)
		assert.Equal(t, string(body), string(expected))
		response := `
		{
			"zone":{
				"name":"test.com",
				"id":"domain"
			}
		}
		`
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	zone, err := zoneService.CreateZone(zoneRequest)

	assert.NilError(t, err)
	assert.Assert(t, zone != nil)
	assert.Equal(t, *zone.Id, "domain")
	assert.Equal(t, *zone.Name, "test.com")
}

func TestCreateZoneError(t *testing.T) {
	name := "test.com"
	ttl := 3600
	zoneRequest := &ZoneRequest{Name: &name, TTL: &ttl}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones", func(w http.ResponseWriter, r *http.Request) {
		response := `
		{
			"zone":{
				"id":"domain",
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	_, err := zoneService.CreateZone(zoneRequest)

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestUpdateZoneWithEmptyZoneId(t *testing.T) {

	service := &zoneService{}
	_, err := service.UpdateZone(nil, &ZoneRequest{})
	assert.Error(t, err, "900 : zoneId is nil")
}

func TestUpdateZone(t *testing.T) {
	id := "domain"
	name := "test.com"
	ttl := 3600
	zoneRequest := &ZoneRequest{Name: &name, TTL: &ttl}

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones/domain", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		expected, _ := json.Marshal(zoneRequest)
		assert.Equal(t, string(body), string(expected))
		response := `
		{
			"zone":{
				"name":"test.com",
				"id":"domain"
			}
		}
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	zone, err := zoneService.UpdateZone(&id, zoneRequest)

	assert.NilError(t, err)
	assert.Assert(t, zone != nil)
	assert.Equal(t, *zone.Id, "domain")
	assert.Equal(t, *zone.Name, "test.com")
}

func TestUpdateZoneError(t *testing.T) {
	id := "domain"
	name := "test.com"
	ttl := 3600
	zoneRequest := &ZoneRequest{Name: &name, TTL: &ttl}
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones/domain", func(w http.ResponseWriter, r *http.Request) {
		response := `
		{
			"zone":{
				"id":"domain",
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	_, err := zoneService.UpdateZone(&id, zoneRequest)

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestDeleteZoneWithEmptyZoneId(t *testing.T) {
	id := "           "
	service := &zoneService{}
	err := service.DeleteZone(&id)
	assert.Error(t, err, "901 : zoneId is empty")
}

func TestDeleteZone(t *testing.T) {
	id := "domain"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones/domain", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	err := zoneService.DeleteZone(&id)

	assert.NilError(t, err)
}

func TestValidateZoneFileWithEmptyFile(t *testing.T) {
	zoneFile := "    "
	service := &zoneService{}
	err := service.ValidateZoneFile(&zoneFile)
	assert.Error(t, err, "901 : zoneFile is empty")
}

func TestValidateZone(t *testing.T) {
	zoneFile := "$ORIGIN opsheaven.space.\n$TTL 3600\n@		IN	SOA	hydrogen.ns.hetzner.com. dns.hetzner.com. 2024022431 86400 10800 3600000 3600"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones/file/validate", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		assert.Equal(t, string(body), string(zoneFile))
		response := `{
			"parsed_records": 0
		  }`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	err := zoneService.ValidateZoneFile(&zoneFile)

	assert.NilError(t, err)
}

func TestValidateZoneError(t *testing.T) {
	id := "domain"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones/file/validate", func(w http.ResponseWriter, r *http.Request) {
		response := `
		{
			"error":{
				"message":"Invalid Zone File"
			}
		}
		`
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	err := zoneService.ValidateZoneFile(&id)

	assert.Error(t, err, "422 Unprocessable Entity")
}

func TestExportZoneFileWithEmptyZoneId(t *testing.T) {
	id := "         "
	service := &zoneService{}
	_, err := service.ExportZoneFile(&id)
	assert.Error(t, err, "901 : zoneId is empty")
}

func TestExportZoneFile(t *testing.T) {
	id := "domain"
	expected := "$ORIGIN opsheaven.space."
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones/domain/export", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		fmt.Fprint(w, expected)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	actual, err := zoneService.ExportZoneFile(&id)

	assert.NilError(t, err)
	assert.Equal(t, *actual, expected)
}

func TestExportZoneFileError(t *testing.T) {
	id := "domain"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones/domain/export", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	_, err := zoneService.ExportZoneFile(&id)

	assert.Error(t, err, "404 Not Found")
}

func TestImportZoneFileWithEmptyZoneId(t *testing.T) {
	id := "  "
	zoneFile := "domain"
	service := &zoneService{}
	_, err := service.ImportZoneFile(&id, &zoneFile)
	assert.Error(t, err, "901 : zoneId is empty")
}

func TestImportZoneFileWithEmptyZoneFile(t *testing.T) {
	id := "domain"
	zoneFile := "   "
	service := &zoneService{}
	_, err := service.ImportZoneFile(&id, &zoneFile)
	assert.Error(t, err, "901 : zoneFile is empty")
}

func TestImportZoneFile(t *testing.T) {
	id := "domain"
	expected := "$ORIGIN opsheaven.space."
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones/domain/import", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		assert.Equal(t, string(body), expected)
		response := `
		{
			"zone":{
				"name":"test.com",
				"id":"domain"
			}
		}
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	zone, err := zoneService.ImportZoneFile(&id, &expected)

	assert.NilError(t, err)
	assert.Assert(t, zone != nil)
	assert.Equal(t, *zone.Id, "domain")
	assert.Equal(t, *zone.Name, "test.com")
}

func TestImporttZoneFileError(t *testing.T) {
	id := "domain"
	zoneFile := "asdadsa"
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones/domain/import", func(w http.ResponseWriter, r *http.Request) {
		response := `{"zone":{"name":"test.com",`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	_, err := zoneService.ImportZoneFile(&id, &zoneFile)

	assert.Error(t, err, "unexpected end of JSON input")
}
