package main

import (
	"fmt"
	"log"
	"os"

	"lemin/funcs"
	"lemin/graph"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		log.Println("Try passing the name of a file with ants")
		return
	}

	err := funcs.CheckInput(args[0])
	if err != nil {
		log.Println(err.Error())
		return
	}

	//Create Graph
	test := &graph.Graph{}

	err = test.NewVertexes()
	if err != nil {
		log.Println(err)
		return
	}
	err = test.NewEdges()
	if err != nil {
		log.Println(err.Error())
		return
	}

	//Will save found paths in pahts and pahts2 arguments
	paths := make(map[string][]*graph.Vertex)
	paths2 := make(map[string][]*graph.Vertex)

	start := ""

	for k, v := range funcs.Nodes {
		if v == "start" {
			start = k
		}
	}
	end := ""

	for k, v := range funcs.Nodes {
		if v == "end" {
			end = k
		}
	}

	//Finds paths with BFS method
	paths, err = test.BFSS(paths, start, end)
	if err != nil {
		log.Println(err.Error())
		return
	}

	//Reset the Graph - edges
	for _, v := range test.Vertices {
		v.Adjacent = []*graph.Vertex{}
	}
	test.ResetEdges()

	//Finds paths with Bhandari method
	paths2 = test.BhanS(start, end)

	//Reset the graph
	test.ResetVertexes()
	test.ResetEdges()

	//Распределяем куда какой муравей пойдет
	var BFS []graph.Pt
	var Bhandari []graph.Pt
	BFS, Bhandari = graph.SortAnts(paths, paths2)

	source := test.GetVertex(start)
	target := test.GetVertex(end)

	//Проверяем BFS метод и Bhandari метод
	//Сравниваем между собой методы и выдаем результат
	//Print input
	for _, ex := range funcs.Example {
		fmt.Println(ex)
	}
	fmt.Println()
	test.MoveAnts(BFS, Bhandari, source, target)
}
