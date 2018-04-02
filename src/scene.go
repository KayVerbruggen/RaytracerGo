package main

import (
	"math/rand"
)

// Scenes can be rendered, they contain a list of objects and a camera.
type scene struct {
	cam     *camera
	objects []*object
}

func (s *scene) hit(r ray, tmin float64, tmax float64, hr *hitRecord) bool {
	hitAny := false       // We have't hit anything.
	closestSoFar := tmax  // We haven't hit anything, so there's no closest.
	tempHr := hitRecord{} // Create an empty hit record.

	for i := 0; i < len(s.objects); i++ {
		if s.objects[i].hit(r, tmin, closestSoFar, &tempHr) {
			// We've hit something!
			hitAny = true
			// We want our objects to be closer than this one.
			closestSoFar = tempHr.t

			// This is currently the closest object so we update the hit record.
			*hr = tempHr
		}
	}

	return hitAny
}

func (s *scene) boundingBox(t0, t1 float64, box *aabb) bool {
	// There needs to be at least one object!
	if len(s.objects) < 1 {
		return false
	}
	// Check if we even hit the first one.
	var tempBox *aabb
	if !s.objects[0].boundingBox(t0, t1, tempBox) {
		return false
	}
	box = tempBox
	// Now create a bounding box for all the objects.
	for i := 1; i < len(s.objects); i++ {
		if s.objects[1].boundingBox(t0, t1, tempBox) {
			box = surroundingBox(box, tempBox)
		} else {
			return false
		}
	}

	return true
}

func randScene() *scene {
	checkerMat := dif(checker(vec(0.2, 0.3, 0.1), vec(0.9, 0.9, 0.9)))
	marbleMat := dif(perlTex(4.0))
	texMat := dif(createImageTex("../res/texture.png"))

	// List of objects.
	objList := []*object{
		sphere(1000.0, vec(0.0, -1000, 0.0), checkerMat),
		sphere(1.0, vec(0.0, 1.0, 0.0), marbleMat),
		sphere(1.0, vec(-4.0, 1.0, 0.0), met(col(0.7, 0.6, 0.5), 0.0)),
		sphere(1.0, vec(4.0, 1.0, 0.0), texMat),
	}

	for a := -2; a < 2; a++ {
		for b := -2; b < 2; b++ {
			chooseMat := rand.Float64()
			center := vec(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+0.9*rand.Float64())

			if center.sub(vec(4.0, 0.2, 0.0)).length() > 0.9 {
				if chooseMat < 0.6 { // Diffuse
					objList = append(objList, sphere(0.2, center, dif(col(rand.Float64()*rand.Float64(), rand.Float64()*rand.Float64(), rand.Float64()*rand.Float64()))))
				} else if chooseMat < 0.8 { // Metal
					objList = append(objList, sphere(0.2, center, met(col(0.5*(1+rand.Float64()), 0.5*(1+rand.Float64()), 0.5*(1+rand.Float64())), 0.5*rand.Float64())))
				} else if chooseMat < 0.9{ // Glass
					objList = append(objList, sphere(0.2, center, glass(1.5)))
				} else { // Marble
					objList = append(objList, sphere(0.2, center, marbleMat))
				}
			}
		}
	}

	// Create a scene, containing a camera and a list of objects to render.
	return &scene{cam(vec(13.0, 2.0, 3.0), vec(0.0, 0.0, 0.0), 20.0, 0.15, 1.0),
		objList}
}
