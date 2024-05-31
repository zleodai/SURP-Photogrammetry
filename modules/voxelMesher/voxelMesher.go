package voxelMesher

import (
	"encoding/json"
	"fmt"
	"math"
	"modules/pointCloudDecoder"
	"os"
	"runtime"
)

type pointVal struct {
	X     int
	Y     int
	Z     int
	Value int
}

type pointList struct {
	Points []pointVal
}

func Test() byte {
	return 0
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Mesh(xArray, yArray, zArray []pointCloudDecoder.Point, voxelSize float64) [][][]int {
	var xSize int = int(math.Floor(math.Abs(xArray[0].X-xArray[len(xArray)-1].X) / voxelSize))
	var ySize int = int(math.Floor(math.Abs(yArray[0].Y-yArray[len(yArray)-1].Y) / voxelSize))
	var zSize int = int(math.Floor(math.Abs(zArray[0].Z-zArray[len(zArray)-1].Z) / voxelSize))
	fmt.Println("Voxel Grid X, Y, Z:", xSize, ySize, zSize)
	fmt.Println("Total Voxel Amount:", xSize*ySize*zSize)

	voxels := make([][][]int, xSize+1)
	for i := 0; i < len(voxels); i++ {
		voxels[i] = make([][]int, ySize+1)
		for j := 0; j < len(voxels[i]); j++ {
			voxels[i][j] = make([]int, zSize+1)
		}
	}

	for _, point := range xArray {
		xIndex := int(math.Floor(math.Abs(xArray[0].X-point.X) / voxelSize))
		yIndex := int(math.Floor(math.Abs(yArray[0].Y-point.Y) / voxelSize))
		zIndex := int(math.Floor(math.Abs(zArray[0].Z-point.Z) / voxelSize))
		voxels[xIndex][yIndex][zIndex] += 1
	}

	//GenerateVoxelJson(voxels)

	runtime.GC()
	return voxels
}

func MinMaxMesh(xMinMax, yMinMax, zMinMax [2]float64, points []pointCloudDecoder.Point, voxelSize float64) [][][]int {
	var xSize int = int(math.Floor(math.Abs(xMinMax[0]-xMinMax[1]) / voxelSize))
	var ySize int = int(math.Floor(math.Abs(yMinMax[0]-yMinMax[1]) / voxelSize))
	var zSize int = int(math.Floor(math.Abs(zMinMax[0]-zMinMax[1]) / voxelSize))
	fmt.Println("Voxel Grid X, Y, Z:", xSize, ySize, zSize)
	fmt.Println("Total Voxel Amount:", xSize*ySize*zSize)

	voxels := make([][][]int, xSize+1)
	for i := 0; i < len(voxels); i++ {
		voxels[i] = make([][]int, ySize+1)
		for j := 0; j < len(voxels[i]); j++ {
			voxels[i][j] = make([]int, zSize+1)
		}
	}

	for _, point := range points {
		xIndex := int(math.Floor(math.Abs(xMinMax[0]-point.X) / voxelSize))
		yIndex := int(math.Floor(math.Abs(yMinMax[0]-point.Y) / voxelSize))
		zIndex := int(math.Floor(math.Abs(zMinMax[0]-point.Z) / voxelSize))
		voxels[xIndex][yIndex][zIndex] += 1
	}

	//GenerateVoxelJson(voxels)

	runtime.GC()
	return voxels
}



func IterativeMesh(xMinMax, yMinMax, zMinMax [2]float64, points []pointCloudDecoder.Point, voxelEndSize float64, Iterations, scaleFactor int) {
	var xSize int = int(math.Floor(math.Abs(xMinMax[0]-xMinMax[1]) / voxelEndSize))
	var ySize int = int(math.Floor(math.Abs(yMinMax[0]-yMinMax[1]) / voxelEndSize))
	var zSize int = int(math.Floor(math.Abs(zMinMax[0]-zMinMax[1]) / voxelEndSize))

	masterVoxels := make([][][]int, xSize+1)
	for i := 0; i < len(masterVoxels); i++ {
		masterVoxels[i] = make([][]int, ySize+1)
		for j := 0; j < len(masterVoxels[i]); j++ {
			masterVoxels[i][j] = make([]int, zSize+1)
		}
	}
	
	voxelSize := voxelEndSize * math.Pow(float64(scaleFactor), float64(Iterations))
	for i := 0; i < Iterations; i++ {
		currVoxels := MinMaxMesh(xMinMax, yMinMax, zMinMax, points, voxelSize)

		numSmallVoxels := math.Pow(float64(scaleFactor), float64(Iterations - i))

		for x, yArray := range masterVoxels {
			for y, zArray := range yArray {
				for z := 0; z < len(zArray); z++ {
					masterVoxels[x][y][z] += currVoxels[int(x / int(numSmallVoxels))][int(y / int(numSmallVoxels))][int(z/int(numSmallVoxels))]
				}
			}
		}
		voxelSize /= float64(scaleFactor)
	}
	GenerateVoxelJson(masterVoxels)
}

func GenerateVoxelJson(voxels [][][]int) {
	file, errs := os.Create("VoxelMatrix.JSON")
	if errs != nil {
		panic("Failed to write to file:" + errs.Error())
	}

	enc := json.NewEncoder(file)
	cleanedPoints := []pointVal{}

	for xIndex, xArray := range voxels {
		for yIndex, yArray := range xArray {
			for zIndex, Value := range yArray {
				if Value != 0 {
					point := pointVal{X: xIndex, Y: yIndex, Z: zIndex, Value: Value}
					cleanedPoints = append(cleanedPoints, point)
				}
			}
		}
	}

	enc.Encode(pointList{Points: cleanedPoints})
}
