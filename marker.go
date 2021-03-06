package honeycombio

import (
	"time"
)

// Compile-time proof of interface implementation.
var _ Markers = (*markers)(nil)

// Markers describes all the markers related methods that Honeycomb supports.
type Markers interface {
	// List all markers present in this dataset.
	List() ([]Marker, error)

	// Get a marker by its ID. Returns nil, ErrNotFound if there is no marker
	// with the given ID.
	//
	// This method calls List internally since there is no API available to
	// directly get a single marker.
	Get(id string) (*Marker, error)

	// Create a new marker in this dataset.
	Create(data MarkerCreateData) (*Marker, error)
}

// markers implements Markers.
type markers struct {
	client *Client
}

// Marker represents a Honeycomb marker, as described by https://docs.honeycomb.io/api/markers/#fields-on-a-marker
type Marker struct {
	ID string `json:"id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// StartTime unix timestamp truncates to seconds
	StartTime int64 `json:"start_time,omitempty"`
	// EndTime unix timestamp truncates to seconds
	EndTime int64 `json:"end_time,omitempty"`
	// Message is optional free-form text associated with the message
	Message string `json:"message,omitempty"`
	// Type is an optional marker identifier, eg 'deploy' or 'chef-run'
	Type string `json:"type,omitempty"`
	// URL is an optional url associated with the marker
	URL string `json:"url,omitempty"`
	// Color is not stored in the marker table but populated by a join
	Color string `json:"color,omitempty"`
}

func (s *markers) List() ([]Marker, error) {
	req, err := s.client.newRequest("GET", "/1/markers/"+s.client.dataset, nil)
	if err != nil {
		return nil, err
	}

	var m []Marker
	err = s.client.do(req, &m)
	return m, err
}

func (s *markers) Get(ID string) (*Marker, error) {
	markers, err := s.List()
	if err != nil {
		return nil, err
	}

	for _, m := range markers {
		if m.ID == ID {
			return &m, nil
		}
	}
	return nil, ErrNotFound
}

// MarkerCreateData holds the data to create a new marker.
type MarkerCreateData struct {
	StartTime int64  `json:"start_time,omitempty"`
	EndTime   int64  `json:"end_time,omitempty"`
	Message   string `json:"message,omitempty"`
	Type      string `json:"type,omitempty"`
	URL       string `json:"url,omitempty"`
}

func (s *markers) Create(d MarkerCreateData) (*Marker, error) {
	req, err := s.client.newRequest("POST", "/1/markers/"+s.client.dataset, d)
	if err != nil {
		return nil, err
	}

	var m Marker
	err = s.client.do(req, &m)
	return &m, err
}
