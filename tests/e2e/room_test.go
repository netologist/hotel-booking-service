package e2e_test

import (
	"context"
	"net/http"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/netologist/hotel-booking/pkg/client"
)

var _ = Describe("Rooms API (/api/v1/rooms/available)", func() {
	BeforeEach(func(ctx context.Context) {
		// Reset and seed database before tests
		resetResp, err := apiClient.ResetDatabaseWithResponse(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(resetResp.StatusCode()).To(Equal(http.StatusOK))

		seedResp, err := apiClient.SeedDatabaseWithResponse(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(seedResp.StatusCode()).To(Equal(http.StatusOK))
	})

	Context("Finding available rooms by capacity and date", func() {
		It("filters rooms based on capacity invariant", func(ctx context.Context) {
			checkIn := openapi_types.Date{Time: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)}
			checkOut := openapi_types.Date{Time: time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC)}

			// 1 guest: all 12 rooms should be available (2 hotels x 6 rooms)
			resp1, err := apiClient.FindAvailableRoomsWithResponse(ctx, &client.FindAvailableRoomsParams{
				CheckInDate:  checkIn,
				CheckOutDate: checkOut,
				Guests:       1,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp1.StatusCode()).To(Equal(http.StatusOK))
			Expect(resp1.JSON200).NotTo(BeNil())
			Expect(resp1.JSON200.Rooms).To(HaveLen(12))

			// 2 guests: only double (2) and deluxe (3) rooms should be available (2 hotels x 4 rooms = 8)
			resp2, err := apiClient.FindAvailableRoomsWithResponse(ctx, &client.FindAvailableRoomsParams{
				CheckInDate:  checkIn,
				CheckOutDate: checkOut,
				Guests:       2,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp2.StatusCode()).To(Equal(http.StatusOK))
			Expect(resp2.JSON200).NotTo(BeNil())
			Expect(resp2.JSON200.Rooms).To(HaveLen(8))
			for _, r := range resp2.JSON200.Rooms {
				Expect(r.Capacity).To(BeNumerically(">=", 2))
				Expect(string(r.RoomType)).NotTo(Equal("single"))
			}

			// 3 guests: only deluxe (3) rooms should be available (2 hotels x 2 rooms = 4)
			resp3, err := apiClient.FindAvailableRoomsWithResponse(ctx, &client.FindAvailableRoomsParams{
				CheckInDate:  checkIn,
				CheckOutDate: checkOut,
				Guests:       3,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp3.StatusCode()).To(Equal(http.StatusOK))
			Expect(resp3.JSON200).NotTo(BeNil())
			Expect(resp3.JSON200.Rooms).To(HaveLen(4))
			for _, r := range resp3.JSON200.Rooms {
				Expect(r.Capacity).To(BeNumerically(">=", 3))
				Expect(string(r.RoomType)).To(Equal("deluxe"))
			}

			// 4 guests: no room can accommodate 4 people (capacity exceeded)
			resp4, err := apiClient.FindAvailableRoomsWithResponse(ctx, &client.FindAvailableRoomsParams{
				CheckInDate:  checkIn,
				CheckOutDate: checkOut,
				Guests:       4,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp4.StatusCode()).To(Equal(http.StatusOK))
			Expect(resp4.JSON200).NotTo(BeNil())
			Expect(resp4.JSON200.Rooms).To(BeEmpty())
		})

		It("filters available rooms by hotel ID", func(ctx context.Context) {
			checkIn := openapi_types.Date{Time: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)}
			checkOut := openapi_types.Date{Time: time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)}

			// Find Continental hotel to get its ID
			hotelsResp, err := apiClient.FindHotelsWithResponse(ctx, &client.FindHotelsParams{Name: hotelContinental})
			Expect(err).NotTo(HaveOccurred())
			Expect(hotelsResp.JSON200.Hotels).NotTo(BeNil())
			Expect(*hotelsResp.JSON200.Hotels).To(HaveLen(1))
			continentalID := (*hotelsResp.JSON200.Hotels)[0].Id

			resp, err := apiClient.FindAvailableRoomsWithResponse(ctx, &client.FindAvailableRoomsParams{
				CheckInDate:  checkIn,
				CheckOutDate: checkOut,
				Guests:       1,
				HotelId:      &continentalID,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			Expect(resp.JSON200.Rooms).To(HaveLen(6))
			for _, r := range resp.JSON200.Rooms {
				Expect(r.HotelId).To(Equal(continentalID))
				Expect(r.HotelName).To(Equal("The Continental"))
			}
		})

		It("returns HTTP 400 when check-out date is not after check-in date", func(ctx context.Context) {
			checkIn := openapi_types.Date{Time: time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC)}
			checkOut := openapi_types.Date{Time: time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC)}

			resp, err := apiClient.FindAvailableRoomsWithResponse(ctx, &client.FindAvailableRoomsParams{
				CheckInDate:  checkIn,
				CheckOutDate: checkOut,
				Guests:       1,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusBadRequest))
			Expect(resp.ContentType()).To(ContainSubstring("application/problem+json"))
		})

		It("returns HTTP 400 when check-out date equals check-in date", func(ctx context.Context) {
			checkIn := openapi_types.Date{Time: time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC)}
			checkOut := openapi_types.Date{Time: time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC)}

			resp, err := apiClient.FindAvailableRoomsWithResponse(ctx, &client.FindAvailableRoomsParams{
				CheckInDate:  checkIn,
				CheckOutDate: checkOut,
				Guests:       1,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusBadRequest))
			Expect(resp.ContentType()).To(ContainSubstring("application/problem+json"))
		})
	})
})
