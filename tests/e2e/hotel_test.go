package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/netologist/hotel-booking/internal/gen/openapi"
)

const (
	fieldHotelName    = "name"
	fieldHotelAddress = "address"
	fieldHotelCity    = "city"
)

var _ = Describe("Hotels API (/api/v1/hotels)", func() {
	Describe("POST /api/v1/hotels (Register Hotel)", func() {
		It("creates a new hotel and returns HTTP 201 with generated fields", func() {
			hotelName := "E2E Hotel " + uuid.NewString()[:8]
			payload, err := json.Marshal(map[string]string{
				fieldHotelName:    hotelName,
				fieldHotelAddress: "100 Grand Boulevard",
				fieldHotelCity:    "Metropolis",
			})
			Expect(err).NotTo(HaveOccurred())

			resp, err := httpClient.Post(
				baseURL+"/api/v1/hotels",
				"application/json",
				bytes.NewReader(payload),
			)
			Expect(err).NotTo(HaveOccurred())
			defer func() { _ = resp.Body.Close() }()

			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			Expect(resp.Header.Get("Content-Type")).To(ContainSubstring("application/json"))

			var created openapi.Hotel
			err = json.NewDecoder(resp.Body).Decode(&created)
			Expect(err).NotTo(HaveOccurred())

			Expect(created.Id).NotTo(Equal(uuid.Nil))
			Expect(created.Name).To(Equal(hotelName))
			Expect(created.Address).To(Equal("100 Grand Boulevard"))
			Expect(created.City).To(Equal("Metropolis"))
			Expect(created.CreatedAt).NotTo(BeNil())
			Expect(created.UpdatedAt).NotTo(BeNil())
		})

		It("returns HTTP 400 when body has validation errors", func() {
			payload := `{"name":"","address":"","city":""}`
			resp, err := httpClient.Post(
				baseURL+"/api/v1/hotels",
				"application/json",
				strings.NewReader(payload),
			)
			Expect(err).NotTo(HaveOccurred())
			defer func() { _ = resp.Body.Close() }()

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			Expect(resp.Header.Get("Content-Type")).To(ContainSubstring("application/problem+json"))

			var prob openapi.Problem
			err = json.NewDecoder(resp.Body).Decode(&prob)
			Expect(err).NotTo(HaveOccurred())

			Expect(prob.Status).To(Equal(http.StatusBadRequest))
			Expect(prob.Title).To(Equal("Bad Request"))
			Expect(prob.Detail).NotTo(BeNil())
		})

		It("returns HTTP 400 when body is malformed JSON", func() {
			resp, err := httpClient.Post(
				baseURL+"/api/v1/hotels",
				"application/json",
				strings.NewReader(`{invalid`),
			)
			Expect(err).NotTo(HaveOccurred())
			defer func() { _ = resp.Body.Close() }()

			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			Expect(resp.Header.Get("Content-Type")).To(ContainSubstring("application/problem+json"))

			var prob openapi.Problem
			err = json.NewDecoder(resp.Body).Decode(&prob)
			Expect(err).NotTo(HaveOccurred())

			Expect(prob.Status).To(Equal(http.StatusBadRequest))
			Expect(prob.Title).To(Equal("Bad Request"))
			Expect(prob.Detail).NotTo(BeNil())
		})

		It("returns HTTP 409 when attempting to register a duplicate hotel name", func() {
			hotelName := "E2E Unique Hotel " + uuid.NewString()[:8]
			payload, err := json.Marshal(map[string]string{
				fieldHotelName:    hotelName,
				fieldHotelAddress: "200 Duplicate Ave",
				fieldHotelCity:    "Gotham",
			})
			Expect(err).NotTo(HaveOccurred())

			resp1, err := httpClient.Post(
				baseURL+"/api/v1/hotels",
				"application/json",
				bytes.NewReader(payload),
			)
			Expect(err).NotTo(HaveOccurred())
			defer func() { _ = resp1.Body.Close() }()
			Expect(resp1.StatusCode).To(Equal(http.StatusCreated))

			// Second request with same name
			resp2, err := httpClient.Post(
				baseURL+"/api/v1/hotels",
				"application/json",
				bytes.NewReader(payload),
			)
			Expect(err).NotTo(HaveOccurred())
			defer func() { _ = resp2.Body.Close() }()

			Expect(resp2.StatusCode).To(Equal(http.StatusConflict))
			Expect(resp2.Header.Get("Content-Type")).To(ContainSubstring("application/problem+json"))

			var prob openapi.Problem
			err = json.NewDecoder(resp2.Body).Decode(&prob)
			Expect(err).NotTo(HaveOccurred())

			Expect(prob.Status).To(Equal(http.StatusConflict))
			Expect(prob.Title).To(Equal("Conflict"))
			Expect(prob.Code).NotTo(BeNil())
			Expect(*prob.Code).To(Equal("HOTEL_ALREADY_EXISTS"))
		})
	})

	Describe("GET /api/v1/hotels (Find Hotels)", func() {
		It("finds hotels matching the search query case-insensitively", func() {
			uniqueSuffix := uuid.NewString()[:8]
			hotelName := "Tracer Hotel " + uniqueSuffix
			payload, err := json.Marshal(map[string]string{
				fieldHotelName:    hotelName,
				fieldHotelAddress: "200 Search Ave",
				fieldHotelCity:    "Query City",
			})
			Expect(err).NotTo(HaveOccurred())

			createResp, err := httpClient.Post(
				baseURL+"/api/v1/hotels",
				"application/json",
				bytes.NewReader(payload),
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(createResp.StatusCode).To(Equal(http.StatusCreated))
			_ = createResp.Body.Close()

			// Tracer query with partial and lower-case search
			searchResp, err := httpClient.Get(baseURL + "/api/v1/hotels?name=" + strings.ToLower(uniqueSuffix))
			Expect(err).NotTo(HaveOccurred())
			defer func() { _ = searchResp.Body.Close() }()

			Expect(searchResp.StatusCode).To(Equal(http.StatusOK))
			Expect(searchResp.Header.Get("Content-Type")).To(ContainSubstring("application/json"))

			var result struct {
				Hotels []openapi.Hotel `json:"hotels"`
			}
			err = json.NewDecoder(searchResp.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())

			Expect(result.Hotels).NotTo(BeEmpty())
			var found bool
			for _, h := range result.Hotels {
				if h.Name == hotelName {
					found = true
					Expect(h.Address).To(Equal("200 Search Ave"))
					Expect(h.City).To(Equal("Query City"))
					break
				}
			}
			Expect(found).To(BeTrue(), "expected created hotel to be found in search results")
		})

		It("returns HTTP 200 with empty list when no hotels match", func() {
			nonExistent := "NonExistentHotel_" + uuid.NewString()
			searchResp, err := httpClient.Get(baseURL + "/api/v1/hotels?name=" + nonExistent)
			Expect(err).NotTo(HaveOccurred())
			defer func() { _ = searchResp.Body.Close() }()

			Expect(searchResp.StatusCode).To(Equal(http.StatusOK))
			Expect(searchResp.Header.Get("Content-Type")).To(ContainSubstring("application/json"))

			var result struct {
				Hotels []openapi.Hotel `json:"hotels"`
			}
			err = json.NewDecoder(searchResp.Body).Decode(&result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Hotels).To(BeEmpty())
		})

		It("returns HTTP 400 when name query parameter is blank", func() {
			searchResp, err := httpClient.Get(baseURL + "/api/v1/hotels?name=")
			Expect(err).NotTo(HaveOccurred())
			defer func() { _ = searchResp.Body.Close() }()

			Expect(searchResp.StatusCode).To(Equal(http.StatusBadRequest))
			Expect(searchResp.Header.Get("Content-Type")).To(ContainSubstring("application/problem+json"))

			var prob openapi.Problem
			err = json.NewDecoder(searchResp.Body).Decode(&prob)
			Expect(err).NotTo(HaveOccurred())
			Expect(prob.Status).To(Equal(http.StatusBadRequest))
			Expect(prob.Code).NotTo(BeNil())
			Expect(*prob.Code).To(Equal("INVALID_QUERY_PARAMETER"))
		})
	})

})
