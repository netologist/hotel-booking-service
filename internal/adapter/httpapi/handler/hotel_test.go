package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/netologist/hotel-booking/internal/adapter/httpapi/handler"
	"github.com/netologist/hotel-booking/internal/domain"
	usecasemocks "github.com/netologist/hotel-booking/internal/gen/mocks/usecase"
	"github.com/netologist/hotel-booking/internal/gen/openapi"
)

const (
	testHotelName          = "Grand Budapest Hotel"
	testHotelAddress       = "100 Alpine Way"
	testHotelCity          = "Lutz"
	validCreateHotelJSON   = `{"name":"Grand Budapest Hotel","address":"100 Alpine Way","city":"Lutz"}`
	contentTypeApplication = "application/json"
	fieldValidationFailed  = "validation_failed"
	fieldName              = "name"
	fieldAddress           = "address"
	fieldCity              = "city"
	errInternalServer      = "INTERNAL_SERVER_ERROR"
)

func TestNewHotelHandler(t *testing.T) {
	mockService := usecasemocks.NewMockHotelService(t)
	h := handler.NewHotelHandler(mockService)
	require.NotNil(t, h)
}

func TestHotelHandler_CreateHotel_Success(t *testing.T) {
	hotelID := uuid.New()
	fixedTime := time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)
	mockService := usecasemocks.NewMockHotelService(t)
	h := handler.NewHotelHandler(mockService)

	mockService.EXPECT().CreateHotel(mock.Anything, testHotelName, testHotelAddress, testHotelCity).Return(domain.Hotel{
		ID:        hotelID,
		Name:      testHotelName,
		Address:   testHotelAddress,
		City:      testHotelCity,
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}, nil).Once()

	resp, err := h.CreateHotel(t.Context(), openapi.CreateHotelRequestObject{
		Body: &openapi.CreateHotelJSONRequestBody{
			Name:    testHotelName,
			Address: testHotelAddress,
			City:    testHotelCity,
		},
	})
	require.NoError(t, err)

	jsonResp, ok := resp.(openapi.CreateHotel201JSONResponse)
	require.True(t, ok)
	assert.Equal(t, hotelID, jsonResp.Id)
	assert.Equal(t, testHotelName, jsonResp.Name)
	assert.Equal(t, testHotelAddress, jsonResp.Address)
	assert.Equal(t, testHotelCity, jsonResp.City)
	assert.True(t, fixedTime.Equal(jsonResp.CreatedAt))
	assert.True(t, fixedTime.Equal(jsonResp.UpdatedAt))
}

func TestHotelHandler_CreateHotel_MalformedJSON(t *testing.T) {
	router, _ := SetupTestRouterWithMock(t)

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/api/v1/hotels",
		strings.NewReader(`{"name": "invalid-json`),
	)
	req.Header.Set("Content-Type", contentTypeApplication)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/problem+json")
	var prob openapi.Problem
	err := json.NewDecoder(rec.Body).Decode(&prob)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, prob.Status)
	assert.Equal(t, "Bad Request", prob.Title)
	require.NotNil(t, prob.Detail)
	assert.NotEmpty(t, *prob.Detail)
}

type validationTestCase struct {
	name          string
	body          openapi.CreateHotelRequest
	expectedField string
}

func missingFieldTestCases() []validationTestCase {
	return []validationTestCase{
		{
			name: "missing_name",
			body: openapi.CreateHotelRequest{
				Name:    "",
				Address: testHotelAddress,
				City:    testHotelCity,
			},
			expectedField: fieldName,
		},
		{
			name: "missing_address",
			body: openapi.CreateHotelRequest{
				Name:    testHotelName,
				Address: "",
				City:    testHotelCity,
			},
			expectedField: fieldAddress,
		},
		{
			name: "missing_city",
			body: openapi.CreateHotelRequest{
				Name:    testHotelName,
				Address: testHotelAddress,
				City:    "",
			},
			expectedField: fieldCity,
		},
	}
}

