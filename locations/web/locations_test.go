package main

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/kchimev/locations-api/locations/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type LocationEntityMock struct {
	mock.Mock
}

func (m *LocationEntityMock) Get(lat float64, lon float64, radius int) ([]models.Location, error) {
	argsList := m.Called(mock.Anything, mock.Anything, mock.Anything)

	return argsList.Get(0).([]models.Location), argsList.Error(1)
}

func (*LocationEntityMock) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (*LocationEntityMock) Close() error {
	return nil
}

type POIEntityMock struct {
	mock.Mock
}

func (m *POIEntityMock) Get(lat float64, lon float64, radius int) ([]models.POI, error) {
	argsList := m.Called(mock.Anything, mock.Anything, mock.Anything)

	return argsList.Get(0).([]models.POI), argsList.Error(1)
}

func (*POIEntityMock) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (*POIEntityMock) Close() error {
	return nil
}

func TestFetchLocationData(t *testing.T) {
	testCases := []struct {
		name     string
		expected *ResponsePayload
		error    error
		prepare  func(*LocationsService)
	}{
		{
			name:  "Valid request found results",
			error: nil,
			expected: &ResponsePayload{
				Locations: []models.Location{
					{
						Name: "Test Location",
					},
				},
				POIs: []models.POI{
					{
						Name: "Test POI",
					},
				},
			},
			prepare: func(loc *LocationsService) {
				locEnt := new(LocationEntityMock)
				poiEnt := new(POIEntityMock)
				locEnt.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(
					[]models.Location{
						{
							Name: "Test Location",
						},
					},
					nil,
				)
				poiEnt.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(
					[]models.POI{
						{
							Name: "Test POI",
						},
					},
					nil,
				)
				loc.locEnt = locEnt
				loc.poiEnt = poiEnt
			},
		},
		{
			name:  "Valid request no results",
			error: nil,
			expected: &ResponsePayload{
				Locations: []models.Location{},
				POIs:      []models.POI{},
			},
			prepare: func(loc *LocationsService) {
				locEnt := new(LocationEntityMock)
				poiEnt := new(POIEntityMock)
				locEnt.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(
					[]models.Location{},
					nil,
				)
				poiEnt.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(
					[]models.POI{},
					nil,
				)
				loc.locEnt = locEnt
				loc.poiEnt = poiEnt
			},
		},
		{
			name:     "Valid request error while fetching",
			expected: nil,
			error:    errors.New("Fetching from database failed"),
			prepare: func(loc *LocationsService) {
				locEnt := new(LocationEntityMock)
				poiEnt := new(POIEntityMock)
				locEnt.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(
					[]models.Location{},
					errors.New("Fetching from database failed"),
				)
				poiEnt.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(
					[]models.POI{},
					nil,
				)
				loc.locEnt = locEnt
				loc.poiEnt = poiEnt
			},
		},
	}

	for _, test := range testCases {
		locSer := &LocationsService{}
		t.Run(test.name, func(t *testing.T) {
			test.prepare(locSer)
			res, err := locSer.fetchLocationData(42.3, 23.3)
			if err != nil {
				assert.Equal(t, err, test.error)
			} else {
				assert.Equal(t, res, test.expected)
			}
		})
	}
}
