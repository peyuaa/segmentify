// Package handlers Segmentify API.
//
// Documentation for Segmentify API.
//
// Schemes: http
// BasePath: /
// Version: 0.0.1
// Contact: Dmitriy Krasnov<dk.peyuaa@gmail.com>
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import "github.com/peyuaa/segmentify/models"

// A segment with the specified slug returns in the response
// swagger:response segmentResponse
type segmentResponse struct {
	// A segment with the specified slug
	// in: body
	Body models.Segment
}

// A created segment returns in the response
// swagger:response createSegmentResponse
type createSegmentResponse struct {
	// A segment with the specified slug
	// in: body
	Body models.Segment

	// A link to the created segment
	Location string
}

// A list of segments returns in the response
// swagger:response segmentsResponse
type segmentsResponse struct {
	// All active segments in the system
	// in: body
	Body models.Segments
}

// swagger:parameters deleteSegment
type segmentSlugParameterWrapper struct {
	// The slug of the segment to delete from the database
	// in: path
	// required: true
	Slug string
}

// swagger:response noContentResponse
type segmentNoContentResponse struct {
}

// swagger:response errorResponse
type segmentErrorResponse struct {
	// The error message
	// in: body
	Body GenericError
}

// swagger:response userHistoryResponse
type userHistoryResponse struct {
	// link to csv file with user's segments history for specified period
	// in: body
	Body models.UserHistoryResponse
}

// swagger:response activeSegmentsResponse
type activeSegmentsResponse struct {
	// list of active segments for specified user
	// in: body
	Body models.ActiveSegmentsResponse
}
