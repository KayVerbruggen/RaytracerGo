package main

import "math/rand"

type material struct {
	diffuse bool
	col     vec3
	fuzz    float64
}

func (m material) scatter(rIn ray, hr *hitRecord, atten *vec3, rOut *ray, rnd *rand.Rand) bool {
	// Difference between diffuse and metallic materials.
	if m.diffuse {
		target := hr.p.add(hr.normal).add(randInUnitSphere(rnd))
		*rOut = ray{hr.p, target.sub(hr.p)}
		*atten = m.col
		return true
	}

	reflected := reflect(rIn.dir, hr.normal)

	// For optimization, there is no point in calculating random in unit sphere,
	// if it's going to be multiplied be 0 anyway. Improvement: 33%	for a material using 0.0 fuzz.
	if m.fuzz == 0.0 {
		*rOut = ray{hr.p, reflected}
	} else {
		*rOut = ray{hr.p, reflected.add(randInUnitSphere(rnd).subScalar(m.fuzz))}
	}
	*atten = m.col

	return dot(rOut.dir, hr.normal) > 0.0
}

func dif(r, g, b float64) material {
	return material{true, v(r, g, b), 0.0}
}

func met(r, g, b, f float64) material {
	return material{false, v(r, g, b), f}
}
