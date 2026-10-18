package usecase

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/netologist/hotel-booking/internal/domain"
)

type BookingNotifier interface {
	SendConfirmationEmail(booking domain.Booking)
}

type AsyncEmailNotifier struct {
	delay time.Duration
	wg    sync.WaitGroup
}

func NewAsyncEmailNotifier(delay time.Duration) *AsyncEmailNotifier {
	return &AsyncEmailNotifier{delay: delay}
}

func (n *AsyncEmailNotifier) SendConfirmationEmail(b domain.Booking) {
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		if n.delay > 0 {
			time.Sleep(n.delay)
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

// Drain waits for all background email simulations to complete, or until ctx expires.
func (n *AsyncEmailNotifier) Drain(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		n.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
