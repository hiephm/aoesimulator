package main

import "fmt"

const (
	// Build Type
	IDLE       = ""
	WORKER     = "worker"
	TOOL_AGE   = "tool_age"
	BRONZE_AGE = "bronze_age"

	// Resource Type
	FOOD = "F"
	WOOD = "W"
	GOLD = "G"

	WORKER_SPAWN_TIME = 20
)

var totalTime = 900
var maxWorkers = 20

var workerCost = 50
var beCost = 30
var bgCost, bsCost = 120, 120
var bbCost, blCost, bmCost = 150, 150, 150

var timeToFindFoodSource0 = 30
var movingTimeToSource = 10
var initialWoodCollectTime = 30
var foodCollectTime = 25
var woodCollectTime = 20
var goldCollectTime = 25

var bgCount, bsCount, bbCount, blCount, bmCount, beCount = 0, 0, 0, 0, 0, 0

type Worker struct {
	MovingTime   int
	CollectTime  int
	ResourceType string // F, W, G
}

type TownCenter struct {
	BuildTime int
	BuildType string // worker, tool_age, bronze_age
}

type FoodSource struct {
	Workers []*Worker
	Capacity int
	Remain int
}

type WoodSource struct {
	Workers []*Worker
}

var currentTime = 0
var idleWorkers []*Worker
var town *TownCenter
var food, wood, gold = 200, 200, 0

var foodSource0, foodSource1, foodSource2 = 900, 600, 750

func init() {
	// Start game with 3 workers
	town = new(TownCenter)
	idleWorkers = []*Worker{
		{}, {}, {},
	}
}

func main() {
	output("Start game")
	for currentTime = 0; currentTime <= totalTime; currentTime++ {
		town.spawnWorkers()
		buildStructures()
		collectFood()
		collectWood()
		advanceAge()
	}
}

func (t *TownCenter) spawnWorkers() {
	if len(idleWorkers) == maxWorkers {
		return
	}

	if t.BuildType == WORKER {
		if t.BuildTime == WORKER_SPAWN_TIME {
			idleWorkers = append(idleWorkers, &Worker{})
			t.BuildTime = 1
			t.BuildType = IDLE
			output("Worker spawned.")
		} else {
			t.BuildTime++
		}
	} else {
		if food < workerCost {
			return
		}
		food -= workerCost
		t.BuildType = WORKER
		t.BuildTime++
	}
}

func buildStructures() {
	if (beCount + 1) * 4 <= len(idleWorkers) + 2 {
		output("BE built.")
		wood -= beCost
		beCount++
	}

	if currentTime < timeToFindFoodSource0 {
		return
	}
	if bgCount == 0 {
		output("BG built.")
		wood -= bgCost
		bgCount++
	}

	if bsCount == 0 && wood > 120 {
		output("First BS built.")
		wood -= bsCost
		bsCount++
	}

	if bsCount == 1 && wood > 120 {
		output("Second BS built.")
		wood -= bsCost
		bsCount++
	}
}

func collectFood() {
	// food source 0
	if bgCount == 1 {
		for i, worker := range idleWorkers {
			if i >= 6 {
				break
			}
			if worker.ResourceType == IDLE {
				output("Assign worker to food maker (source 0)")
			}
			worker.ResourceType = FOOD
			if worker.MovingTime < movingTimeToSource {
				worker.MovingTime++
				continue
			}
			if worker.CollectTime < foodCollectTime {
				worker.CollectTime++
				continue
			}
			food += 10
			foodSource0 -= 10
			worker.CollectTime = 1
		}
	}

	if bsCount == 2 {
		for _, worker := range idleWorkers {
			if worker.ResourceType == IDLE {
				output("Assign worker to food maker (source 1)")
				worker.ResourceType = FOOD
				continue
			}

			if worker.ResourceType == FOOD {
				if worker.MovingTime < movingTimeToSource {
					worker.MovingTime++
					continue
				}
				if worker.CollectTime < foodCollectTime {
					worker.CollectTime++
					continue
				}
				food += 10
				foodSource1 -= 10
				worker.CollectTime = 1
			}
		}
	}
}

func collectWood() {
	// Collecting initial woods without BS
	if bsCount == 0 {
		for i, worker := range idleWorkers {
			if i < 7 {
				continue
			}
			if worker.ResourceType == IDLE {
				output("Assign worker to wood cutter")
			}
			worker.ResourceType = WOOD
			if worker.CollectTime < initialWoodCollectTime {
				worker.CollectTime++
				continue
			}
			wood += 10
			worker.CollectTime = 1
		}
	} else {
		// Collecting woods with BS
		if bsCount == 0 {
			return
		}
		for i, worker := range idleWorkers {
			if i < 7 || i > 12 {
				if i != 7 && worker.ResourceType == WOOD {
					worker.ResourceType = IDLE
					output("Change wood cutter to idle")
				}
				continue
			}
			if worker.ResourceType == IDLE {
				output("Assign worker to wood cutter")
			}
			worker.ResourceType = WOOD
			if worker.CollectTime < woodCollectTime {
				worker.CollectTime++
				continue
			}
			wood += 10
			worker.CollectTime = 1
		}
	}
}

func advanceAge() {
	if food >= 500 {
		output("Advance to Bronze age")
		food -= 500
	}
}

func output(msg string) {
	foodMakers, woodCutters, goldMiners := countWorkers()
	minute := currentTime / 60
	second := currentTime % 60
	fmt.Printf("[%02d:%02d] [BE:%2d] [W:%2d FM:%2d WC:%2d GM:%2d] [F:%3d W:%3d G:%3d] [F0:%3d F1:%3d F2:%3d]. Msg: %s\n",
		minute, second,
		(beCount+1)*4,
		len(idleWorkers), foodMakers, woodCutters, goldMiners,
		food, wood, gold,
		foodSource0, foodSource1, foodSource2,
		msg)
}

func countWorkers() (foodMakers, woodCutters, goldMiners int) {
	for _, worker := range idleWorkers {
		switch worker.ResourceType {
		case FOOD:
			foodMakers++
		case WOOD:
			woodCutters++
		case GOLD:
			goldMiners++
		}
	}
	return
}