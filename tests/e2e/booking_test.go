package e2e_test

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/netologist/hotel-booking/pkg/client"
)

var _ = Describe("Bookings API (/api/v1/bookings)", func() {
	BeforeEach(func(ctx context.Context) {
		resetResp, err := apiClient.ResetDatabaseWithResponse(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(resetResp.StatusCode()).To(Equal(http.StatusOK))

		seedResp, err := apiClient.SeedDatabaseWithResponse(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(seedResp.StatusCode()).To(Equal(http.StatusOK))
	})

	Context("POST /api/v1/bookings (Book a Room)", func() {
		It("books a room successfully and returns HTTP 201 with booking reference", func(ctx context.Context) {
			// Find available rooms for 2 guests
			checkIn := openapi_types.Date{Time: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)}
			checkOut := openapi_types.Date{Time: time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)}

			availResp, err := apiClient.FindAvailableRoomsWithResponse(ctx, &client.FindAvailableRoomsParams{
				CheckInDate:  checkIn,
				CheckOutDate: checkOut,
				Guests:       2,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(availResp.StatusCode()).To(Equal(http.StatusOK))
			Expect(availResp.JSON200.Rooms).NotTo(BeEmpty())

			targetRoom := availResp.JSON200.Rooms[0]

			// Book target room
			bookResp, err := apiClient.BookRoomWithResponse(ctx, client.BookRoomJSONRequestBody{
				RoomId:         targetRoom.Id,
				CheckInDate:    checkIn,
				CheckOutDate:   checkOut,
				NumberOfGuests: 2,
				GuestName:      "Alice Wonderland",
				GuestEmail:     openapi_types.Email("alice@example.com"),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(bookResp.StatusCode()).To(Equal(http.StatusCreated))
			Expect(bookResp.JSON201).NotTo(BeNil())

			booking := bookResp.JSON201
			Expect(booking.Id).NotTo(Equal(uuid.Nil))
			Expect(booking.Reference).To(HavePrefix("BK-"))
			Expect(booking.RoomId).To(Equal(targetRoom.Id))
			Expect(booking.HotelId).To(Equal(targetRoom.HotelId))
			Expect(booking.GuestName).To(Equal("Alice Wonderland"))
			Expect(booking.GuestEmail).To(Equal("alice@example.com"))
			Expect(booking.NumberOfGuests).To(Equal(2))

			// Room should no longer be available for overlapping dates
			availAfter, err := apiClient.FindAvailableRoomsWithResponse(ctx, &client.FindAvailableRoomsParams{
				CheckInDate:  checkIn,
				CheckOutDate: checkOut,
				Guests:       2,
				HotelId:      &targetRoom.HotelId,
			})
			Expect(err).NotTo(HaveOccurred())
			for _, r := range availAfter.JSON200.Rooms {
				Expect(r.Id).NotTo(Equal(targetRoom.Id), "booked room should not be in available rooms")
			}
		})

		It("prevents double-booking the same room for overlapping dates (HTTP 409)", func(ctx context.Context) {
			checkIn1 := openapi_types.Date{Time: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)}
			checkOut1 := openapi_types.Date{Time: time.Date(2026, 11, 4, 0, 0, 0, 0, time.UTC)}

			availResp, err := apiClient.FindAvailableRoomsWithResponse(ctx, &client.FindAvailableRoomsParams{
				CheckInDate:  checkIn1,
				CheckOutDate: checkOut1,
				Guests:       1,
			})
			Expect(err).NotTo(HaveOccurred())
			targetRoom := availResp.JSON200.Rooms[0]

			// First booking succeeds
			book1, err := apiClient.BookRoomWithResponse(ctx, client.BookRoomJSONRequestBody{
				RoomId:         targetRoom.Id,
				CheckInDate:    checkIn1,
				CheckOutDate:   checkOut1,
				NumberOfGuests: 1,
				GuestName:      "First Guest",
				GuestEmail:     openapi_types.Email("first@example.com"),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(book1.StatusCode()).To(Equal(http.StatusCreated))

			// Second booking overlaps (Nov 2 to Nov 5 overlaps Nov 1 to Nov 4 on Nov 2 & Nov 3 nights)
			checkIn2 := openapi_types.Date{Time: time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC)}
			checkOut2 := openapi_types.Date{Time: time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC)}

			book2, err := apiClient.BookRoomWithResponse(ctx, client.BookRoomJSONRequestBody{
				RoomId:         targetRoom.Id,
				CheckInDate:    checkIn2,
				CheckOutDate:   checkOut2,
				NumberOfGuests: 1,
				GuestName:      "Second Guest",
				GuestEmail:     openapi_types.Email("second@example.com"),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(book2.StatusCode()).To(Equal(http.StatusConflict))
			Expect(book2.ContentType()).To(ContainSubstring("application/problem+json"))
			Expect(book2.ApplicationproblemJSON409).NotTo(BeNil())
			Expect(*book2.ApplicationproblemJSON409.Code).To(Equal("ROOM_ALREADY_BOOKED"))
		})

		It("allows booking adjacent non-overlapping dates for the same room", func(ctx context.Context) {
			checkIn1 := openapi_types.Date{Time: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)}
			checkOut1 := openapi_types.Date{Time: time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)}

			availResp, err := apiClient.FindAvailableRoomsWithResponse(ctx, &client.FindAvailableRoomsParams{
				CheckInDate:  checkIn1,
				CheckOutDate: checkOut1,
				Guests:       1,
			})
			Expect(err).NotTo(HaveOccurred())
			targetRoom := availResp.JSON200.Rooms[0]

			// First booking: Nov 1 to Nov 3
			book1, err := apiClient.BookRoomWithResponse(ctx, client.BookRoomJSONRequestBody{
				RoomId:         targetRoom.Id,
				CheckInDate:    checkIn1,
				CheckOutDate:   checkOut1,
				NumberOfGuests: 1,
				GuestName:      "Guest 1",
				GuestEmail:     openapi_types.Email("g1@example.com"),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(book1.StatusCode()).To(Equal(http.StatusCreated))

			// Second booking adjacent: Nov 3 to Nov 5 (check-in == previous check-out)
			checkIn2 := openapi_types.Date{Time: time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)}
			checkOut2 := openapi_types.Date{Time: time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC)}

			book2, err := apiClient.BookRoomWithResponse(ctx, client.BookRoomJSONRequestBody{
				RoomId:         targetRoom.Id,
				CheckInDate:    checkIn2,
				CheckOutDate:   checkOut2,
				NumberOfGuests: 1,
				GuestName:      "Guest 2",
				GuestEmail:     openapi_types.Email("g2@example.com"),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(book2.StatusCode()).To(Equal(http.StatusCreated))
		})

		It("rejects booking exceeding room capacity (HTTP 400)", func(ctx context.Context) {
			checkIn := openapi_types.Date{Time: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)}
			checkOut := openapi_types.Date{Time: time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)}

			// Find single room (capacity 1)
			availResp, err := apiClient.FindAvailableRoomsWithResponse(ctx, &client.FindAvailableRoomsParams{
				CheckInDate:  checkIn,
				CheckOutDate: checkOut,
				Guests:       1,
			})
			Expect(err).NotTo(HaveOccurred())

			var singleRoom client.AvailableRoom
			for _, r := range availResp.JSON200.Rooms {
				if r.Capacity == 1 {
					singleRoom = r
					break
				}
			}
			Expect(singleRoom.Id).NotTo(Equal(uuid.Nil))

			// Try booking for 2 guests
			bookResp, err := apiClient.BookRoomWithResponse(ctx, client.BookRoomJSONRequestBody{
				RoomId:         singleRoom.Id,
				CheckInDate:    checkIn,
				CheckOutDate:   checkOut,
				NumberOfGuests: 2,
				GuestName:      "Over Capacity Party",
				GuestEmail:     openapi_types.Email("over@example.com"),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(bookResp.StatusCode()).To(Equal(http.StatusBadRequest))
			Expect(bookResp.ContentType()).To(ContainSubstring("application/problem+json"))
		})

		It("handles concurrent booking attempts without double booking", func(ctx context.Context) {
			checkIn := openapi_types.Date{Time: time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)}
			checkOut := openapi_types.Date{Time: time.Date(2026, 12, 5, 0, 0, 0, 0, time.UTC)}

			availResp, err := apiClient.FindAvailableRoomsWithResponse(ctx, &client.FindAvailableRoomsParams{
				CheckInDate:  checkIn,
				CheckOutDate: checkOut,
				Guests:       1,
			})
			Expect(err).NotTo(HaveOccurred())
			targetRoom := availResp.JSON200.Rooms[0]

			// Launch concurrent bookings for the exact same room and dates
			const concurrency = 4
			statusCodes := make([]int, concurrency)
			var wg sync.WaitGroup

			for i := range concurrency {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					resp, bookErr := apiClient.BookRoomWithResponse(ctx, client.BookRoomJSONRequestBody{
						RoomId:         targetRoom.Id,
						CheckInDate:    checkIn,
						CheckOutDate:   checkOut,
						NumberOfGuests: 1,
						GuestName:      "Concurrent Guest",
						GuestEmail:     openapi_types.Email("concurrent@example.com"),
					})
					if bookErr == nil {
						statusCodes[idx] = resp.StatusCode()
					}
				}(i)
			}

			wg.Wait()

			createdCount := 0
			conflictCount := 0
			for _, sc := range statusCodes {
				switch sc {
				case http.StatusCreated:
					createdCount++
				case http.StatusConflict:
					conflictCount++
				}
			}

			Expect(createdCount).To(Equal(1), "exactly one concurrent booking must succeed")
			Expect(conflictCount).To(Equal(concurrency-1), "all other concurrent attempts must be rejected with 409 Conflict")
		})
	})
})
