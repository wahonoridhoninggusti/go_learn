package main

import (
	"errors"
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
	fmt.Stringer
}

type Rectangle struct {
	Width, Height float64
}

func NewRectangle(width, height float64) (*Rectangle, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("rectangle dimensions must be positive")
	}
	return &Rectangle{width, height}, nil
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Height + r.Width)
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle [Width=%.2f, Height=%.2f]", r.Width, r.Height)
}

type Circle struct {
	Radius float64
}

func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
		return nil, errors.New("circle dimensions must be positive")
	}
	return &Circle{radius}, nil
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return math.Pi * (2 * c.Radius)
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle [Radius=%.2f]", c.Radius)
}

type Triangle struct {
	SideA, SideB, SideC float64
}

func NewTriangle(sideA, sideB, sideC float64) (*Triangle, error) {
	if sideA <= 0 || sideB <= 0 || sideC <= 0 {
		return nil, errors.New("Triangle dimensions must be positive")
	}

	if sideA+sideB == sideC || sideA+sideB < sideC {
		return nil, errors.New("invalid triangle")
	}
	return &Triangle{sideA, sideB, sideC}, nil
}

func (t Triangle) Area() float64 {
	half := (t.SideA + t.SideB + t.SideC) / 2
	return math.Sqrt(half * (half - t.SideA) * (half - t.SideB) * (half - t.SideC))
}

func (t Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

func (t Triangle) String() string {
	return fmt.Sprintf("Triangle sides [Side A=%.2f, Side B=%.2f, Side C=%.2f]", t.SideA, t.SideB, t.SideC)
}

type ShapeCalculator struct{}

func (cal *ShapeCalculator) PrintProperties(s Shape) {
	fmt.Println(s)
	fmt.Printf("  Area: %.2f\n", s.Area())
	fmt.Printf("  Perimeter: %.2f\n", s.Perimeter())
}

func (cal *ShapeCalculator) TotalArea(shape []Shape) float64 {
	total := 0.0
	for _, s := range shape {
		total += s.Area()
	}
	return total
}

func (cal *ShapeCalculator) LargestShape(shape []Shape) Shape {
	largest := shape[0]
	for _, s := range shape[1:] {
		if s.Area() > largest.Area() {
			largest = s
		}
	}

	return largest
}

func (cal *ShapeCalculator) SortByArea(shape []Shape, ascending bool) []Shape {
	sorted := make([]Shape, len(shape))
	copy(sorted, shape)
	n := len(sorted)

	for i := 1; i < n; i++ {
		key := sorted[i]
		j := i - 1

		if ascending {
			for j >= 0 && sorted[j].Area() > key.Area() {
				sorted[j+1] = sorted[j]
				j = j - 1
			}
			sorted[j+1] = key
		} else {
			for j >= 0 && sorted[j].Area() < key.Area() {
				sorted[j+1] = sorted[j]
				j--
			}
			sorted[j+1] = key
		}

	}
	return sorted
}

func NewShapeCalculator() *ShapeCalculator {
	return &ShapeCalculator{}
}

func main() {
	rect, _ := NewRectangle(4, 2)
	circle, _ := NewCircle(3)
	triangle, _ := NewTriangle(2, 3, 4)
	shapes := []Shape{rect, circle, triangle}
	sc := NewShapeCalculator()

	totalArea := sc.TotalArea(shapes)
	largestArea := sc.LargestShape(shapes)
	sortedArea := sc.SortByArea(shapes, true)
	fmt.Println("total area ", totalArea, "largest: ", largestArea)

	fmt.Println("->", sortedArea)
	for _, s := range shapes {
		sc.PrintProperties(s)
		fmt.Println()
	}
}
