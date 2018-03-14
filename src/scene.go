package main

type scene struct {
	cam     camera
	objects []object
}

func (s scene) hit(r ray, tmin float64, tmax float64, hr *hitRecord) bool {
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
