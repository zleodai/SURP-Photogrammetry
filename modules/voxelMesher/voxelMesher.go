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

type VoxelBB struct {
	minX float64
	maxX float64
	minY float64
	maxY float64
	minZ float64
	maxZ float64
}

func generateVoxels(xSize, ySize, zSize int, minX, minY, minZ, voxelSize float64) []VoxelBB {
	voxels := []VoxelBB{}
	for xIndex := 0; xIndex < xSize; xIndex++ {
		for yIndex := 0; yIndex < ySize; yIndex++ {
			for zIndex := 0; zIndex < zSize; zIndex++ {
				var upperX float64 = minX + (float64(xIndex+1) * voxelSize)
				var upperY float64 = minY + (float64(yIndex+1) * voxelSize)
				var upperZ float64 = minZ + (float64(zIndex+1) * voxelSize)
				var lowerX float64 = minX + (float64(xIndex) * voxelSize)
				var lowerY float64 = minY + (float64(yIndex) * voxelSize)
				var lowerZ float64 = minZ + (float64(zIndex) * voxelSize)
				newVoxel := VoxelBB{
					minX: lowerX,
					maxX: upperX,
					minY: lowerY,
					maxY: upperY,
					minZ: lowerZ,
					maxZ: upperZ,
				}
				voxels = append(voxels, newVoxel)
			}
		}
	}

	return voxels
}

func checkPointInside(voxel VoxelBB, point pointCloudDecoder.Point) bool {
	if point.X >= voxel.minX && point.X <= voxel.maxX && point.Y >= voxel.minY && point.Y <= voxel.maxY && point.Z >= voxel.minZ && point.Z <= voxel.maxZ {
		return true
	}
	return false
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
