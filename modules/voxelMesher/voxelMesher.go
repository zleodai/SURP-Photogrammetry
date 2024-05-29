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

type Voxel struct {
	coords [4]float64
	center int
}

func generateVoxels () {

}

func checkPointInside () bool {

}

func Mesh(xArray, yArray, zArray []pointCloudDecoder.Point, voxelSize float64) {
	var xSize int = int(math.Floor(math.Abs(xArray[0].X-xArray[len(xArray)-1].X) / voxelSize))
	var ySize int = int(math.Floor(math.Abs(yArray[0].Y-yArray[len(yArray)-1].Y) / voxelSize))
	var zSize int = int(math.Floor(math.Abs(zArray[0].Z-zArray[len(zArray)-1].Z) / voxelSize))
	fmt.Println("Voxel Grid X, Y, Z:", xSize, ySize, zSize)
	fmt.Println("Total Voxel Amount:", xSize*ySize*zSize)

	
	runtime.GC()
}
