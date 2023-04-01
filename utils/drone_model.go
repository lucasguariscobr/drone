package utils

import (
	"encoding/json"
	"strconv"

	"github.com/rs/zerolog/log"
)

type DroneType string

const (
	Unknown          DroneType = ""
	QuadcopterSmall  DroneType = "quadcopter-small"
	QuadcopterLarge  DroneType = "quadcopter-large"
	PlaneSmall       DroneType = "plane-small"
	SingleRotorLarge DroneType = "single-rotor-large"
)

func fromString(s string) (DroneType, error) {
	switch s {
	case string(QuadcopterSmall):
		return QuadcopterSmall, nil
	case string(QuadcopterLarge):
		return QuadcopterLarge, nil
	case string(PlaneSmall):
		return PlaneSmall, nil
	case string(SingleRotorLarge):
		return SingleRotorLarge, nil
	}
	return Unknown, ErrCreateDroneType
}

type drone struct {
	Id                  string
	InstructionIndex    interface{}
	instructionIndexInt int
	Name                string
	Plan                []string
	DroneTypeRaw        string `json:"type"`
	droneType           DroneType
	Cost                json.RawMessage
	Status              string
}

// ValidateDroneModelJson receives a JSON RawMessage and validates it.
// The validation rules map the constraints from the drone model documentation.
// Returns an error if any issue is found.
func ValidateDroneModelJson(jsonRawData *json.RawMessage) error {
	log.Debug().Msg("Validating Drone Model JSON Payload...")

	var droneModel drone
	jsonErr := json.Unmarshal(*jsonRawData, &droneModel)
	if jsonErr != nil {
		return jsonErr
	}

	if string(droneModel.Cost) != "" {
		return ErrCreateDroneCost
	}

	if droneModel.Status != "" {
		return ErrCreateDroneStatus
	}

	instructionIndexError := validateInstructionIndex(&droneModel)
	if instructionIndexError != nil {
		return instructionIndexError
	}

	planError := validatePlan(&droneModel)
	if planError != nil {
		return planError
	}

	droneTypeError := validateType(&droneModel)
	if droneTypeError != nil {
		return droneTypeError
	}

	log.Debug().Msg("Validation Completed")
	return nil
}

func validateInstructionIndex(droneModel *drone) error {
	if droneModel.InstructionIndex == "" {
		return nil
	}

	switch droneModel.InstructionIndex.(type) {
	case int:
		droneModel.instructionIndexInt = droneModel.InstructionIndex.(int)
	case string:
		var instrIndexErr error
		droneModel.instructionIndexInt, instrIndexErr = strconv.Atoi(droneModel.InstructionIndex.(string))
		if instrIndexErr != nil {
			return ErrCreateDroneInstructionIndex
		}
	}

	return nil
}

func validatePlan(droneModel *drone) error {
	if len(droneModel.Plan) == 0 {
		return ErrCreateDronePlanLength
	}

	lastPlanInstruction := droneModel.Plan[len(droneModel.Plan)-1]
	// Assuming that's the only valid landing command
	if lastPlanInstruction != "land-drone" {
		return ErrCreateDronePlanLastInstruction
	}
	return nil
}

func validateType(droneModel *drone) error {
	if droneModel.DroneTypeRaw == "" {
		return ErrCreateDroneMissingType
	}

	var droneTypeErr error
	droneModel.droneType, droneTypeErr = fromString(droneModel.DroneTypeRaw)
	if droneTypeErr != nil {
		return droneTypeErr
	}

	return nil
}
