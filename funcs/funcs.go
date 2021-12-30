package funcs

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

//Example saves input
var Example = []string{}

//Ants - number of ants in the graph
var Ants int

//Nodes - names of vertexes
var Nodes = make(map[string]string)

//Connections - links between rooms
var Connections = make(map[string][]string)

//CheckInput checks file for valid input
func CheckInput(file string) error {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return err
	}

	var input []string

	scanner := bufio.NewScanner(f)

	//Сохраняет все строки в массиве без \n
	for scanner.Scan() {
		input = append(input, scanner.Text())
		Example = append(Example, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if len(input) == 0 {
		return errors.New("ERROR: invalid data format, empty file")
	}

	nants, err := strconv.Atoi(input[0])
	if err != nil {
		return errors.New("ERROR: invalid data format, invalid number of Ants")
	}

	if nants <= 0 {
		return errors.New("ERROR: invalid data format, invalid number of Ants")
	}

	if nants > 100000 {
		return errors.New("ERROR: invalid data format, too many Ants")
	}

	Ants = nants
	isstart := false
	isend := false
	cstart := 0
	cend := 0
	//Проверка на наличие стартовой и конечной комнаты и на количество таких комнат
	for _, v := range input {
		if v == "##start" {
			cstart++
			isstart = true
		}
		if v == "##end" {
			cend++
			isend = true
		}
	}

	if !isstart || !isend {
		return errors.New("ERROR: invalid data format, no start room or no end room found")
	}
	if cstart > 1 || cend > 1 {
		return errors.New("ERROR: invalid data format, more than 1 start or end rooms found")
	}

	//names of rooms
	var nrooms []string
	//names of rooms with coordinates
	var crooms []string

	var coords []string
	isDuplicate := false
	iofend := 0
	iofstart := 0
	// loop:
	lastroom := 0
	//Проверка на дублирующие комнаты
	//Проверка на наличие стартовой и конечных комнат, а также связей между комнатами
	for i, v := range input {
		if i == 0 {
			continue
		}

		if v == "##start" {
			iofstart = i
			if i+1 != len(input) {
				nstart := input[iofstart+1]
				if strings.HasPrefix(nstart, " ") || strings.HasPrefix(nstart, "L") || strings.HasPrefix(nstart, "#") {
					return errors.New("ERROR: invalid data format, invalid start room")
				}
			}
			continue
			// break
		}

		if v == "##end" {
			iofend = i
			if i+1 != len(input) {
				nend := input[iofend+1]
				if nend == "" || strings.HasPrefix(nend, " ") || strings.HasPrefix(nend, "L") || strings.HasPrefix(nend, "#") {
					return errors.New("ERROR: invalid data format, invalid end room")
				}
			}
			continue
			// break
		}

		if toIgnore(v) {
			continue
		}

		if !isRoom(v) {
			break
		}

		lastroom = i
		isUnique := true

		nroom := ""
		whitesp := 0
		ix := 0
		for indx, l := range v {
			if indx == 0 && (l == ' ' || l == 'L' || l == '#') {
				return errors.New("ERROR: invalid data format, invalid name of a room")
			}

			if (l != ' ') && (indx+1 == len(v)) {
				return errors.New("ERROR: invalid data format, invalid name of a room")
			}

			if l == ' ' {
				ix = indx
				break
			}
			nroom += string(l)
		}

		for ix < len(v) {
			if v[ix] == ' ' {
				whitesp++
			}
			ix++
		}
		if (len(v) - len(nroom)) == whitesp {
			return errors.New("ERROR: invalid data format, invalid name of a room")
		}

		croom := ""
		for _, l := range v {
			croom += string(l)
		}
		crooms = append(crooms, croom)

		xy := ""
		c := 0
		for _, l := range v {
			if l == ' ' && c == 0 {
				c++
				continue
			}

			if c == 0 {
				continue
			}

			xy += string(l)
		}
		coords = append(coords, xy)

		for _, r := range nrooms {
			if nroom == r {
				isUnique = false
				break
			}
		}

		if !isUnique {
			isDuplicate = true
			break
		}
		nrooms = append(nrooms, nroom)
	}

	if isDuplicate {
		return errors.New("ERROR: invalid data format, duplicated rooms with the identical names")
	}

	if len(input) == iofstart+1 {
		return errors.New("ERROR: invalid data format, no start room and no links")
	}

	if len(input) == iofend+1 {
		return errors.New("ERROR: invalid data format, no end room and no links")
	}

	for _, r := range nrooms {
		if input[iofend+1] == r {
			return errors.New("ERROR: invalid data format, a room is the same as the end room")
		}
	}
	for _, r := range nrooms {
		if input[iofstart+1] == r {
			return errors.New("ERROR: invalid data format, a room is the same as the start room")
		}
	}

	//Проверка на комнаты с одинаковыми координатами
	coord := []string{}
	isCoords := true
	for _, c := range coords {
		isUnique := true

		for _, r := range coord {
			if r == c {
				isUnique = false
				break
			}
		}

		if !isUnique {
			isCoords = false
			break
		}

		coord = append(coord, c)
	}

	if !isCoords {
		return errors.New("ERROR: invalid data format, different rooms with the same coordinates")
	}

	//Проверка на наличие связей между комнатами
	areLinks := false
	for j, l := range input {
		if j < lastroom {
			continue
		}

		ignor := false
		if toIgnore(l) {
			ignor = true
		}

		if !ignor && !isRoom(l) {
			areLinks = true
			break
		}
	}

	if !areLinks {
		return errors.New("ERROR: invalid data format, no links found")
	}
	isCoor := true
	//Проверка недопустимых координат комнат
loop:
	for _, r := range crooms {
		name := ""
		x := ""
		y := ""

		n := 0
		for _, l := range r {
			if l == ' ' {
				n++
				continue
			}

			if n == 0 {
				name += string(l)
			}
			if n == 1 {
				x += string(l)
			}
			if n == 2 {
				y += string(l)
			}
			if n > 2 {
				isCoor = false
				break loop
			}
		}
		if len(name) == 0 {
			return errors.New("ERROR: invalid data format, empty line")

		}
		if name[0] == 'L' || name[0] == '#' {
			return errors.New("ERROR: invalid data format, name of the room can not start with letter L or #")
		}

		xn, err := strconv.Atoi(x)
		if err != nil {
			isCoor = false
			break
		}

		if (x != "0" && xn == 0) || xn < 0 {
			isCoor = false
			break
		}

		yn, err := strconv.Atoi(y)
		if err != nil {
			isCoor = false
			break
		}

		if (y != "0" && yn == 0) || yn < 0 {
			isCoor = false
			break
		}

	}

	if !isCoor {
		return errors.New("ERROR: invalid data format, invalid coordinates")
	}
	//Добавляет имена нод для дальнейшего добавления в граф
	err = AddNodes(input, nrooms, iofstart, iofend)
	if err != nil {
		return err
	}

	//Проверка на наличие связей
	//Проверка на валидность связей
	var links []string
	isOkLinks := true
	for j, l := range input {
		if j <= lastroom {
			continue
		}

		if l == "##start" || l == "##end" {
			return errors.New("ERROR: invalid data format, invalid start room or end room")

		}
		if toIgnore(l) {
			continue
		}

		if isRoom(l) {
			isOkLinks = false
			break
		}

		ok := checkLrooms(l, nrooms)

		if !ok {
			isOkLinks = false

			break
		}

		links = append(links, l)
	}

	if len(links) < 1 {
		return errors.New("ERROR: invalid data format, no links found or unacceptable order of input")
	}

	if !isOkLinks {
		return errors.New("ERROR: invalid data format, invalid links")

	}

	return nil
}

//checkLrooms проверяет на валидность ребер между вершинами
func checkLrooms(s string, nrooms []string) bool {
	r1 := ""
	r2 := ""
	c := 0
	ok := true
	for _, l := range s {
		if l == '-' {
			c++
			continue
		}

		if c == 0 {
			r1 += string(l)
			continue
		}

		r2 += string(l)

		//more "-" than allowed
		if c > 1 {
			ok = false

			break
		}
	}
	//Link to themselves
	if r1 == r2 {

		ok = false
	}

	if !ok {

		return false
	}

	//links to unknown rooms
	isR1 := false
	isR2 := false
	for _, r := range nrooms {
		if r == r1 {
			isR1 = true
		}

		if r == r2 {
			isR2 = true
		}
	}

	if !isR1 || !isR2 {
		return false
	}

	//Сохраняет связь между ребрами в мапу для дальнейшего добавления в граф
	AddConnection(r1, r2)

	return true
}

//находит комментарии в инпуте
func toIgnore(s string) bool {
	if strings.HasPrefix(s, "#") {
		return true
	}

	return false
}

//
func isRoom(s string) bool {
	if strings.Contains(s, "-") {
		return false
	}

	return true
}

//AddNodes находит и сохраняет комнаты в глобальной map
//А также находит стартовую и конечную комнаты и помечает их
func AddNodes(input, nrooms []string, iofstart, iofend int) error {
	startr := ""
	for _, r := range input[iofstart+1] {
		if r == ' ' {
			break
		}
		startr += string(r)
	}

	endr := ""
	for _, r := range input[iofend+1] {
		if r == ' ' {
			break
		}
		endr += string(r)
	}

	for _, v := range nrooms {
		if v == startr {
			Nodes[startr] = "start"
			continue
		}

		if v == endr {
			Nodes[endr] = "end"
			continue
		}

		Nodes[v] = "r"
	}
	if startr == "" || endr == "" {
		return errors.New("Invalid start room or end room")
	}
	return nil
}

//AddConnection adds connection between rooms to the map
func AddConnection(r1, r2 string) {
	Connections[r1] = append(Connections[r1], r2)
}
