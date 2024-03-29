package main

type StreetTime struct {
	name               string
	greenLigthDuration int
}

type output struct {
	intersectionId int
	intersection   Intersection
	streetsTime    []StreetTime
}

func algorithm(
	config Config,
	streets []*Street,
	carsPaths []*CarsPaths,
	streetsMap map[string]*Street,
	intersectionMap map[int]*Intersection,
	intersectionsList []*Intersection,
) []output {
	out := make([]output, 0)

	for _, carPath := range carsPaths {
		var intersectionId int
		var totTime int
		for _, street := range carPath.streetNames {
			street := streetsMap[street]
			totTime += street.timeNeeded

			intersectionId = street.endIntersection
		}

		if totTime <= config.simuDuration {
			intersection := intersectionMap[intersectionId]
			intersection.arrivingCars++
			intersectionMap[intersectionId] = intersection
		}
	}

	for _, intersection := range intersectionsList {
		visited := make(map[int]bool)
		dfs(
			visited,
			config.simuDuration,
			intersection,
			intersection.arrivingCars,
			intersection,
			intersectionMap,
		)
	}

	for _, intersection := range intersectionsList {
		streetTimes := make([]StreetTime, 0)
		totScore := 0

		for _, street := range intersection.incomingStreets {
			totScore += int(street.score)
		}
		if totScore == 0 {
			continue
		}

		a := len(intersection.incomingStreets) / len(intersection.outcomingStreets)
		if a == 0 {
			a = 1
		}

		for _, street := range intersection.incomingStreets {
			if street.score == 0 {
				continue
			}

			streetTimes = append(streetTimes, StreetTime{
				name:               street.name,
				greenLigthDuration: a, // / len(intersection.incomingStreets),
			})
		}

		out = append(out, output{
			intersectionId: intersection.id,
			streetsTime:    streetTimes,
		})
	}

	return out
}

func dfs(
	visited map[int]bool,
	remainingTime int,
	// incomingStreet *Street,
	intersection *Intersection,
	score int,
	startIntersection *Intersection,

	intersectionMap map[int]*Intersection,
) int {
	if visited := visited[intersection.id]; visited {
		return score
	}

	visited[intersection.id] = true

	// score += intersection.arrivingCars // TODO

	for streetName, incomingStreet := range intersection.incomingStreets {
		streetScore := dfs(
			visited,
			remainingTime-incomingStreet.timeNeeded,
			// incomingStreet,
			intersectionMap[incomingStreet.startIntersection],
			score+int(intersectionMap[incomingStreet.startIntersection].arrivingCars),

			startIntersection,
			intersectionMap,
		)
		// score += streetScore
		incomingStreet.score = streetScore
		intersection.incomingStreets[streetName] = incomingStreet
	}

	// visited[intersection.id] = false
	return score
}

func pickBestIntersection(intersectionsList []*Intersection) *Intersection {
	return intersectionsList[0]
}

type IntersectionNode struct {
	data *Intersection

	nextIntersections []*Intersection
}
