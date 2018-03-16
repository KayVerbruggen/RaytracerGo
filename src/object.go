package main

import (
	"math"
)

type hitRecord struct {
	t         float64
	p, normal vec3
	mat       *material
}

// Use these to differentiate the different shapes of objects.
const (
	shapeCircle uint8 = 0
)

// Objects can be hit by rays.
type object struct {
	shape  uint8
	radius float64
	center vec3
	mat    *material
}

func (o *object) hit(r ray, tmin float64, tmax float64, hr *hitRecord) bool {
	// Different implementations for different shapes.
	switch o.shape {
	// In case it's a circle, the only one we have right now.
	case shapeCircle:
		// Variables that are necessary for the ABC formula.
		oc := r.origin.sub(o.center)
		a := dot(r.dir, r.dir)
		b := dot(oc, r.dir)
		c := dot(oc, oc) - o.radius*o.radius

		// Use the ABC formula to figure out if we're hitting the sphere.
		discriminant := b*b - a*c
		if discriminant > 0.0 {
			// In case I'll ever wonder why we first try the minus variant,
			// it's because you want the front side of the sphere.
			temp := (-b - math.Sqrt(discriminant)) / a
			if temp > tmin && temp < tmax {
				hr.t = temp
				hr.p = r.pointAtParam(hr.t)
				hr.normal = hr.p.sub(o.center).divScalar(o.radius)
				hr.mat = o.mat

				return true
			}

			// If there was no solution we now try it with the plus variant.
			temp = (-b + math.Sqrt(discriminant)) / a
			if temp > tmin && temp < tmax {
				hr.t = temp
				hr.p = r.pointAtParam(hr.t)
				hr.normal = hr.p.sub(o.center).divScalar(o.radius)
				hr.mat = o.mat

				return true
			}
		}

		// Return false, because there is no solution.
		return false

	// This should never happen, but whatever.
	default:
		return false
	}
}
