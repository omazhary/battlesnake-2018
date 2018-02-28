package main

import (
	"gopkg.in/oleiade/lane.v1"
	"math"
)

var movement_map = map[int]string{
	1: "up",
	2: "right",
	3: "down",
	4: "left",
}

var move_queue = lane.NewDeque()

func Strategize(world *MoveRequest) string {

	myHeadLocation := world.You.Head()
	myTailLocation := world.You.Body.Data[world.You.Length-1]
	nearTailLocations := GetValidAdjacentPoints(myTailLocation, world)
	foodLocation := FindFood(myHeadLocation, world)
	var path_map []int
	if world.You.Health < 50 || world.You.Length < 10 {
		path_map = ShortestPath(myHeadLocation, foodLocation, world)
		if path_map == nil {
			for _, possible_target_destination := range nearTailLocations {
				path_map = ShortestPath(myHeadLocation, possible_target_destination, world)
				if path_map != nil {
					break
				}
			}
		}
	} else {
		for _, possible_target_destination := range nearTailLocations {
			path_map = ShortestPath(myHeadLocation, possible_target_destination, world)
			if path_map != nil {
				break
			}
		}
	}
	return movement_map[path_map[0]]
}

func IsGoingToHitHimselfAtPoint(p Point, world *MoveRequest) bool {
	for _, bodyPoints := range world.You.Body.Data {
		if p.Equals(bodyPoints) {
			return true
		}
	}
	return false
}

func IsGoingToHitOthersAtPoint(p Point, world *MoveRequest) bool {
	for _, snake := range world.Snakes.Data {
		for _, bodyPoints := range snake.Body.Data {
			if p.Equals(bodyPoints) {
				return true
			}
		}

		for _, position_next_to_enemys_head := range GetAdjacentPoints(snake.Head(), world) {
			if position_next_to_enemys_head.Equals(p) && snake.Health > world.You.Health {
				return true
			}
		}
	}
	return false
}

func IsValidPointToMoveTo(p Point, world *MoveRequest) bool {
	if p.IsOutOfMapBounds(world) {
		return false
	}
	if IsGoingToHitHimselfAtPoint(p, world) {
		return false
	}
	if IsGoingToHitOthersAtPoint(p, world) {
		return false
	}
	return true
}

func GetNextPointBasedOnDirection(direction int, currentPoint Point) Point {
	nextPoint := currentPoint
	switch direction {
	case 1: //up
		nextPoint.Y = nextPoint.Y - 1
	case 2: //right
		nextPoint.X++
	case 3: // down
		nextPoint.Y = nextPoint.Y + 1
	case 4: // left
		nextPoint.X = nextPoint.X - 1
	default:
		nextPoint.X = nextPoint.X - 1
	}
	return nextPoint
}

func SimplePath(source Point, destination Point) {
	xdimension := source.X - destination.X
	ydimension := source.Y - destination.Y

	if xdimension < 0 {
		for i := xdimension; i < 0; i++ {
			//fmt.Println("right")
			move_queue.Prepend(2)
		}
	} else {
		for i := 0; i < xdimension; i++ {
			//fmt.Println("left")
			move_queue.Prepend(4)
		}
	}
	if ydimension < 0 {
		for i := ydimension; i < 0; i++ {
			//fmt.Println("down")
			move_queue.Prepend(3)
		}
	} else {
		for i := 0; i < ydimension; i++ {
			//fmt.Println("up")
			move_queue.Prepend(1)
		}
	}
}

func FindFood(location Point, world *MoveRequest) Point {
	closest_distance := math.MaxFloat64
	closest_food := world.Food.Data[0]
	for _, food_source := range world.Food.Data {
		dist := location.distance(food_source)
		if dist < closest_distance {
			closest_distance = dist
			closest_food = food_source
		}
	}
	return closest_food
}

// returns a slice with the directions towards the destination. If no path to destination, returns nil
func ShortestPath(source Point, destination Point, world *MoveRequest) []int {
	queue := lane.NewDeque()
	var parent Point
	visited := map[Point]bool{}
	plan_to_visit := map[Point]bool{}
	possible_directions := []int{1, 2, 3, 4}
	o_path := map[Point][]int{}

	queue.Prepend(source)

	for !queue.Empty() {
		parent, _ = queue.Pop().(Point)
		plan_to_visit[parent] = false
		if parent.Equals(destination) {
			break
		}

		for _, next_move := range possible_directions {
			next_position := GetNextPointBasedOnDirection(next_move, parent)
			if !IsValidPointToMoveTo(next_position, world) {
				continue
			}
			if visited[next_position] {
				continue
			}
			if !plan_to_visit[next_position] {
				queue.Prepend(next_position)
				plan_to_visit[next_position] = true
				o_path[next_position] = append(o_path[parent], next_move)
			}
		}
		visited[parent] = true
	}
	return o_path[destination]
}

func GetTailPosition(snake Snake) Point {
	return snake.Body.Data[snake.Length-1]
}

// give this function a point on the map, and it will return the adjacent valid positions
func GetAdjacentPoints(point Point, world *MoveRequest) []Point {
	output := []Point{
		Point{X: point.X + 1, Y: point.Y},
		Point{X: point.X, Y: point.Y + 1},
		Point{X: point.X - 1, Y: point.Y},
		Point{X: point.X, Y: point.Y - 1},
	}
	return output
}

// this one only returns valid points
func GetValidAdjacentPoints(point Point, world *MoveRequest) []Point {
	var output []Point
	for _, adj_point := range GetAdjacentPoints(point, world) {
		if IsValidPointToMoveTo(adj_point, world) {
			output = append(output, adj_point)
		}
	}
	return output
}
