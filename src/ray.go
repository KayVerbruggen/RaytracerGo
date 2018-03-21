package main

import (
	"math"
	"math/rand"
)

// ray is used for tracing a line from a origin(O) in a direction(Dir).
type ray struct {
	origin, dir vec3
	time float64
}

// PointAtParam gets a vec3 position at a certain distance across the line.
func (r *ray) point(t float64) vec3 {
	return r.origin.add(r.dir.mulScalar(t))
}

// Color returns a color based on what the ray hits.
func (r *ray) color(s *scene, depth int64, rnd *rand.Rand) vec3 {
	hr := hitRecord{}
	if s.hit(*r, 0.001, math.MaxFloat64, &hr) {
		scattered := ray{}
		attenuation := vec3{}
		if depth < 50 && hr.mat.scatter(*r, &hr, &attenuation, &scattered, rnd) {
			return attenuation.mul(scattered.color(s, depth+1, rnd))
		}
		return vec(0.0, 0.0, 0.0)
	}

	nd := r.dir.normalize()
	t := 0.5 * (nd.y + 1.0)

	// 		(1.0-t) * (1.0, 1.0, 1.0) + t * (0.5, 0.7, 1.0)
	temp := vec(1.0, 1.0, 1.0).mulScalar(1.0 - t).add(vec(0.5, 0.7, 1.0).mulScalar(t))

	return temp
}

func randInUnitSphere(rnd *rand.Rand) vec3 {
	p := vec(rnd.Float64(), rnd.Float64(), rnd.Float64()).mulScalar(2.0).sub(vec(1.0, 1.0, 1.0))
	for p.lengthSqr() >= 1.0 {
		p = vec(rnd.Float64(), rnd.Float64(), rnd.Float64()).mulScalar(2.0).sub(vec(1.0, 1.0, 1.0))
	}

	return p
}