func boundaryFieldTestCases() []validationTestCase {
	return []validationTestCase{
		{
			name: "name_too_short",
			body: openapi.CreateHotelRequest{
				Name:    "A",
				Address: testHotelAddress,
				City:    testHotelCity,
			},
			expectedField: fieldName,
		},
		{
			name: "name_too_long",
			body: openapi.CreateHotelRequest{
				Name:    strings.Repeat("A", 256),
				Address: testHotelAddress,
				City:    testHotelCity,
			},
			expectedField: fieldName,
		},
		{
			name: "address_too_short",
			body: openapi.CreateHotelRequest{
				Name:    testHotelName,
				Address: "A",
				City:    testHotelCity,
			},
			expectedField: fieldAddress,
		},
		{
			name: "address_too_long",
			body: openapi.CreateHotelRequest{
				Name:    testHotelName,
				Address: strings.Repeat("A", 256),
				City:    testHotelCity,
			},
			expectedField: fieldAddress,
		},
		{
			name: "city_too_short",
			body: openapi.CreateHotelRequest{
				Name:    testHotelName,
				Address: testHotelAddress,
				City:    "A",
			},
			expectedField: fieldCity,
		},
		{
			name: "city_too_long",
			body: openapi.CreateHotelRequest{
				Name:    testHotelName,
				Address: testHotelAddress,
				City:    strings.Repeat("A", 101),
			},
			expectedField: fieldCity,
		},
	}
}

func TestHotelRouter_CreateHotel_ValidationFailure(t *testing.T) {
	testCases := append(missingFieldTestCases(), boundaryFieldTestCases()...)
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			router, _ := SetupTestRouterWithMock(t)

			payload, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req := httptest.NewRequestWithContext(
				t.Context(),
				http.MethodPost,
				"/api/v1/hotels",
				bytes.NewReader(payload),
			)
			req.Header.Set("Content-Type", contentTypeApplication)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Header().Get("Content-Type"), "application/problem+json")

			var prob openapi.Problem
			err = json.NewDecoder(rec.Body).Decode(&prob)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, prob.Status)
			assert.Equal(t, "Bad Request", prob.Title)
			require.NotNil(t, prob.Detail)
			assert.NotEmpty(t, *prob.Detail)
		})
	}
}

func TestHotelHandler_CreateHotel_ServiceError(t *testing.T) {
	mockService := usecasemocks.NewMockHotelService(t)
	h := handler.NewHotelHandler(mockService)

	serviceErr := errors.New("database connection failed")
	mockService.EXPECT().
		CreateHotel(mock.Anything, testHotelName, testHotelAddress, testHotelCity).
		Return(domain.Hotel{}, serviceErr).
		Once()

	resp, err := h.CreateHotel(t.Context(), openapi.CreateHotelRequestObject{
		Body: &openapi.CreateHotelJSONRequestBody{
			Name:    testHotelName,
			Address: testHotelAddress,
			City:    testHotelCity,
		},
	})
	require.ErrorIs(t, err, serviceErr)
	assert.Nil(t, resp)
}

func TestHotelHandler_CreateHotel_DomainAndBusinessErrors(t *testing.T) {
	testCases := []struct {
		name           string
		mockErr        error
		expectedStatus int
		expectedTitle  string
		expectedCode   string
	}{
		{
			name:           "domain_invariant_error",
			mockErr:        domain.ErrInvalidHotel,
			expectedStatus: http.StatusBadRequest,
			expectedTitle:  "Bad Request",
			expectedCode:   "INVALID_HOTEL",
		},
		{
			name:           "conflict_duplicate_error",
			mockErr:        domain.ErrHotelAlreadyExists,
			expectedStatus: http.StatusConflict,
			expectedTitle:  "Conflict",
			expectedCode:   "HOTEL_ALREADY_EXISTS",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := usecasemocks.NewMockHotelService(t)
			h := handler.NewHotelHandler(mockService)

			mockService.EXPECT().
				CreateHotel(mock.Anything, testHotelName, testHotelAddress, testHotelCity).
				Return(domain.Hotel{}, tc.mockErr).
				Once()

			resp, err := h.CreateHotel(t.Context(), openapi.CreateHotelRequestObject{
				Body: &openapi.CreateHotelJSONRequestBody{
					Name:    testHotelName,
					Address: testHotelAddress,
					City:    testHotelCity,
				},
			})
			require.NoError(t, err)

			var prob openapi.Problem
			switch r := resp.(type) {
			case openapi.CreateHotel400ApplicationProblemPlusJSONResponse:
				prob = openapi.Problem(r.BadRequestApplicationProblemPlusJSONResponse)
			case openapi.CreateHotel409ApplicationProblemPlusJSONResponse:
				prob = openapi.Problem(r.ConflictApplicationProblemPlusJSONResponse)
			default:
				t.Fatalf("unexpected response type: %T", resp)
			}
			assert.Equal(t, tc.expectedStatus, prob.Status)
			assert.Equal(t, tc.expectedTitle, prob.Title)
			require.NotNil(t, prob.Detail)
			assert.Equal(t, tc.mockErr.Error(), *prob.Detail)
			require.NotNil(t, prob.Code)
			assert.Equal(t, tc.expectedCode, *prob.Code)
		})
	}
}

