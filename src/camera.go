package main

import (
	"math"
	"math/rand"
)

type camera struct {
	lowerLeft, hor, vert, origin vec3
	u, v, w                      vec3
	lensRadius, shutter          float64
}

func cam(lookFrom, lookAt vec3, fov, aperture, shutter float64) *camera {
	c := &camera{}
	c.lensRadius = aperture / 2.0
	c.shutter = shutter
	theta := fov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := (float64(width) / float64(height)) * halfHeight
	focusDist := 10.0

	// This is used to calculate the direction of the camera.
	c.w = lookFrom.sub(lookAt).normalize() // The difference from the target and position, will give the direction.
	c.u = cross(vec(0.0, 1.0, 0.0), c.w).normalize()
	c.v = cross(c.w, c.u)

	c.origin = lookFrom
	c.lowerLeft = c.origin.sub(c.u.mulScalar(halfWidth * focusDist)).sub(c.v.mulScalar(halfHeight * focusDist)).sub(c.w.mulScalar(focusDist))
	c.hor = c.u.mulScalar(2.0 * halfWidth * focusDist)
	c.vert = c.v.mulScalar(2.0 * halfHeight * focusDist)

	return c
}

func (c *camera) ray(s, t float64, rnd *rand.Rand) ray {
	// Better performance, because there are no random numbers needed.
	// Even if we would calculate them, we would multiply by zero, so this is useless.
	if c.lensRadius == 0.0 {
		return ray{
			c.origin,
			c.lowerLeft.add(c.hor.mulScalar(s).add(c.vert.mulScalar(t))).sub(c.origin),
			rnd.Float64() * c.shutter,
		}
	}

	// Add blur, when there is a higher radius.
	rd := randInDisk(rnd).mulScalar(c.lensRadius)
	offset := c.u.mulScalar(rd.x).add(c.v.mulScalar(rd.y))

	return ray{
		c.origin.add(offset),
		c.lowerLeft.add(c.hor.mulScalar(s).add(c.vert.mulScalar(t))).sub(c.origin).sub(offset),
		rnd.Float64() * c.shutter,
	}
}

func randInDisk(rnd *rand.Rand) vec3 {
	p := vec(rnd.Float64(), rnd.Float64(), 0.0).mulScalar(2.0).sub(vec(1.0, 1.0, 0.0))
	for dot(p, p) >= 1.0 {
		p = vec(rnd.Float64(), rnd.Float64(), 0.0).mulScalar(2.0).sub(vec(1.0, 1.0, 0.0))
	}
	return p
}
