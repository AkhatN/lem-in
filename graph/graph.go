package graph

import (
	"errors"
	"fmt"
	"lemin/funcs"
)

//Graph struct
type Graph struct {
	Vertices []*Vertex
}

//Vertex struct
type Vertex struct {
	Name     string
	Adjacent []*Vertex
	Dist     int     // Distance from start
	Prev     *Vertex // Previous vertex
	Nofants  int     // number of ants in the room
	Busy     bool    // shows if the room is occupied by an ant
	Nofant   []int   // shows the Name of the ant that is in the room
	Room     string  // shows what type of room it is
}

//Pt shows length of the path and how many ants in the path already, and indexes of them, and path with vertexes
type Pt struct {
	Length int
	nAnts  int
	iAnts  []int
	Pth    []*Vertex
}

//Queue struct
type Queue struct {
	items []*Vertex
}

//Ants1 ..
var Ants1 = []int{}

//AddVertex adds vertexes to the graph
func (g *Graph) AddVertex(k, v string) error {
	if contains(g.Vertices, k) {
		err := fmt.Errorf("ERROR: invalid data format, vertex %v not added because it is an existing key", k)
		fmt.Println(err.Error())
		return err
	}

	g.Vertices = append(g.Vertices, &Vertex{Name: k, Room: v})
	return nil
}

//UpVertexes updates vertexes to the graph
func (g *Graph) UpVertexes(k, v string) {
	if contains(g.Vertices, k) {
		return
	}

	g.Vertices = append(g.Vertices, &Vertex{Name: k, Room: v})
}

//AddEdge adds edges to vertexes
func (g *Graph) AddEdge(from, to string) error {
	fromVertex := g.GetVertex(from)
	toVertex := g.GetVertex(to)

	if fromVertex == nil || toVertex == nil {
		err := fmt.Errorf("ERROR: invalid data format, invalid edge to non-existing vertexes (%v--%v)", from, to)
		// fmt.Println(err.Error())
		return err
	}

	if contains(fromVertex.Adjacent, to) {
		err := fmt.Errorf("ERROR: invalid data format, multi edge is invalid (%v--%v)", from, to)
		// fmt.Println(err.Error())
		return err
	}

	if fromVertex == toVertex {
		err := fmt.Errorf("ERROR: invalid data format, self loop is invalid (%v--%v)", from, to)
		// fmt.Println(err.Error())
		return err
	}

	fromVertex.Adjacent = append(fromVertex.Adjacent, toVertex)
	toVertex.Adjacent = append(toVertex.Adjacent, fromVertex)
	return nil
}

//UpEdge updates edges to vertexes
func (g *Graph) UpEdge(from, to string) {
	fromVertex := g.GetVertex(from)
	toVertex := g.GetVertex(to)

	if fromVertex == nil || toVertex == nil {
		return
	}

	if contains(fromVertex.Adjacent, to) {
		return
	}

	if fromVertex == toVertex {
		return
	}

	fromVertex.Adjacent = append(fromVertex.Adjacent, toVertex)
	toVertex.Adjacent = append(toVertex.Adjacent, fromVertex)
	return
}

//GetVertex finds the vertex and returns it
func (g *Graph) GetVertex(k string) *Vertex {
	for i, v := range g.Vertices {
		if v.Name == k {
			return g.Vertices[i]
		}
	}
	return nil
}

//checks if vertex already exists
func contains(s []*Vertex, k string) bool {
	for _, v := range s {
		if k == v.Name {
			return true
		}
	}
	return false
}

func (g *Graph) createVisited() map[*Vertex]bool {
	visited := make(map[*Vertex]bool, len(g.Vertices))
	for _, v := range g.Vertices {
		visited[v] = false
	}
	return visited
}

//Enqueue adds new item to the queue
func (q *Queue) Enqueue(i *Vertex) {
	q.items = append(q.items, i)
}

//Dequeue removes an item from the queue and returns the item
func (q *Queue) Dequeue() *Vertex {
	toRemove := q.items[0]
	q.items = q.items[1:]
	return toRemove
}

