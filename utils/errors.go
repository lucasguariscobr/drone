package utils

import (
	"errors"
	"net/http"
)

// ErrorBuilder will return custom errors for the HTTP
// Status Codes that were mapped by the API.
func ErrorBuilder(httpStatusCode int) error {
	switch httpStatusCode {
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	case http.StatusInternalServerError:
		return ErrInternalServer
	}

	return nil
}

var ErrUnauthorized = errors.New("unauthorized access. Check that DRONE_TOKEN was correctly set and that DRONE_ADDR is pointing to the correct backend")
var ErrTooManyRequests = errors.New("too many requests. Consider setting DRONE_MAX_RETRIES to define a maximum number of retries when certain errors codes are encountered")
var ErrBadRequest = errors.New("bad Request")
var ErrInternalServer = errors.New("internal Server Error")
var ErrMissingToken = errors.New("no Authorization token available. Configure DRONE_TOKEN")
var ErrMissingAddr = errors.New("no API Address available. Configure DRONE_ADDR")
var ErrCreateDroneCost = errors.New("new drones should not include cost")
var ErrCreateDroneStatus = errors.New("new drones should not include status")
var ErrCreateDronePlanLength = errors.New("the drone plan must be at least one instruction long")
var ErrCreateDronePlanLastInstruction = errors.New("the last instruction from the drone plan must be a landing command")
var ErrCreateDroneType = errors.New("the drone type must be one of: quadcopter-small, quadcopter-large, plane-small, single-rotor-large")
var ErrCreateDroneInstructionIndex = errors.New("the instruction index couldn't be converted to a numeric index")
var ErrCreateDroneMissingType = errors.New("the type is required to create a drone")
