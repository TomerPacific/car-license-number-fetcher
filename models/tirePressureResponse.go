package tirePressure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

type TirePressurePayload struct {
	Make         string `json:"make"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	Modification string `json:"modification,omitempty"`
	TireSize     string `json:"tireSize,omitempty"`
}

type TirePressureResponse struct {
	Source   string      `json:"source"`
	FrontPsi *float64    `json:"frontPsi,omitempty"`
	RearPsi  *float64    `json:"rearPsi,omitempty"`
	Unit     string      `json:"unit,omitempty"`
	Note     string      `json:"note,omitempty"`
	Raw      interface{} `json:"raw,omitempty"`
}