//NewVertexes adds vertexes to the graph
func (g *Graph) NewVertexes() error {
	for k, v := range funcs.Nodes {
		err := g.AddVertex(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

//NewEdges adds edges between vertexes
func (g *Graph) NewEdges() error {
	for k, v := range funcs.Connections {
		for _, r := range v {
			err := g.AddEdge(k, r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//ResetVertexes adds vertexes to the graph
func (g *Graph) ResetVertexes() {
	for k, v := range funcs.Nodes {
		g.UpVertexes(k, v)
	}
}

//ResetEdges adds edges between vertexes
func (g *Graph) ResetEdges() {
	for k, v := range funcs.Connections {
		for _, r := range v {
			g.UpEdge(k, r)
		}
	}
}

//BFSS finds paths with BFS method
func (g *Graph) BFSS(paths map[string][]*Vertex, start, end string) (map[string][]*Vertex, error) {
	i := 1
	pp1, ispath11 := g.Findpath(paths, start, end)

	paths["1"] = pp1
	if ispath11 != false {
		for pp, ispath := g.Findpath(paths, start, end); ispath != false || pp != nil; {
			path1, ispath1 := g.Findpath(paths, start, end)

			if ispath1 == false {
				break
			}
			i++

			name := fmt.Sprintf("%v", i)
			paths[name] = path1
		}
	}

	if len(pp1) == 0 {

		return paths, errors.New("Error: invalid data format, no path between start and end rooms")
	}

	return paths, nil
}

//BhanS finds paths with Bhandari method
func (g *Graph) BhanS(start, end string) map[string][]*Vertex {
	paths2 := make(map[string][]*Vertex)
	paths3 := make(map[string][]*Vertex)
	p1, ispath2 := g.Findpath(paths2, start, end)
	paths2["1"] = p1
	g.RevEdges(p1, start, end)
	i := 2

	if ispath2 != false {
		for g.Bellman(start, end) != nil {
			path1 := g.Bellman(start, end)
			g.RevEdges(path1, start, end)

			name := fmt.Sprintf("%v", i)
			paths2[name] = path1
			i++

		}

		g.NewGraph(paths2)

		i = 1
		for pp2, ispath := g.Findpath(paths3, start, end); ispath != false || pp2 != nil; {
			path1, ispath3 := g.Findpath(paths3, start, end)
			if ispath3 == false {
				break
			}
			name := fmt.Sprintf("%v", i)
			paths3[name] = path1
			i++

		}

	} else {
		paths3["1"] = p1
	}
	return paths3
}

//SortAnts sorts ants between paths
func SortAnts(paths, paths2 map[string][]*Vertex) ([]Pt, []Pt) {
	var BFS []Pt
	var Bhandari []Pt
	x := 0
	short1 := ""
	lengthp := 0
	mapabfs := make(map[string]int)
	for mm := range paths {
		mapabfs[mm] = 0
	}

	for len(BFS) != len(paths) {
		lengthp = 9999
		short1 = ""
		for n1, v := range paths {
			if mapabfs[n1] == 0 && x == 0 {
				short1 = n1
				x++
				lengthp = len(v)
				continue
			}
			if mapabfs[n1] == 1 {
				continue
			}
			if short1 == "" {
				continue
			}

			if len(v) < lengthp {
				lengthp = len(v)
				short1 = n1
			}
		}
		l11 := lengthp - 2
		ps := []*Vertex{}
		for _, v1 := range paths[short1] {
			ps = append(ps, v1)
		}

		pass1 := Pt{l11, 0, []int{}, ps}
		BFS = append(BFS, pass1)
		mapabfs[short1] = 1
		x = 0
	}

	x = 0
	short1 = ""
	lengthp = 0
	mapabfs1 := make(map[string]int)
	for mm := range paths2 {
		mapabfs1[mm] = 0
	}
	for len(Bhandari) != len(paths2) {
		lengthp = 9999
		short1 = ""
		for n1, v := range paths2 {
			if mapabfs1[n1] == 1 {
				continue
			}
			if mapabfs1[n1] == 0 && x == 0 {
				short1 = n1
				x++
				lengthp = len(v)
				continue
			}

			if short1 == "" {
				continue
			}

			if len(v) < lengthp {
				lengthp = len(v)
				short1 = n1
			}
		}

		l22 := lengthp - 2
		ps := []*Vertex{}
		for _, v1 := range paths2[short1] {
			ps = append(ps, v1)
		}

		pass1 := Pt{l22, 0, []int{}, ps}
		Bhandari = append(Bhandari, pass1)
		mapabfs1[short1] = 1
		x = 0
	}

	for j := 1; j <= funcs.Ants; j++ {
		Ants1 = append(Ants1, j)
	}

	//BFS
	//Нужно отсортировать пути для BFS  Bhandari структур по возрастанию длины путей
	//Отправим первого муравья в самый короткий первый путь
	inpath := 0
	short := 0
	for pp11, nn := range BFS {
		if pp11 == 0 {
			short = len(nn.Pth)
			continue
		}

		if len(nn.Pth) < short {
			short = len(nn.Pth)
			inpath = pp11
		}
	}
	BFS[inpath].nAnts++
	BFS[inpath].iAnts = append(BFS[inpath].iAnts, 1)
	//Считаем сколько куда отправить муравьев

	b := 1
	// curPath := 0
	if len(Ants1) > 1 {
		if len(BFS) > 1 {
			for b < len(Ants1) {
				currentAnt := Ants1[b]
				currentPath := 0
				nRAnts := 0
				for ind, p := range BFS {
					if ind == 0 {
						currentPath = 0
						nRAnts = p.nAnts + p.Length
						continue
					}

					if p.nAnts+p.Length < nRAnts {
						currentPath = ind
						nRAnts = p.nAnts + p.Length
					}
				}
				BFS[currentPath].iAnts = append(BFS[currentPath].iAnts, currentAnt)
				BFS[currentPath].nAnts++
				b++
			}
		} else {
			for b < len(Ants1) {
				BFS[0].iAnts = append(BFS[0].iAnts, Ants1[b])
				BFS[0].nAnts++
				b++
			}
		}
	}

	//Bhandari
	//Отправим первого муравья в самый короткий первый путь
	inpath = 0
	short = 0
	for pp11, nn := range Bhandari {
		if pp11 == 0 {
			short = len(nn.Pth)
			continue
		}

		if len(nn.Pth) < short {
			short = len(nn.Pth)
			inpath = pp11
		}
	}

	Bhandari[inpath].nAnts++
	Bhandari[inpath].iAnts = append(Bhandari[inpath].iAnts, 1)
	//Считаем сколько куда отправить муравьев

	b = 1
	if len(Ants1) > 1 {
		if len(Bhandari) > 1 {
			for b < len(Ants1) {
				currentAnt := Ants1[b]
				currentPath := 0
				nRAnts := 0
				for ind, p := range Bhandari {
					if ind == 0 {
						currentPath = 0
						nRAnts = p.nAnts + p.Length
						continue
					}

					if p.nAnts+p.Length < nRAnts {
						currentPath = ind
						nRAnts = p.nAnts + p.Length
					}
				}
				Bhandari[currentPath].iAnts = append(Bhandari[currentPath].iAnts, currentAnt)
				Bhandari[currentPath].nAnts++
				b++
			}
		} else {
			for b < len(Ants1) {
				Bhandari[0].iAnts = append(Bhandari[0].iAnts, Ants1[b])
				Bhandari[0].nAnts++
				b++
			}
		}
	}

	return BFS, Bhandari
}

//Remactive removes an ant from the active queue
func Remactive(aq *ActiveQ, ant int) {
	ants1 := []int{}
	for _, a := range aq.ants {
		if a != ant {
			ants1 = append(ants1, a)
		}
	}
	aq.ants = ants1
}

//Addactiveants adds ants
func (qu *Que) Addactiveants(activeants ActiveQ) {
	for _, ant := range activeants.ants {
		qu.ants = append(qu.ants, ant)
	}
}

//ActiveQ очередь из активных муравьев
type ActiveQ struct {
	ants []int
}

//Que очередь из муравьев которые буду идти в комнаты
type Que struct {
	ants []int
}

//Enq adds ants to the queue
func (qu *Que) Enq(ant int) {
	qu.ants = append(qu.ants, ant)
}

//Deq removes first ant from the queue and returns the ant
func (qu *Que) Deq() int {
	toRemove := qu.ants[0]
	qu.ants = qu.ants[1:]
	return toRemove
}

//Findpath - search for one shortest path with BFS method
func (g *Graph) Findpath(paths map[string][]*Vertex, source, target string) ([]*Vertex, bool) {
	path := []*Vertex{}
	visited := g.createVisited()
	for _, v := range g.Vertices {
		v.Dist = 99999
		v.Prev = nil
	}

	start := g.GetVertex(source)
	end := g.GetVertex(target)

	start.Dist = 0

	//проходимся по известным путям, и помечаем эти комнаты как посещенные, чтобы при поиске новых путей в них не заходить
	for _, ps := range paths {
		for _, v := range ps {
			visited[v] = true
		}
	}

	visited[end] = false
	visited[start] = true

	q := Queue{items: []*Vertex{}}
	q.Enqueue(start)

	for len(q.items) > 0 {
		currentV := q.items[0]
		q.items = q.items[1:]

		if currentV == end {
			path = Savepath(start, currentV)
			ispath := true
			if len(path) == 2 {
				ispath = false
			}
			return path, ispath
		}

		// Проходимся по всем примыкающим соседям/нодам текущей ноды, которые не были посещены
		for _, v := range currentV.Adjacent {
			for w, b := range visited {
				if w == v && b == false {
					// d - new Distance to the vertex from the start
					// v.Dist - known Distance to the vertex from the start
					d := currentV.Dist + 1
					// Если новая дистанция до ноды меньше известной, то обновляем дистанцию до этой ноды и меняем предыдущую ноду
					if d < v.Dist {
						v.Dist = d
						v.Prev = currentV
					}
					break
				}
			}
		}

		// добавляем соседей текущей ноды в очередь для посещения в след цикле. Без разницы в каком порядке, потому что у всех дистанция одинаковая
		for _, node := range currentV.Adjacent {
			if visited[node] == false {
				visited[node] = true
				q.items = append(q.items, node)
			}
		}
	}

	return nil, false
}

//Savepath saves found path
func Savepath(start, end *Vertex) []*Vertex {
	p := []*Vertex{}
	u := end

	if u.Prev != nil || u == start {
		for u != nil {
			p = append(p, u)
			u = u.Prev
		}
	}

	path := []*Vertex{}

	for i := len(p) - 1; i >= 0; i-- {
		path = append(path, p[i])
	}

	return path
}

//RevEdges reverses edges
func (g *Graph) RevEdges(path []*Vertex, source, target string) {
	start := g.GetVertex(source)
	end := g.GetVertex(target)

	p := []*Vertex{}
	u := end

	if u.Prev != nil || u == start {
		for u != nil {
			p = append(p, u)
			u = u.Prev
		}
	}

	for _, v := range p {
		next := v
		if v.Prev != nil {
			current := v.Prev
			for _, adj := range current.Adjacent {
				if adj == next {
					next.Prev = nil
					g.deleteEdge(current, next)
					break
				}
			}
		}
	}
}

//deleteEdge deletes edge
func (g *Graph) deleteEdge(current, next *Vertex) {
	newadj := []*Vertex{}
	for _, v := range current.Adjacent {
		if v != next {
			newadj = append(newadj, v)
		}
	}
	current.Adjacent = newadj

}

//Bellman finds all intersecting paths for Bhandari method
func (g *Graph) Bellman(source, target string) []*Vertex {
	path := []*Vertex{}
	visited := g.createVisited()

	for _, v := range g.Vertices {
		v.Dist = 99999
		v.Prev = nil
	}

	start := g.GetVertex(source)
	end := g.GetVertex(target)

	start.Dist = 0

	visited[start] = true

	q := Queue{items: []*Vertex{}}
	q.Enqueue(start)

	for len(q.items) > 0 {
		currentV := q.items[0]
		q.items = q.items[1:]

		if currentV == end {
			path = Savepath(start, currentV)
			return path
		}

		// Проходимся по всем примыкающим соседям/нодам текущей ноды, которые не были посещены
		for _, v := range currentV.Adjacent {
			for w, b := range visited {
				if w == v && b == false {
					// d - new Distance to the vertex from the start
					// v.Dist - known Distance to the vertex from the start
					d := currentV.Dist + 1
					// Если новая дистанция до ноды меньше известной, то обновляем дистанцию до этой ноды и меняем предыдущую ноду
					if d < v.Dist {
						v.Dist = d
						v.Prev = currentV
					}
					break
				}
			}
		}

		// добавляем соседей текущей ноды в очередь для посещения в след цикле. Без разницы в каком порядке, потому что у всех дистанция одинаковая
		for _, node := range currentV.Adjacent {
			if visited[node] == false {
				visited[node] = true
				q.items = append(q.items, node)
			}
		}
	}

	return nil
}

//NewGraph creates new graph according to Bellman graph for Bhandari method
func (g *Graph) NewGraph(paths map[string][]*Vertex) {
	for _, v := range g.Vertices {
		v.Dist = 99999
		v.Prev = nil
		v.Adjacent = []*Vertex{}
	}

	for _, ps := range paths {
		for i, v := range ps {
			if i == len(ps)-1 {
				break
			}

			next := ps[i+1]
			v.Adjacent = append(v.Adjacent, next)
		}
	}

	//Идем циклом по всем нодам в графе
	for _, vert := range g.Vertices {
	loop:
		//идем циклом по всем соседям текущей ноды
		for _, v := range vert.Adjacent {
			// проходимся циклом по соседям следующей ноды
			// Если там у следующей ноды в соседях есть текущая нода, то значит они инверсны и нужно удалить связь между текущей и следующей
			for _, w := range v.Adjacent {
				if vert == w {
					g.deleteInverse(vert, v)
					break loop
				}
			}
		}
	}
}

func (g *Graph) deleteInverse(current, next *Vertex) {
	adj1 := []*Vertex{}
	adj2 := []*Vertex{}

	for _, v1 := range current.Adjacent {
		if v1 != next {
			adj1 = append(adj1, v1)
		}
	}
	current.Adjacent = adj1

	for _, v2 := range next.Adjacent {
		if v2 != current {
			adj2 = append(adj2, v2)
		}
	}

	next.Adjacent = adj2
}

//MoveAnts moves ants between rooms from start room to the end room
func (g *Graph) MoveAnts(BFS, Bhandari []Pt, source, target *Vertex) {
	//здесь храним строки из результатов для вывода и сравнения количества строк
	BFSPath := []string{}
	BhPath := []string{}

	for _, an := range Ants1 {
		source.Nofant = append(source.Nofant, an)
	}

	//Нужна очередь для активных муравьев, которые не дошли до конца
	//Эта очередь обновляется, когда идет проверка для следующих муравьев, которые только начнут идти
	//И очередь обнуляется в конце, затем обновляется добавляя муравьев, которые все еще активны

	activeants := ActiveQ{ants: []int{}}
	qu := Que{ants: []int{}}
	for _, p := range BFS {
		if len(p.iAnts) != 0 {
			activeants.ants = append(activeants.ants, p.iAnts[0])
		}
	}
	qu.Addactiveants(activeants)

	indexofq := 1

	for len(target.Nofant) != len(Ants1) {
		if len(qu.ants) == 0 {
			qu.Addactiveants(activeants)
		}
		str := ""
		for len(qu.ants) != 0 {
			curant := qu.ants[0]
			// Убрать из очереди
			qu.ants = qu.ants[1:]
			curplace := &Vertex{}

			//Находит где находится сейчас муравей
		l1:
			for _, v := range g.Vertices {
				for _, curant1 := range v.Nofant {
					if curant1 == curant {
						curplace = v
						break l1
					}
				}
			}

			//находит индекс пути, то есть по которому пути идти муравью
			index := 0 //index of path
		l2:
			for j, p := range BFS {
				for _, inx := range p.iAnts {
					if curant == inx {
						index = j
						break l2
					}
				}
			}

			nextroom := &Vertex{}

			//Находит следующую комнату
			for k, r := range BFS[index].Pth {
				if r == curplace {
					nextroom = BFS[index].Pth[k+1]
					break
				}
			}

			//Теперь нужно разместить текущего муравья в следующей комнате
			nextroom.Nofant = append(nextroom.Nofant, curant)

			//Убрать из текущей комнаты
			if len(curplace.Nofant) > 0 {
				Nofant := []int{}
				for _, a := range curplace.Nofant {
					if a != curant {
						Nofant = append(Nofant, a)
					}
				}
				curplace.Nofant = Nofant
			} else {
				curplace.Nofant = []int{}
			}

			// Пометить текущую комнату свободной, пометить след комнату занятой
			curplace.Busy = false
			nextroom.Busy = true

			// Если след комната финиш, то убрать муравья из активной очереди
			if nextroom == target {
				Remactive(&activeants, curant)
			}
			str += fmt.Sprintf("L%v-%v", curant, nextroom.Name)
			if len(qu.ants) != 0 {
				str += " "
			}
			if len(BFS[0].Pth) == 2 && curant != len(Ants1) {
				str += " "
			}
		}
		if len(BFS[0].Pth) != 2 {
			str += "\n"
		}

		BFSPath = append(BFSPath, str)
		//Если есть муравьие которые все еще не активны, то есть находятся в старте, то добавить ждущих муравьев в активную очередь
		if len(source.Nofant) > 0 {
			for _, p := range BFS {
				for inx, a := range p.iAnts {
					if inx == indexofq {
						activeants.ants = append(activeants.ants, a)
						break
					}
				}
			}
		}
		indexofq++
	}

	for _, v := range g.Vertices {
		v.Nofant = []int{}
	}

	for _, an := range Ants1 {
		source.Nofant = append(source.Nofant, an)
	}
	//Нужна очередь для активных муравьев, которые не дошли до конца
	//Эта очередь обновляется, когда идет проверка для следующих муравьев, которые только начнут идти
	//И очередь обнуляется в конце, затем обновляется добавляя муравьев, которые все еще активны

	activeants = ActiveQ{ants: []int{}}
	qu = Que{ants: []int{}}
	for _, p := range Bhandari {
		if len(p.iAnts) != 0 {
			activeants.ants = append(activeants.ants, p.iAnts[0])
		}
	}
	qu.Addactiveants(activeants)

	indexofq = 1

	//Как обновлять очередь
	for len(target.Nofant) != len(Ants1) {
		if len(qu.ants) == 0 {
			qu.Addactiveants(activeants)
		}
		str := ""
		for len(qu.ants) != 0 {
			curant := qu.ants[0]
			// Убрать из очереди
			qu.ants = qu.ants[1:]
			curplace := &Vertex{}

			//Находит где находится сейчас муравей
		l3:
			for _, v := range g.Vertices {
				for _, curant1 := range v.Nofant {
					if curant1 == curant {
						curplace = v
						break l3
					}

				}
			}

			//находит индекс пути, то есть по которому пути идти муравью
			index := 0
		l4:
			for j, p := range Bhandari {
				for _, inx := range p.iAnts {
					if curant == inx {
						index = j
						break l4
					}
				}
			}

			nextroom := &Vertex{}

			//Находит следующую комнату
			for k, r := range Bhandari[index].Pth {
				if r == curplace {
					nextroom = Bhandari[index].Pth[k+1]
					break
				}
			}

			//Теперь нужно разместить текущего муравья в следующей комнате
			nextroom.Nofant = append(nextroom.Nofant, curant)

			//Убрать из текущей комнаты
			if len(curplace.Nofant) > 0 {
				Nofant := []int{}
				for _, a := range curplace.Nofant {
					if a != curant {
						Nofant = append(Nofant, a)
					}
				}
				curplace.Nofant = Nofant
			} else {
				curplace.Nofant = []int{}
			}

			// Пометить текущую комнату свободной, пометить след комнату занятой
			curplace.Busy = false
			nextroom.Busy = true

			// Если след комната финиш, то убрать муравья из активной очереди
			if nextroom == target {
				Remactive(&activeants, curant)
			}
			str += fmt.Sprintf("L%v-%v", curant, nextroom.Name)
			if len(qu.ants) != 0 {
				str += " "
			}
			if len(Bhandari[0].Pth) == 2 {
				str += " "
			}
		}
		if len(Bhandari[0].Pth) != 2 {
			str += "\n"
		}

		BhPath = append(BhPath, str)
		//Если есть муравьие которые все еще не активны, то есть находятся в старте, то добавить ждущих муравьев в активную очередь
		if len(source.Nofant) > 0 {
			for _, p := range Bhandari {
				for inx, a := range p.iAnts {
					if inx == indexofq {
						activeants.ants = append(activeants.ants, a)
						break
					}
				}
			}
		}
		indexofq++
	}

	if len(BhPath) < len(BFSPath) {
		for _, s := range BhPath {
			fmt.Print(s)
		}
		if len(Bhandari[0].Pth) == 2 {
			fmt.Println()
		}
		return
	}
	for _, s := range BFSPath {
		fmt.Print(s)
	}
	if len(BFS[0].Pth) == 2 {
		fmt.Println()
	}
}
