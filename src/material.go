package main

import (
	"math"
	"math/rand"
)

type material struct {
	matType  uint8
	col      vec3
	fuzz     float64
	refIndex float64
}

const (
	matDiffuse = 0
	matMetal   = 1
	matGlass   = 2
)

func dif(r, g, b float64) *material {
	m := material{}

	m.matType = matDiffuse
	m.col = v(r, g, b)

	return &m
}

func met(r, g, b, f float64) *material {
	m := material{}

	m.matType = matMetal
	m.col = v(r, g, b)
	m.fuzz = f

	return &m
}

func glass(index float64) *material {
	m := material{}

	m.matType = matGlass
	m.refIndex = index

	return &m
}

func (m *material) scatter(rIn ray, hr *hitRecord, atten *vec3, rOut *ray, rnd *rand.Rand) bool {
	// Difference between diffuse and metallic materials.
	switch m.matType {
	case matDiffuse:
		target := hr.p.add(hr.normal).add(randInUnitSphere(rnd))
		*rOut = ray{hr.p, target.sub(hr.p)}
		*atten = m.col
		return true

	case matMetal:
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

	case matGlass:
		reflected := reflect(rIn.dir, hr.normal)
		*atten = v(1.0, 1.0, 1.0)
		refracted := vec3{}
		outwardNormal := vec3{}

		var niOverNt float64
		var cosine float64
		var reflectProbe float64

		if dot(rIn.dir, hr.normal) > 0.0 {
			outwardNormal = hr.normal.mulScalar(-1.0)
			niOverNt = m.refIndex
			cosine = m.refIndex * dot(rIn.dir, hr.normal) / rIn.dir.length()
		} else {
			outwardNormal = hr.normal
			niOverNt = 1.0 / m.refIndex
			cosine = -dot(rIn.dir, hr.normal) / rIn.dir.length()
		}

		if refract(rIn.dir, outwardNormal, niOverNt, &refracted) {
			reflectProbe = schlick(cosine, m.refIndex)
		} else {
			*rOut = ray{hr.p, reflected}
			return true
		}

		if rnd.Float64() < reflectProbe {
			*rOut = ray{hr.p, reflected}
		} else {
			*rOut = ray{hr.p, refracted}
		}

		return true
	}

	return false
}

func schlick(cosine, index float64) float64 {
	r0 := (1 - index) / (1 + index)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

func reflect(v vec3, n vec3) vec3 {
	return v.sub(n.mulScalar(2 * dot(v, n)))
}

func refract(v, n vec3, rf float64, vOut *vec3) bool {
	nv := v.normalize()
	dt := dot(nv, n)
	discriminant := 1.0 - rf*rf*(1-dt*dt)
	if discriminant > 0.0 {
		*vOut = nv.sub(n.mulScalar(dt)).mulScalar(rf).sub(n.mulScalar(math.Sqrt(discriminant)))
		return true
	}
	return false
}
