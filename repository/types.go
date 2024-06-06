// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type CreateEstateInput struct {
	Id     string
	Width  int
	Length int
}

type GetEstateByIdInput struct {
	Id string
}

type GetEstateByIdOutput struct {
	Id            string
	Width         int
	Length        int
	Count         int
	Max           int
	Min           int
	Median        float64
	DroneDistance int
}

type CountCoordinateTreeInput struct {
	X int
	Y int
}

type CountCoordinateTreeOutput struct {
	Count int
}

type CreateTreeInput struct {
	Id     string
	X      int
	Y      int
	Height int

	EstateId        string
	DroneDistFactor int
}

type GetPrevNextTreeInput struct {
	PrevX int
	PrevY int

	NextX int
	NextY int
}

type GetPrevNextTreeOutput struct {
	PrevTreeHeight int
	NextTreeHeight int
}

type GetHeightEstateTreesInput struct {
	EstateId string
}

type GetHeightEstateTreesOutput struct {
	Heights []int
}

type StoreMedianEstateInput struct {
	EstateId string
	Median   float64
}

type GetAllEstateTreesInput struct {
	EstateId string

	Orders []string
	Limit  int
}

type EstateTreeOutput struct {
	X      int
	Y      int
	Height int
}

type GetAllEstateTreesOutput struct {
	EstateTrees []EstateTreeOutput
}
