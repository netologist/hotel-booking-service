package usecase

import (
	"log/slog"
	"time"

	"github.com/netologist/hotel-booking/internal/domain"
)

type BookingNotifier interface {
	SendConfirmationEmail(booking domain.Booking)
}

type AsyncEmailNotifier struct {
	delay time.Duration
}

func NewAsyncEmailNotifier(delay time.Duration) *AsyncEmailNotifier {
	return &AsyncEmailNotifier{delay: delay}
}

func (n *AsyncEmailNotifier) SendConfirmationEmail(b domain.Booking) {
	delay := n.delay
	go func() {
		if delay > 0 {
			time.Sleep(delay)
		}
		slog.Info("booking confirmation email sent asynchronously",
			"reference", b.Reference,
			"guest_name", b.GuestName,
			"guest_email", b.GuestEmail,
			"room_id", b.RoomID,
			"hotel_id", b.HotelID,
			"check_in", b.CheckInDate.Format("2006-01-02"),
			"check_out", b.CheckOutDate.Format("2006-01-02"),
		)
	}()
}
