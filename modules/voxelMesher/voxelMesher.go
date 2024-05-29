package voxelMesher

import (
	"fmt"
	"math"
	"modules/pointCloudDecoder"
	"runtime"
)

func Test() byte {
	return 0
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
		xIndex := int(math.Floor(math.Abs(xArray[0].X - point.X) / voxelSize))
		yIndex := int(math.Floor(math.Abs(yArray[0].Y - point.Y) / voxelSize))
		zIndex := int(math.Floor(math.Abs(zArray[0].Z - point.Z) / voxelSize))
		voxels[xIndex][yIndex][zIndex] += 1
	}

	fmt.Println(voxels)

	runtime.GC()
}
