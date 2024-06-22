package pointSorter

import (
	"container/heap"
	"modules/pointCloudDecoder"
)

type HeapItem struct {
	Value    [4]int
	Priority float64
	Index    int
}

type PriorityQueue []*HeapItem

func Test() byte {
	return 0
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*HeapItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *HeapItem, value [4]int, priority float64) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.Index)
}

func MinMaxPoints(pointData pointCloudDecoder.PointData) ([2]float64, [2]float64, [2]float64) {
	var points []pointCloudDecoder.Point = pointData.Points

	xMinMax := [2]float64{0, 0}
	yMinMax := [2]float64{0, 0}
	zMinMax := [2]float64{0, 0}

	for _, point := range points {
		if point.X < xMinMax[0] {
			xMinMax[0] = point.X
		}
		if point.X > xMinMax[1] {
			xMinMax[1] = point.X
		}
		if point.Y < yMinMax[0] {
			yMinMax[0] = point.Y
		}
		if point.Y > yMinMax[1] {
			yMinMax[1] = point.Y
		}
		if point.Z < zMinMax[0] {
			zMinMax[0] = point.Z
		}
		if point.Z > zMinMax[1] {
			zMinMax[1] = point.Z
		}
	}
	return xMinMax, yMinMax, zMinMax
}

func SortVertices(verticies []HeapItem) ([][3]int, []int) {
	vertexHeap := make(PriorityQueue, len(verticies))

	for _, vertex := range verticies {
		vertexHeap[vertex.Index] = &vertex
	}
	
	heap.Init(&vertexHeap)

	heapArray := [][3]int{}
	indexArray := []int{}

	for index := 0; index < len(verticies); index++ {
		heapItem := heap.Pop(&vertexHeap).(*HeapItem)
		vertexValue := [3]int{heapItem.Value[0], heapItem.Value[1], heapItem.Value[2]}
		heapArray = append(heapArray, vertexValue)
		indexArray = append(indexArray, heapItem.Value[3])
	}

	return heapArray,indexArray
}

// func SortPointData(pointData pointCloudDecoder.PointData) ([]pointCloudDecoder.Point, []pointCloudDecoder.Point, []pointCloudDecoder.Point) {
// 	var points []pointCloudDecoder.Point = pointData.Points

// 	xPQ := make(PriorityQueue, len(pointData.Points))
// 	yPQ := make(PriorityQueue, len(pointData.Points))
// 	zPQ := make(PriorityQueue, len(pointData.Points))

// 	for index, point := range points {
// 		xPQ[index] = &HeapItem{
// 			value:    point,
// 			priority: point.X,
// 			index:    index,
// 		}
// 		yPQ[index] = &HeapItem{
// 			value:    point,
// 			priority: point.Y,
// 			index:    index,
// 		}
// 		zPQ[index] = &HeapItem{
// 			value:    point,
// 			priority: point.Z,
// 			index:    index,
// 		}
// 	}

// 	heap.Init(&xPQ)
// 	heap.Init(&yPQ)
// 	heap.Init(&zPQ)

// 	xArray := []pointCloudDecoder.Point{}
// 	yArray := []pointCloudDecoder.Point{}
// 	zArray := []pointCloudDecoder.Point{}

// 	for index := 0; index < len(points); index++ {
// 		xItem := heap.Pop(&xPQ).(*HeapItem)
// 		yItem := heap.Pop(&yPQ).(*HeapItem)
// 		zItem := heap.Pop(&zPQ).(*HeapItem)

// 		xArray = append(xArray, xItem.value)
// 		yArray = append(yArray, yItem.value)
// 		zArray = append(zArray, zItem.value)
// 	}
// 	return xArray, yArray, zArray
// }
