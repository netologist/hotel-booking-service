package e2e_test

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/netologist/hotel-booking/pkg/client"
)

var _ = Describe("Hotels API (/api/v1/hotels)", func() {
	Describe("POST /api/v1/hotels (Register Hotel)", func() {
		It("creates a new hotel and returns HTTP 201 with generated fields", func(ctx context.Context) {
			hotelName := "E2E Hotel " + uuid.NewString()[:8]
			resp, err := apiClient.CreateHotelWithResponse(ctx, client.CreateHotelJSONRequestBody{
				Name:    hotelName,
				Address: "100 Grand Boulevard",
				City:    "Metropolis",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
			Expect(resp.ContentType()).To(ContainSubstring("application/json"))
			Expect(resp.JSON201).NotTo(BeNil())

			created := resp.JSON201
			Expect(created.Id).NotTo(Equal(uuid.Nil))
			Expect(created.Name).To(Equal(hotelName))
			Expect(created.Address).To(Equal("100 Grand Boulevard"))
			Expect(created.City).To(Equal("Metropolis"))
			Expect(created.CreatedAt).NotTo(BeZero())
			Expect(created.UpdatedAt).NotTo(BeZero())
			Expect(created.Rooms).NotTo(BeNil())
			Expect(*created.Rooms).To(HaveLen(6))
		})
		It("returns HTTP 400 when body has validation errors", func(ctx context.Context) {
			resp, err := apiClient.CreateHotelWithResponse(ctx, client.CreateHotelJSONRequestBody{
				Name:    "",
				Address: "",
				City:    "",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.StatusCode()).To(Equal(http.StatusBadRequest))
			Expect(resp.ContentType()).To(ContainSubstring("application/problem+json"))
			Expect(resp.ApplicationproblemJSON400).NotTo(BeNil())

			prob := resp.ApplicationproblemJSON400
			Expect(prob.Status).To(Equal(http.StatusBadRequest))
			Expect(prob.Title).To(Equal("Bad Request"))
			Expect(prob.Detail).NotTo(BeNil())
		})

		It("returns HTTP 400 when body is malformed JSON", func(ctx context.Context) {
			resp, err := apiClient.CreateHotelWithBodyWithResponse(
				ctx,
				"application/json",
				strings.NewReader(`{invalid`),
			)
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.StatusCode()).To(Equal(http.StatusBadRequest))
			Expect(resp.ContentType()).To(ContainSubstring("application/problem+json"))
			Expect(resp.ApplicationproblemJSON400).NotTo(BeNil())

			prob := resp.ApplicationproblemJSON400
			Expect(prob.Status).To(Equal(http.StatusBadRequest))
			Expect(prob.Title).To(Equal("Bad Request"))
			Expect(prob.Detail).NotTo(BeNil())
		})

		It("returns HTTP 409 when attempting to register a duplicate hotel name", func(ctx context.Context) {
			hotelName := "E2E Unique Hotel " + uuid.NewString()[:8]
			reqBody := client.CreateHotelJSONRequestBody{
				Name:    hotelName,
				Address: "200 Duplicate Ave",
				City:    "Gotham",
			}

			resp1, err := apiClient.CreateHotelWithResponse(ctx, reqBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp1.StatusCode()).To(Equal(http.StatusCreated))

			// Second request with same name
			resp2, err := apiClient.CreateHotelWithResponse(ctx, reqBody)
			Expect(err).NotTo(HaveOccurred())

			Expect(resp2.StatusCode()).To(Equal(http.StatusConflict))
			Expect(resp2.ContentType()).To(ContainSubstring("application/problem+json"))
			Expect(resp2.ApplicationproblemJSON409).NotTo(BeNil())

			prob := resp2.ApplicationproblemJSON409
			Expect(prob.Status).To(Equal(http.StatusConflict))
			Expect(prob.Title).To(Equal("Conflict"))
			Expect(prob.Code).NotTo(BeNil())
			Expect(*prob.Code).To(Equal("HOTEL_ALREADY_EXISTS"))
		})
	})

	Describe("GET /api/v1/hotels (Find Hotels)", func() {
		It("finds hotels matching the search query case-insensitively", func(ctx context.Context) {
			uniqueSuffix := uuid.NewString()[:8]
			hotelName := "Tracer Hotel " + uniqueSuffix
			createResp, err := apiClient.CreateHotelWithResponse(ctx, client.CreateHotelJSONRequestBody{
				Name:    hotelName,
				Address: "200 Search Ave",
				City:    "Query City",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(createResp.StatusCode()).To(Equal(http.StatusCreated))

			// Tracer query with partial and lower-case search
			searchParam := strings.ToLower(uniqueSuffix)
			searchResp, err := apiClient.FindHotelsWithResponse(ctx, &client.FindHotelsParams{
				Name: searchParam,
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(searchResp.StatusCode()).To(Equal(http.StatusOK))
			Expect(searchResp.ContentType()).To(ContainSubstring("application/json"))
			Expect(searchResp.JSON200).NotTo(BeNil())
			Expect(searchResp.JSON200.Hotels).NotTo(BeNil())
			hotels := *searchResp.JSON200.Hotels
			Expect(hotels).NotTo(BeEmpty())

			var found bool
			for _, h := range hotels {
				if h.Name == hotelName {
					found = true
					Expect(h.Address).To(Equal("200 Search Ave"))
					Expect(h.City).To(Equal("Query City"))
					break
				}
			}
			Expect(found).To(BeTrue(), "expected created hotel to be found in search results")
		})

		It("returns HTTP 200 with empty list when no hotels match", func(ctx context.Context) {
			nonExistent := "NonExistentHotel_" + uuid.NewString()
			searchResp, err := apiClient.FindHotelsWithResponse(ctx, &client.FindHotelsParams{
				Name: nonExistent,
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(searchResp.StatusCode()).To(Equal(http.StatusOK))
			Expect(searchResp.ContentType()).To(ContainSubstring("application/json"))
			Expect(searchResp.JSON200).NotTo(BeNil())
			Expect(searchResp.JSON200.Hotels).NotTo(BeNil())
			Expect(*searchResp.JSON200.Hotels).To(BeEmpty())
		})

		It("returns HTTP 400 when name query parameter is blank", func(ctx context.Context) {
			searchResp, err := apiClient.FindHotelsWithResponse(ctx, &client.FindHotelsParams{
				Name: "",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(searchResp.StatusCode()).To(Equal(http.StatusBadRequest))
			Expect(searchResp.ContentType()).To(ContainSubstring("application/problem+json"))
			Expect(searchResp.ApplicationproblemJSON400).NotTo(BeNil())

			prob := searchResp.ApplicationproblemJSON400
			Expect(prob.Status).To(Equal(http.StatusBadRequest))
			Expect(prob.Code).NotTo(BeNil())
			Expect(*prob.Code).To(Equal("INVALID_QUERY_PARAMETER"))
		})
	})

})
