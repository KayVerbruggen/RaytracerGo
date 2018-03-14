package main

import "math"

// vec3 uses float64.
type vec3 struct {
	x, y, z float64
}

func v(x, y, z float64) vec3 {
	return vec3{x, y, z}
}

// Addvec3 is used for adding two vec3's
func (v vec3) add(v3 vec3) vec3 {
	return vec3{
		x: v.x + v3.x,
		y: v.y + v3.y,
		z: v.z + v3.z,
	}
}

// Subvec3 is used for subtracting two vec3's
func (v vec3) sub(v3 vec3) vec3 {
	return vec3{
		x: v.x - v3.x,
		y: v.y - v3.y,
		z: v.z - v3.z,
	}
}

// Mulvec3 is used for multiplying two vec3's
func (v vec3) mul(v3 vec3) vec3 {
	return vec3{
		x: v.x * v3.x,
		y: v.y * v3.y,
		z: v.z * v3.z,
	}
}

// Divvec3 is used for dividing two vec3's
func (v vec3) div(v3 vec3) vec3 {
	return vec3{
		x: v.x / v3.x,
		y: v.y / v3.y,
		z: v.z / v3.z,
	}
}

// AddFloat is used for adding a float
func (v vec3) addScalar(f float64) vec3 {
	return vec3{
		x: v.x + f,
		y: v.y + f,
		z: v.z + f,
	}
}

// SubFloat is used for subtracting a float
func (v vec3) subScalar(f float64) vec3 {
	return vec3{
		x: v.x - f,
		y: v.y - f,
		z: v.z - f,
	}
}

// MulFloat is used for multiplying by a float
func (v vec3) mulScalar(f float64) vec3 {
	return vec3{
		x: v.x * f,
		y: v.y * f,
		z: v.z * f,
	}
}

// DivFloat is used for dividing by a float
func (v vec3) divScalar(f float64) vec3 {
	return vec3{
		x: v.x / f,
		y: v.y / f,
		z: v.z / f,
	}
}

// Length will return the length of the vec3
func (v vec3) length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

// LengthSqr will turn the squared length of the vec3, this will be faster than Length().
func (v vec3) lengthSqr() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

// Normalize gives the vector a length of 1.
func (v vec3) normalize() vec3 {
	return v.divScalar(v.length())
}

// Dot does something with vectors.
func dot(v1 vec3, v2 vec3) float64 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

// Cross also does something with vectors.
func cross(v1 vec3, v2 vec3) vec3 {
	return vec3 {
		x: v1.y*v2.z - v1.z*v2.y,
		y: -(v1.x*v2.z - v1.z*v2.x),
		z: v1.x*v2.y - v1.y*v2.x,
	}
}

func reflect(v vec3, n vec3) vec3 {
	return v.sub(n.mulScalar(2*dot(v, n)))
}