package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/netologist/hotel-booking/pkg/client"
)

func TestClient_CreateHotel(t *testing.T) {
	t.Parallel()

	hotelID := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/hotels", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var req client.CreateHotelJSONRequestBody
		if err := json.NewDecoder(r.Body).Decode(&req); !assert.NoError(t, err) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		created := client.Hotel{
			Id:        hotelID,
			Name:      req.Name,
			Address:   req.Address,
			City:      req.City,
			CreatedAt: now,
			UpdatedAt: now,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(created)
	}))
	defer server.Close()

	c, err := client.NewClientWithResponses(server.URL)
	require.NoError(t, err)

	resp, err := c.CreateHotelWithResponse(context.Background(), client.CreateHotelJSONRequestBody{
		Name:    "Grand Budapest Hotel",
		Address: "100 Alpine Way",
		City:    "Lutz",
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode())
	require.NotNil(t, resp.JSON201)
	assert.Equal(t, "Grand Budapest Hotel", resp.JSON201.Name)
	assert.Equal(t, "100 Alpine Way", resp.JSON201.Address)
	assert.Equal(t, "Lutz", resp.JSON201.City)
}
