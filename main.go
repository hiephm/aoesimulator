package main

import (
	"fmt"
)

const (
	// Build Type
	IDLE       = ""
	WORKER     = "worker"
	TOOL_AGE   = "tool_age"
	BRONZE_AGE = "bronze_age"

	// Work type
	FOOD  = "F"
	WOOD  = "W"
	GOLD  = "G"
	SCOUT = "S"

	WORKER_SPAWN_TIME = 20
)

var totalTime = 900
var maxWorkers = 20

var workerCost = 50
var beCost = 30
var bgCost, bsCost = 120, 120
var bbCost, blCost, bmCost = 150, 150, 150

var timeToFindFoodSource0 = 60
var movingTimeToSource = 10
var initialWoodCollectTime = 30
var foodCollectTime = 25
var woodCollectTime = 20
var woodCollectAmount = 10
var goldCollectTime = 25

var bgCount, bsCount, bbCount, blCount, bmCount, beCount = 0, 0, 0, 0, 0, 0

type Worker struct {
	MovingTime  int
	CollectTime int
	Work        string // F, W, G, S, IDLE
}

type TownCenter struct {
	BuildTime int
	BuildType string // worker, tool_age, bronze_age
}

type FoodSource struct {
	Name    string
	Workers []*Worker
	Amount  int
}

type WoodSource struct {
	Workers []*Worker
}

var currentTime = 0
var workers []*Worker
var town *TownCenter
var food, wood, gold = 200, 200, 0
var foodSources []*FoodSource
var woodSource *WoodSource

func init() {
	// Start game with 3 workers
	town = new(TownCenter)
	workers = []*Worker{
		{}, {}, {},
	}
	foodSources = []*FoodSource{
		{Amount: 900, Name: "Food 0"},
		{Amount: 600, Name: "Food 1"},
		{Amount: 750, Name: "Food 2"},
	}
	woodSource = new(WoodSource)
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
	if len(workers) == maxWorkers {
		return
	}

	if t.BuildType == WORKER {
		if t.BuildTime == WORKER_SPAWN_TIME {
			workers = append(workers, &Worker{})
			t.BuildTime = 1
			t.BuildType = IDLE
			output("Worker spawned.")
			if len(workers) == 7 {
				workers[6].Work = SCOUT
				output("Assigned worker to scout.")
			}
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
	if (beCount+1)*4 <= len(workers)+2 {
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

func (f *FoodSource) collectFood() {
	for _, worker := range f.Workers {
		if worker.MovingTime < movingTimeToSource {
			worker.MovingTime++
			continue
		}
		if worker.CollectTime < foodCollectTime {
			worker.CollectTime++
			continue
		}
		food += 10
		f.Amount -= 10
		worker.CollectTime = 1
	}
}

func (f *FoodSource) assignWorker(w *Worker) {
	w.Work = FOOD
	f.Workers = append(f.Workers, w)
	output("Assigned worker to food maker: " + f.Name)
}

func collectFood() {
	// food source 0
	if bgCount == 1 {
		for i, worker := range workers {
			if i < 6 && worker.Work == IDLE {
				foodSources[0].assignWorker(worker)
			}
		}
	}

	// food source 1
	if bsCount == 2 {
		for i, worker := range workers {
			if i >= 7 && worker.Work == IDLE {
				foodSources[1].assignWorker(worker)
			}
		}
	}

	for _, foodSource := range foodSources {
		foodSource.collectFood()
	}
}

func (ws *WoodSource) assignWorker(w *Worker) {
	w.Work = WOOD
	ws.Workers = append(ws.Workers, w)
	output("Assigned worker to wood cutter")
}

func (ws *WoodSource) adjustWorkers() {
	for len(ws.Workers) > 6 {
		var w *Worker
		w, ws.Workers = ws.Workers[0], ws.Workers[1:]
		w.Work = IDLE
		output("Change wood cutter to idle")
	}
}

func (ws *WoodSource) collectWood(collectTime int) {
	for _, worker := range ws.Workers {
		if worker.CollectTime < collectTime {
			worker.CollectTime++
			continue
		}
		wood += woodCollectAmount
		worker.CollectTime = 1
	}
}

func collectWood() {
	if bsCount == 0 {
		// Collecting initial woods without BS
		for i, worker := range workers {
			if i >= 7 && worker.Work == IDLE {
				woodSource.assignWorker(worker)
			}
		}
		woodSource.collectWood(initialWoodCollectTime)
	} else {
		// Collecting woods with BS
		woodSource.adjustWorkers()
		woodSource.collectWood(woodCollectTime)
	}
}

func advanceAge() {
	if food >= 500 {
		output("Advance to Bronze age")
		food -= 500
	}
}

func output(msg string) {
	foodMakers, woodCutters, goldMiners, scout := countWorkers()
	minute := currentTime / 60
	second := currentTime % 60
	fmt.Printf("[%02d:%02d] [Pop:%2d/%2d] [FM:%2d WC:%2d GM:%2d S:%d] [F:%3d W:%3d G:%3d] [F0:%3d F1:%3d F2:%3d]. Msg: %s\n",
		minute, second,
		len(workers), (beCount+1)*4,
		foodMakers, woodCutters, goldMiners, scout,
		food, wood, gold,
		foodSources[0].Amount, foodSources[1].Amount, foodSources[2].Amount,
		msg)
}

func countWorkers() (foodMakers, woodCutters, goldMiners, scout int) {
	for _, worker := range workers {
		switch worker.Work {
		case FOOD:
			foodMakers++
		case WOOD:
			woodCutters++
		case GOLD:
			goldMiners++
		case SCOUT:
			scout++
		}
	}
	return
}
