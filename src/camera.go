package main

import (
	"math"
)

type camera struct {
	lowerLeft, hor, vert, origin vec3
}

func cam(lookFrom, lookAt vec3, fov, aspect float64) camera {
	theta := fov * math.Pi / 180
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := aspect * halfHeight

	// This is used to calculate the direction of the camera.
	w := lookFrom.sub(lookAt).normalize()	// The difference from the target and position, will give the direction.
	u := cross(v(0.0, 1.0, 0.0), w).normalize()
	v := cross(w, u)

	return camera{
		lookFrom.sub(u.mulScalar(halfWidth)).sub(v.mulScalar(halfHeight)).sub(w),
		u.mulScalar(2.0 * halfWidth),
		v.mulScalar(2.0 * halfHeight),
		lookFrom,
	}
}

func (c camera) getRay(u, v float64) ray {
	return ray{
		c.origin,
		c.lowerLeft.add(c.hor.mulScalar(u).add(c.vert.mulScalar(v))).sub(c.origin),
	}
}
