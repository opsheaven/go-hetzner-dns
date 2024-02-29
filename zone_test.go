package hetznerdns

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
	assert.Equal(t, zones[0].Id, "1")
	assert.Equal(t, zones[0].Name, "a")
	assert.Equal(t, zones[1].Id, "2")
	assert.Equal(t, zones[1].Name, "b")
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
	service := &zoneService{}
	_, err := service.GetZoneById("    ")
	assert.Error(t, err, "zone_id is invalid because cannot be empty")
}

func TestGetZone(t *testing.T) {
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
	zone, err := zoneService.GetZoneById("1")

	assert.NilError(t, err)
	assert.Assert(t, zone != nil)
	assert.Equal(t, zone.Id, "1")
	assert.Equal(t, zone.Name, "test.com")
}

func TestGetZoneError(t *testing.T) {
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
	_, err := zoneService.GetZoneById("1")

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestCreateZone(t *testing.T) {
	zone := &Zone{
		Name:    "test.com",
		Owner:   "test",
		Project: "project",
	}

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
		expected, _ := json.Marshal(zone)
		assert.Equal(t, string(body), string(expected))
		response := `
		{
			"zone":{
				"name":"test.com",
				"owner":"test",
				"project":"project",
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
	zone, err := zoneService.CreateZone(*zone)

	assert.NilError(t, err)
	assert.Assert(t, zone != nil)
	assert.Equal(t, zone.Id, "domain")
	assert.Equal(t, zone.Name, "test.com")
	assert.Equal(t, zone.Project, "project")
	assert.Equal(t, zone.Owner, "test")
}

func TestCreateZoneError(t *testing.T) {
	zone := &Zone{
		Name:    "test.com",
		Owner:   "test",
		Project: "project",
	}
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
	_, err := zoneService.CreateZone(*zone)

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestUpdateZoneWithEmptyZoneId(t *testing.T) {
	service := &zoneService{}
	_, err := service.UpdateZone("    ", Zone{})
	assert.Error(t, err, "zone_id is invalid because cannot be empty")
}

func TestUpdateZone(t *testing.T) {
	zone := &Zone{
		Name:    "test.com",
		Owner:   "test",
		Project: "project",
		Id:      "domain",
	}

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
		expected, _ := json.Marshal(zone)
		assert.Equal(t, string(body), string(expected))
		response := `
		{
			"zone":{
				"name":"test.com",
				"owner":"test",
				"project":"project",
				"id":"domain"
			}
		}
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	zone, err := zoneService.UpdateZone("domain", *zone)

	assert.NilError(t, err)
	assert.Assert(t, zone != nil)
	assert.Equal(t, zone.Id, "domain")
	assert.Equal(t, zone.Name, "test.com")
	assert.Equal(t, zone.Project, "project")
	assert.Equal(t, zone.Owner, "test")
}

func TestUpdateZoneError(t *testing.T) {
	zone := &Zone{
		Name:    "test.com",
		Owner:   "test",
		Project: "project",
		Id:      "domain",
	}
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
	_, err := zoneService.UpdateZone("domain", *zone)

	assert.Error(t, err, "unexpected end of JSON input")
}

func TestDeleteZoneWithEmptyZoneId(t *testing.T) {
	service := &zoneService{}
	err := service.DeleteZone("    ")
	assert.Error(t, err, "zone_id is invalid because cannot be empty")
}

func TestDeleteZone(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones/domain", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	err := zoneService.DeleteZone("domain")

	assert.NilError(t, err)
}

func TestValidateZoneFileWithEmptyFile(t *testing.T) {
	service := &zoneService{}
	err := service.ValidateZoneFile("    ")
	assert.Error(t, err, "zonefile is invalid because cannot be empty")
}

func TestValidateZone(t *testing.T) {
	zonefile := "$ORIGIN pinchflat.dev.\n$TTL 3600\n@		IN	SOA	hydrogen.ns.hetzner.com. dns.hetzner.com. 2024022431 86400 10800 3600000 3600"
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
		assert.Equal(t, string(body), string(zonefile))
		response := `{
			"parsed_records": 0
		  }`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	err := zoneService.ValidateZoneFile(zonefile)

	assert.NilError(t, err)
}

func TestValidateZoneError(t *testing.T) {
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
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	err := zoneService.ValidateZoneFile("domain")

	assert.Error(t, err, "Invalid Zone File")
}

func TestExportZoneFileWithEmptyZoneId(t *testing.T) {
	service := &zoneService{}
	_, err := service.ExportZoneFile("    ")
	assert.Error(t, err, "zone_id is invalid because cannot be empty")
}

func TestExportZoneFile(t *testing.T) {
	expected := "$ORIGIN pinchflat.dev."
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
	actual, err := zoneService.ExportZoneFile("domain")

	assert.NilError(t, err)
	assert.Equal(t, *actual, expected)
}

func TestExportZoneFileError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	mux.HandleFunc("/api/v1/zones/domain/export", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	_, err := zoneService.ExportZoneFile("domain")

	assert.Error(t, err, "404 Not Found")
}

func TestImportZoneFileWithEmptyZoneId(t *testing.T) {
	service := &zoneService{}
	_, err := service.ImportZoneFile("    ", "asdadas")
	assert.Error(t, err, "zone_id is invalid because cannot be empty")
}

func TestImportZoneFileWithEmptyZoneFile(t *testing.T) {
	service := &zoneService{}
	_, err := service.ImportZoneFile("domain", "     ")
	assert.Error(t, err, "zonefile is invalid because cannot be empty")
}

func TestImportZoneFile(t *testing.T) {
	expected := "$ORIGIN pinchflat.dev."
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
				"owner":"test",
				"project":"project",
				"id":"domain"
			}
		}
		`
		fmt.Fprint(w, response)
	})

	client := newClient()
	client.setBaseURL(server.URL)
	zoneService := &zoneService{client: client}
	zone, err := zoneService.ImportZoneFile("domain", expected)

	assert.NilError(t, err)
	assert.Assert(t, zone != nil)
	assert.Equal(t, zone.Id, "domain")
	assert.Equal(t, zone.Name, "test.com")
	assert.Equal(t, zone.Project, "project")
	assert.Equal(t, zone.Owner, "test")
}

func TestImporttZoneFileError(t *testing.T) {
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
	_, err := zoneService.ImportZoneFile("domain", "asdada")

	assert.Error(t, err, "unexpected end of JSON input")
}
