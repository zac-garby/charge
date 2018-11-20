package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

const (
	k          = 10000
	width      = 256
	height     = 256
	depth      = 256
	gap        = 1.0
	numCharges = 8
	numWorkers = 4
)

type vec struct {
	x, y, z float64
}

func (v *vec) add(o *vec) *vec {
	return &vec{x: v.x + o.x, y: v.y + o.y, z: v.z + o.z}
}

func (v *vec) sub(o *vec) *vec {
	return &vec{x: v.x - o.x, y: v.y - o.y, z: v.z - o.z}
}

func (v *vec) scale(s float64) *vec {
	return &vec{x: v.x * s, y: v.y * s, z: v.z * s}
}

func (v *vec) magSqr() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

func (v *vec) mag() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v *vec) norm() *vec {
	mag := v.mag()
	if mag == 0 {
		return v
	}
	return v.scale(1 / mag)
}

type charge struct {
	pos    *vec
	charge float64
}

type charges []*charge

func (c charges) field(pos *vec) *vec {
	r := new(vec)

	for _, charge := range c {
		diff := charge.pos.sub(pos)
		distSqr := diff.magSqr()
		if distSqr == 0 {
			continue
		}

		force := (k * charge.charge) / distSqr
		norm := diff.norm()
		r = r.add(norm.scale(math.Copysign(math.Log(math.Abs(force)), force)))
	}

	return r
}

type slice struct {
	z float64
	i int
}

func (c charges) worker(id int, slices []slice, finished chan bool) {
	fmt.Println("worker", id, "initialised.", len(slices), "slices in the batch")

	img := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{width, height},
	})

	for _, s := range slices {
		for x := 0.0; x < width; x += gap {
			for y := 0.0; y < height; y += gap {
				field := c.field(&vec{x, y, s.z})

				var cl color.RGBA

				switch int(field.mag()) % 2 {
				case 0:
					cl = color.RGBA{0xff, 0x00, 0x00, 255}
				case 1:
					cl = color.RGBA{0xff, 0xff, 0xff, 255}
				}

				for a := 0; a < gap; a++ {
					for b := 0; b < gap; b++ {
						img.Set(int(x)+a, int(y)+b, cl)
					}
				}
			}
		}

		f, _ := os.Create(fmt.Sprintf("out/%d.png", s.i))
		png.Encode(f, img)
		fmt.Println(id, "did", s.z)
	}

	finished <- true
}

func main() {
	rand.Seed(time.Now().UnixNano())

	c := charges{}

	for i := 0; i < numCharges; i++ {
		c = append(c, &charge{
			pos: &vec{
				x: rand.Float64() * width,
				y: rand.Float64() * height,
				z: rand.Float64() * depth,
			},
			charge: rand.Float64()*2000 - 1000,
		})
	}

	batches := [][]slice{}

	for n := 0; n < numWorkers; n++ {
		batches = append(batches, []slice{})
	}

	i := 0
	for z := 0.0; z < depth; z += gap {
		batches[i%numWorkers] = append(batches[i%numWorkers], slice{
			i: i + 1,
			z: z,
		})
		i++
	}

	finished := make(chan bool)

	for id, batch := range batches {
		go c.worker(id, batch, finished)
	}

	for i := 0; i < numWorkers; i++ {
		<-finished
	}
}