func TestHotelRouter_Integration_CreateHotel(t *testing.T) {
	hotelID := uuid.New()
	fixedTime := time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)
	router, mockService := SetupTestRouterWithMock(t)

	mockService.EXPECT().CreateHotel(mock.Anything, testHotelName, testHotelAddress, testHotelCity).Return(domain.Hotel{
		ID:        hotelID,
		Name:      testHotelName,
		Address:   testHotelAddress,
		City:      testHotelCity,
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}, nil).Once()

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/api/v1/hotels",
		strings.NewReader(validCreateHotelJSON),
	)
	req.Header.Set("Content-Type", contentTypeApplication)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var got openapi.Hotel
	err := json.NewDecoder(rec.Body).Decode(&got)
	require.NoError(t, err)
	assert.Equal(t, hotelID, got.Id)
	assert.Equal(t, testHotelName, got.Name)
}
func TestHotelHandler_FindHotels_Success(t *testing.T) {
	hotelID := uuid.New()
	fixedTime := time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)
	mockService := usecasemocks.NewMockHotelService(t)
	h := handler.NewHotelHandler(mockService)

	hotels := []domain.Hotel{
		{
			ID:        hotelID,
			Name:      testHotelName,
			Address:   testHotelAddress,
			City:      testHotelCity,
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
	}

	mockService.EXPECT().FindHotels(mock.Anything, "Budapest").Return(hotels, nil).Once()

	resp, err := h.FindHotels(t.Context(), openapi.FindHotelsRequestObject{
		Params: openapi.FindHotelsParams{Name: "Budapest"},
	})
	require.NoError(t, err)

	jsonResp, ok := resp.(openapi.FindHotels200JSONResponse)
	require.True(t, ok)
	require.NotNil(t, jsonResp.Hotels)
	require.Len(t, *jsonResp.Hotels, 1)
	assert.Equal(t, hotelID, (*jsonResp.Hotels)[0].Id)
	assert.Equal(t, testHotelName, (*jsonResp.Hotels)[0].Name)
}

func TestHotelHandler_FindHotels_EmptyQuery(t *testing.T) {
	mockService := usecasemocks.NewMockHotelService(t)
	h := handler.NewHotelHandler(mockService)

	resp, err := h.FindHotels(t.Context(), openapi.FindHotelsRequestObject{
		Params: openapi.FindHotelsParams{Name: "   "},
	})
	require.NoError(t, err)

	badReq, ok := resp.(openapi.FindHotels400ApplicationProblemPlusJSONResponse)
	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, badReq.Status)
	assert.Equal(t, "Bad Request", badReq.Title)
	require.NotNil(t, badReq.Detail)
	assert.Equal(t, "name query parameter is required", *badReq.Detail)
}

func TestHotelHandler_FindHotels_ServiceError(t *testing.T) {
	mockService := usecasemocks.NewMockHotelService(t)
	h := handler.NewHotelHandler(mockService)

	serviceErr := errors.New("database failure")
	mockService.EXPECT().FindHotels(mock.Anything, "Budapest").Return(nil, serviceErr).Once()

	resp, err := h.FindHotels(t.Context(), openapi.FindHotelsRequestObject{
		Params: openapi.FindHotelsParams{Name: "Budapest"},
	})
	require.ErrorIs(t, err, serviceErr)
	assert.Nil(t, resp)
}

func TestHotelRouter_Integration_FindHotels(t *testing.T) {
	hotelID := uuid.New()
	fixedTime := time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)
	router, mockService := SetupTestRouterWithMock(t)

	hotels := []domain.Hotel{
		{
			ID:        hotelID,
			Name:      testHotelName,
			Address:   testHotelAddress,
			City:      testHotelCity,
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
	}

	mockService.EXPECT().FindHotels(mock.Anything, "Budapest").Return(hotels, nil).Once()

	req := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		"/api/v1/hotels?name=Budapest",
		nil,
	)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var got struct {
		Hotels []openapi.Hotel `json:"hotels"`
	}
	err := json.NewDecoder(rec.Body).Decode(&got)
	require.NoError(t, err)
	require.Len(t, got.Hotels, 1)
	assert.Equal(t, hotelID, got.Hotels[0].Id)
	assert.Equal(t, testHotelName, got.Hotels[0].Name)
}
