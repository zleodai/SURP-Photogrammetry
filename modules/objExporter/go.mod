module objExporter

require (
	github.com/emirpasic/gods v1.18.1
	modules/greedyMesher v1.0.0
)

replace modules/greedyMesher => ../greedyMesher

require modules/pointSorter v1.0.0
replace modules/pointSorter => ../pointSorter

go 1.22.3
