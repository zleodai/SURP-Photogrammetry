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

func Mesh(xArray, yArray, zArray []pointCloudDecoder.Point, voxelSize float64) {
	var xSize int = int(math.Floor(math.Abs(xArray[0].X-xArray[len(xArray)-1].X) / voxelSize))
	var ySize int = int(math.Floor(math.Abs(yArray[0].Y-yArray[len(yArray)-1].Y) / voxelSize))
	var zSize int = int(math.Floor(math.Abs(zArray[0].Z-zArray[len(zArray)-1].Z) / voxelSize))
	fmt.Println("Voxel Grid X, Y, Z:", xSize, ySize, zSize)
	fmt.Println("Total Voxel Amount:", xSize*ySize*zSize)

	voxels := make([][][]int, xSize+1)
	for yIndex, _ := range voxels {
		voxels[yIndex] = make([][]int, ySize+1)
		for zIndex, _ := range voxels[yIndex] {
			voxels[yIndex][zIndex] = make([]int, zSize+1)
		}
	}

	for _, point := range xArray {
		xIndex := int(math.Floor(math.Abs(xArray[0].X-point.X) / voxelSize))
		yIndex := int(math.Floor(math.Abs(yArray[0].Y-point.Y) / voxelSize))
		zIndex := int(math.Floor(math.Abs(zArray[0].Z-point.Z) / voxelSize))
		voxels[xIndex][yIndex][zIndex] += 1
	}

	//GenerateVoxelJson(voxels, voxelSize)

	runtime.GC()
}

func GenerateVoxelJson(voxels [][][]int, voxelSize float64) {
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
