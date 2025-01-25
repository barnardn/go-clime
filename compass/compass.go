package compass

import (
	"fmt"
)

type CompassDirection struct {
	Degrees float32
	Bearing CompassBearing
}

type CompassBearing int8

const (
	N CompassBearing = iota
	NBE
	NNE
	NEBN
	NE
	NEBE
	ENE
	EBN
	E
	EBS
	ESE
	SEBE
	SE
	SEBS
	SSE
	SBE
	S
	SBW
	SSW
	SWBS
	SW
	SWBW
	WSW
	WBS
	W
	WBN
	WNW
	NWBW
	NW
	NWBN
	NNW
	NBW
)

func FromDegrees(deg float32) *CompassDirection {
	return &CompassDirection{
		Degrees: deg,
		Bearing: degreeToBearing(deg),
	}
}

func degreeToBearing(deg float32) CompassBearing {
	if deg >= 360.0 {
		deg -= 360.0
	}
	for idx := N; idx <= NBW; idx++ {
		azimuth := float32(idx) * 11.25
		deltaMin, deltaMax := azimuth-5.625, azimuth+5.625
		if deg > deltaMin && deg < deltaMax {
			return idx
		}
	}
	return N
}

func (c *CompassBearing) String() string {
	direction := ""
	switch *c {
	case N:
		direction = "North"
	case NEBN:
		direction = "NE by North"
	case NEBE:
		direction = "NE by East"
	case NE:
		direction = "Northeast"
	case NBE:
		direction = "N by East"
	case NW:
		direction = "Northwest"
	case NBW:
		direction = "N by West"
	case NNW:
		direction = "North Northwest"
	case NNE:
		direction = "North Northeast"
	case NWBN:
		direction = "NW by North"
	case EBN:
		direction = "East by N"
	case ENE:
		direction = "East by NE"
	case E:
		direction = "East"
	case ESE:
		direction = "East SouthEast"
	case EBS:
		direction = "E by South"
	case W:
		direction = "West"
	case WBN:
		direction = "W by North"
	case WNW:
		direction = "West Northwest"
	case WBS:
		direction = "W by South"
	case WSW:
		direction = "West Southwest"
	case S:
		direction = "South"
	case SE:
		direction = "Southeast"
	case SW:
		direction = "Southwest"
	case SBE:
		direction = "S by East"
	case SBW:
		direction = "S by West"
	case SSE:
		direction = "South Southeast"
	case SEBE:
		direction = "SE by South"
	case SEBS:
		direction = "SE by South"
	case SSW:
		direction = "South Southwest"
	case SWBW:
		direction = "SW by West"
	case SWBS:
		direction = "SW by South"
	default:
		direction = fmt.Sprintf("Missed %d", *c)
	}
	return direction
}

func (c *CompassDirection) String() string {
	return fmt.Sprintf("%0.2f° %s", c.Degrees, c.Bearing.String())
}
