package audience

import "fmt"

// Payload is a response json.
type Payload struct {
	Segment Segment `json:"segment"`
}

// Print prints payload.
func (p *Payload) Print() string {
	format := `
id: %d,
type: %s,
status: %s,
has_guests: %t,
guest_quantity: %d,
can_create_dependent: %t,
has_derivatives: %t,
hashed: %t,
item_quantity: %d,
guest: %t
`
	return fmt.Sprintf(format,
		p.Segment.ID,
		p.Segment.Type,
		p.Segment.Status,
		p.Segment.HasGuests,
		p.Segment.GuestQuantity,
		p.Segment.CanCreateDependent,
		p.Segment.HasDerivatives,
		p.Segment.Hashed,
		p.Segment.ItemQuantity,
		p.Segment.Guest,
	)
}

// Segment represents segment params.
type Segment struct {
	ID                 int    `json:"id"`
	Type               string `json:"type"`
	Status             string `json:"status"`
	HasGuests          bool   `json:"has_guests"`
	GuestQuantity      int    `json:"guest_quantity"`
	CanCreateDependent bool   `json:"can_create_dependent"`
	HasDerivatives     bool   `json:"has_derivatives"`
	Hashed             bool   `json:"hashed"`
	ItemQuantity       int    `json:"item_quantity"`
	Guest              bool   `json:"guest"`
}

// Confirm represents confirm payload.
type Confirm struct {
	Segment SegmentConfirm `json:"segment"`
}

// SegmentConfirm represents confirm segment request.
type SegmentConfirm struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Hashed      uint   `json:"hashed"`
	ContentType string `json:"content_type"`
}

// ErrorResponse is returned in case of error.
type ErrorResponse struct {
	Errors  []Error `json:"errors"`
	Code    int     `json:"code"`
	Message string  `json:"message"`
}

// Error contains error type and message.
type Error struct {
	Type    string `json:"error_type"`
	Message string `json:"message"`
}
