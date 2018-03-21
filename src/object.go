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
	shape            uint8
	radius           float64
	center0, center1 vec3
	time0, time1     float64
	mat              *material
}

func sphere(radius float64, center vec3, mat *material) *object {
	return &object{
		shapeCircle, radius, center, center, 0.0, 1.0, mat,
	}
}

func movingSphere(radius float64, center0, center1 vec3, time0, time1 float64, mat *material) *object {
	return &object{
		shapeCircle, radius, center0, center1, time0, time1, mat,
	}
}

func (o *object) hit(r ray, tmin float64, tmax float64, hr *hitRecord) bool {
	// Different implementations for different shapes.
	switch o.shape {
		// In case it's a circle, the only one we have right now.
	case shapeCircle:
		// Variables that are necessary for the ABC formula.
		oc := r.origin.sub(o.center(r.time))
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
				hr.p = r.point(hr.t)
				hr.normal = hr.p.sub(o.center(r.time)).divScalar(o.radius)
				hr.mat = o.mat

				return true
			}

			// If there was no solution we now try it with the plus variant.
			temp = (-b + math.Sqrt(discriminant)) / a
			if temp > tmin && temp < tmax {
				hr.t = temp
				hr.p = r.point(hr.t)
				hr.normal = hr.p.sub(o.center(r.time)).divScalar(o.radius)
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

func (o *object) center(t float64) vec3 {
	return o.center0.add(o.center1.sub(o.center0).mulScalar((t - o.time0) / (o.time1 - o.time0)))
}

// Create the bounding box for an object.
func (o *object) boundingBox(t0, t1 float64, box *aabb) bool {
	// Make the box for the begin and end center.
	box0 := &aabb{o.center(o.time0).subScalar(o.radius), o.center(o.time0).addScalar(o.radius)}
	box1 := &aabb{o.center(o.time1).subScalar(o.radius), o.center(o.time1).addScalar(o.radius)}

	// Combine the two boxes.
	box = surroundingBox(box0, box1)
	return true
}

// This is used for moving objects, the bounding box will be the entire path.
func surroundingBox(b0, b1 *aabb) *aabb {
	small := vec(ffmin(b0.min.x, b1.min.x),
		ffmin(b0.min.y, b1.min.y),
		ffmin(b0.min.z, b0.min.z))
	big := vec(ffmax(b0.max.x, b1.max.x),
		ffmax(b0.max.y, b1.max.y),
		ffmax(b0.max.z, b0.max.z))
	return &aabb{small, big}
}

type aabb struct {
	min, max vec3
}

func (b *aabb) hit(r ray, tmin, tmax float64) bool {
	// TODO: Check if this needs to be replaced by the other method given.
	for a := 0; a < 3; a++ {
		t0 := ffmin((b.min.get(a)-r.origin.get(a))/r.dir.get(a),
			(b.max.get(a)-r.origin.get(a))/r.dir.get(a))
		t1 := ffmax((b.min.get(a)-r.origin.get(a))/r.dir.get(a),
			(b.max.get(a)-r.origin.get(a))/r.dir.get(a))

		tmin = ffmax(t0, tmin)
		tmax = ffmin(t1, tmax)
		if tmax <= tmin {
			return false
		}
	}
	return true
}

// TODO: Do we really need a function for this?
func ffmin(a, b float64) float64 {
	if a > b {
		return b
	}
	return a
}

// TODO: Do we really need a function for this?
func ffmax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

type bvhNode struct {
	left, right object
	box         aabb
}

// TODO: Finish this system!
// TODO: Finish this system!
// TODO: Finish this system!
func (b *bvhNode) hit(r ray, tmin, tmax float64, hr *hitRecord) bool {
	// Check if the bounding box has been hit, if it doesn't hit the box,
	// it will definitely not hit the object.
	if b.box.hit(r, tmin, tmax) {
		// We've hit the box, to get more information about the actual objects, we create to hitrecords.
		var lhr, rhr *hitRecord
		// Check which specific object we hit.
		hitLeft := b.left.hit(r, tmin, tmax, lhr)
		hitRight := b.right.hit(r, tmin, tmax, rhr)
		if hitLeft && hitRight {
			// If we hit both we need to figure out which one is in the front.

		} else if hitLeft {
			// We only hit the left one.
			hr = lhr
			return true
		} else if hitRight {
			// We only hit the right one.
			hr = rhr
			return true
		}
		// We hit nothing this can happen,
		// because bounding boxes are not perfectly aligned with objects.
		return false
	}

	return false
}
