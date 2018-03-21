package main

import (
	"math"
	"math/rand"
)

type material struct {
	matType  uint8
	tex      texture
	fuzz     float64
	refIndex float64
}

const (
	matDiffuse = 0
	matMetal   = 1
	matGlass   = 2
)

func dif(tex texture) *material {
	m := material{}

	m.matType = matDiffuse
	m.tex = tex

	return &m
}

func met(tex texture, f float64) *material {
	m := material{}

	m.matType = matMetal
	m.tex = tex
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
		*rOut = ray{hr.p, target.sub(hr.p), rIn.time}
		*atten = m.tex.value(0.0, 0.0, hr.p)
		return true

	case matMetal:
		reflected := reflect(rIn.dir, hr.normal)

		// For optimization, there is no point in calculating random in unit sphere,
		// if it's going to be multiplied be 0 anyway. Improvement: 33%	for a material using 0.0 fuzz.
		if m.fuzz == 0.0 {
			*rOut = ray{hr.p, reflected, rIn.time}
		} else {
			*rOut = ray{hr.p, reflected.add(randInUnitSphere(rnd).subScalar(m.fuzz)), rIn.time}
		}
		*atten = m.tex.value(0.0, 0.0, hr.p)

		return dot(rOut.dir, hr.normal) > 0.0

	case matGlass:
		reflected := reflect(rIn.dir, hr.normal)
		*atten = vec(1.0, 1.0, 1.0)
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
			*rOut = ray{hr.p, reflected, rIn.time}
			return true
		}

		if rnd.Float64() < reflectProbe {
			*rOut = ray{hr.p, reflected, rIn.time}
		} else {
			*rOut = ray{hr.p, refracted, rIn.time}
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

type texture interface {
	value(u, v float64, p vec3) vec3
}

type solidColor struct {
	col vec3
}

func (c solidColor) value(u, v float64, p vec3) vec3 {
	return c.col
}

func col(r, g, b float64) texture {
	return solidColor{vec(r, g, b)}
}

type checkerTex struct {
	colOdd, colEven vec3
}

func checker(col0, col1 vec3) texture {
	return checkerTex{col0, col1}
}

func (c checkerTex) value(u, v float64, p vec3) vec3 {
	sines := math.Sin(10.0*p.x) * math.Sin(10.0*p.y) * math.Sin(10.0*p.z)
	if sines < 0 {
		return c.colOdd
	}

	return c.colEven
}

type noiseTex struct {
	noise *perlin
	scale float64
}

func perlTex(s float64) *noiseTex {
	return &noiseTex{per(), s}
}

func (t *noiseTex) value(u, v float64, p vec3) vec3 {
	return vec(1.0, 1.0, 1.0).mulScalar(0.5).mulScalar(1.0 + math.Sin(t.scale * p.z + 10*t.noise.turb(p)))
}

// Perlin noise is blurred white noise.
type perlin struct {
	permX, permY, permZ [256]int32
	ranvec              [256]vec3
}

func per() *perlin {
	return &perlin{
		perlinGenPerm(),
		perlinGenPerm(),
		perlinGenPerm(),
		perlinGen(),
	}
}

func (p *perlin) noise(perm vec3) float64 {
	u := perm.x - math.Floor(perm.x)
	v := perm.y - math.Floor(perm.y)
	w := perm.z - math.Floor(perm.z)

	i := int(math.Floor(perm.x))
	j := int(math.Floor(perm.y))
	k := int(math.Floor(perm.z))

	var c [2][2][2]vec3
	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				c[di][dj][dk] = p.ranvec[p.permX[(i+di) & 255] ^ p.permY[(j+dj) & 255] ^ p.permZ[(k+dk) & 255]]
			}
		}
	}

	return trilinearInterp(c, u, v, w)
}

func (p *perlin) turb(perm vec3) float64 {
	accum := 0.0
	temp := perm
	weight := 1.0
	for i := 0; i < 7; i++ {
		accum += weight*p.noise(temp)
		weight *= 0.5
		temp = temp.mulScalar(2.0)
	}

	return math.Abs(accum)
}

func trilinearInterp(c [2][2][2]vec3, u, v, w float64) float64 {
	uu := u * u * (3.0 - 2.0*u)
	vv := v * v * (3.0 - 2.0*v)
	ww := w * w * (3.0 - 2.0*w)

	accum := 0.0
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				weight := vec(u-float64(i), v-float64(j), w-float64(k))
				accum += (float64(i)*uu + (1.0-float64(i))*(1.0-uu)) *
					(float64(j)*vv + (1.0-float64(j))*(1.0-vv)) *
					(float64(k)*ww + (1.0-float64(k))*(1.0-ww)) * dot(c[i][j][k], weight)
			}
		}
	}

	return accum
}

func permute(p *[256]int32, n int32) {
	for i := n - 1; i > 0; i-- {
		target := int32(rand.Float64() * float64(i+1))
		tmp := p[i]
		p[i] = p[target]
		p[target] = tmp
	}
}

func perlinGenPerm() [256]int32 {
	var p [256]int32
	for i := 0; i < 256; i++ {
		p[i] = int32(i)
	}
	permute(&p, 256)
	return p
}

// Generate the perlin noise.
func perlinGen() [256]vec3 {
	var p [256]vec3

	for i := 0; i < 256; i++ {
		p[i] = vec(-1.0+2.0*rand.Float64(), -1.0+2.0*rand.Float64(), -1.0+2.0*rand.Float64()).normalize()
	}

	return p
}
