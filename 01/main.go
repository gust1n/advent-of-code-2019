package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const input = `97587
100963
85693
127587
147839
134075
128598
91290
119100
138824
56295
132118
105018
143092
89032
104836
138278
60416
62570
142110
55179
80891
99106
122863
108894
112654
117175
96093
76214
65412
124388
66465
115850
147531
87643
75882
62912
76100
102120
83803
139304
139325
126412
145152
136247
68246
130156
59097
79024
62480
121847
54739
118690
116247
117283
144827
147562
126796
148210
109099
98831
59412
141077
121786
142878
140144
57855
59571
118451
149097
145088
76882
53732
70543
89874
114366
115683
99139
108440
76964
134451
109250
66021
132683
149013
122917
137810
108451
109606
94396
106926
100901
108587
99847
64257
147162
133698
140775
129466
72487`

func main() {
	inputLines := strings.Split(input, "\n")
	var modules []int

	for _, inputLine := range inputLines {
		i, err := strconv.Atoi(inputLine)
		if err != nil {
			log.Fatal(err)
		}
		modules = append(modules, i)
	}

	var totalMass int

	for _, module := range modules {
		totalMass = totalMass + calculateMass(0, module)
	}

	fmt.Println(totalMass)
}

func calculateMass(totalModuleMass int, module int) int {
	mass := (module / 3) - 2
	if mass <= 0 {
		return totalModuleMass
	}
	return calculateMass(totalModuleMass+mass, mass)
}
