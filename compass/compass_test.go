package compass

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestFromDegreesComppassPoints(t *testing.T) {
	deg := float32(0)
	for cp := N; cp <= NBW; cp++ {
		cd := FromDegrees(deg)
		assert.Equalf(t, cd.Bearing, cp, "%.2f mapped to %d, expected %d", cd.Degrees, cd.Bearing, cp)
		deg += 11.25
	}
}

func TestFromDegreesNorth(t *testing.T) {
	cd := FromDegrees(0.0)
	assert.Equalf(t, cd.Bearing, N, "%.2f mapped to %d, expected %d", cd.Degrees, cd.Bearing, N)
	cd = FromDegrees(360.0)
	assert.Equalf(t, cd.Bearing, N, "%.2f mapped to %d, expected %d", cd.Degrees, cd.Bearing, N)
}
