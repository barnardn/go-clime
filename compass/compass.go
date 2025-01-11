package compass

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
		deltaMin, deltaMax := azimuth - 5.625, azimuth + 5.625
		if deg > deltaMin && deg < deltaMax {
			return idx
		}
	}
	return N
}
