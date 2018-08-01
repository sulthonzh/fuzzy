package algorithm

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/charlesvdv/fuzmatch" // Find rule  matching with fuzzy
	"github.com/sulthonzh/match"
)

type (
	// Tsukamoto is a struct for ...
	Tsukamoto struct {
		Inputs        []Variable //
		Output        Variable   //
		Rules         []string   //
		AlphaPredikat []float64  //
		Zx            []float64  //
	}
	// Variable is a struct for ...
	Variable struct {
		Denomination string    // Variable Name
		Name         string    // Variable Name
		Mode         int       // 0:Fix 1:Min-Max
		Domains      []string  // { "Dekat" | "Sedang" | "Jauh" } | { "Baru" | "Sedang" | "Lama" }
		Data         float64   // Input | Output(From spesific periodic)
		Min          float64   // Min Data
		Max          float64   // Max Data
		MIU          []float64 // miu
		Graph        string
	}
)

// New is a func for ...
func (t *Tsukamoto) New() {
	t.Inputs = []Variable{}
	t.Output = Variable{}
}

// Calc is a func for ...
func (t *Tsukamoto) Calc(inputs map[string]float64) {
	for i, v := range t.Inputs {
		// Data input
		v.Data = inputs[v.Name]
		// Hitung miu
		v.calcMIU()
		// Update data
		t.Inputs[i] = v
	}

	t.calcAlphaPredicate()
	t.calcZx()
	t.Output.Data = t.calcZ()

}

// CalcZx is a func for ...
func (t *Tsukamoto) calcZx() {
	// Store miu to data interface for combine processing
	iData := make([][]match.AnyType, len(t.Inputs))
	for i, v := range t.Inputs {
		nd := make([]match.AnyType, len(v.Domains))

		for j, d := range v.Domains {
			nd[j] = fmt.Sprintf("%v_%v", v.Name, d)
		}

		iData[i] = nd
	}

	if len(iData) == 0 {
		return
	}

	// Combine miu
	c := [][]match.AnyType{}
	match.Combine(func(row []match.AnyType) {
		cTemp := make([]match.AnyType, len(row))
		copy(cTemp, row)
		c = append(c, cTemp)
	}, iData...)

	if len(c) == 0 {
		return
	}

	// Get string rule
	f := []string{}
	for _, g := range c {
		d := make([]string, len(g))
		for i, h := range g {
			d[i] = h.(string)
		}
		f = append(f, strings.Join(d, " "))
	}

	if len(f) == 0 {
		return
	}

	// Hasil validasi rules
	rulesResult := make([]string, len(f))
	s := make(map[int]string, len(f))

	r, _ := regexp.Compile("==(.*)$")
	for i, q := range f {
		for _, r := range t.Rules {
			// Fuzzy string matching for search rule
			s[fuzmatch.TokenSetRatio(q, r)] = r
		}
		rulesResult[i] = strings.TrimSpace(r.FindStringSubmatch(s[100])[1])
	}

	if len(rulesResult) == 0 {
		return
	}

	// Menentukan titik harga setiap rule
	rulesPoint := make([]int, len(f))
	for h, r := range rulesResult {
		for i, d := range t.Output.Domains {
			if r == fmt.Sprintf("%v_%v", t.Output.Name, d) {
				rulesPoint[h] = i
				break
			}
		}
	}
	if len(rulesPoint) == 0 {
		return
	}

	z := make([]float64, len(f))
	for i, p := range rulesPoint {
		z[i] = t.Output.z(i, t.AlphaPredikat[i], p)
	}

	t.Zx = z

}

// CalcZ is a function for calculate Z => zigma[z*a]/zigma[a]
func (t *Tsukamoto) calcZ() (result float64) {
	var sumZA float64
	var sumA float64
	for i := 0; i < len(t.Zx); i++ {
		sumZA += t.Zx[i] * t.AlphaPredikat[i] //zigma[z*a]
		sumA += t.AlphaPredikat[i]            //zigma[a]
	}
	return sumZA / sumA //zigma[z*a]/zigma[a]
}

// calcAlphaPredicate is a func for
func (t *Tsukamoto) calcAlphaPredicate() {
	// Store miu to data interface for combine processing
	iData := make([][]match.AnyType, len(t.Inputs))
	for a, b := range t.Inputs {
		c := make([]match.AnyType, len(b.MIU))

		for i, f := range b.MIU {
			c[i] = f
		}
		iData[a] = c
	}
	if len(iData) == 0 {
		return
	}

	// Combine miu
	e := [][]match.AnyType{}
	match.Combine(func(row []match.AnyType) {
		eTemp := make([]match.AnyType, len(row))
		copy(eTemp, row)
		e = append(e, eTemp)
	}, iData...)

	if len(e) == 0 {
		return
	}

	// Get smallest value from miu combination
	f := []float64{}
	for _, g := range e {
		d := make([]float64, len(g))
		for i, h := range g {
			d[i] = h.(float64)
		}

		if len(d) > 0 {
			sort.Float64s(d)
			f = append(f, d[0])
		}
	}

	t.AlphaPredikat = f
}

// CalcMIU is a func for
func (v *Variable) calcMIU() {
	if v.Mode == 1 {
		dl := len(v.Domains) // 6
		if dl > 2 {
			lenght := v.Max - v.Min            // 1000 - 500 = 500
			distance := lenght / float64(dl-1) // 500 / 5 = 100
			points := make([]float64, dl)
			miu := make([]float64, dl)

			for i := 0; i < dl; i++ {
				points[i] = v.Min + (float64(i) * distance) // 500 + (0 * 100) = 500 => (500,600,700,800,900,1000)
			}

			for i := 0; i < dl; i++ {
				if i == 0 {
					if v.Data <= points[i] {

						miu[i] = 1

						break
					} else if (v.Data > points[i]) && (v.Data < points[i+1]) {

						miu[i] = (points[i+1] - v.Data) / (points[i+1] - points[i])
						miu[i+1] = (v.Data - points[i]) / (points[i+1] - points[i])

						break
					}
				} else if i == dl-1 {
					miu[i] = 1
					break
				} else {
					if (v.Data > points[i]) && (v.Data < points[i+1]) {

						miu[i] = (points[i+1] - v.Data) / (points[i+1] - points[i])
						miu[i+1] = (v.Data - points[i]) / (points[i+1] - points[i])

						break
					} else if v.Data == points[i] {
						miu[i] = 1
						break
					}
				}

			}

			v.MIU = miu
		}
	}
}

// Z is function for ...
func (v *Variable) z(apIndex int, apValue float64, vPoint int) (result float64) {
	dl := len(v.Domains) // 6
	if dl > 2 {
		lenght := v.Max - v.Min            // 1000 - 500 = 500
		distance := lenght / float64(dl-1) // 500 / 5 = 100
		points := make([]float64, dl)

		for i := 0; i < dl; i++ {
			points[i] = v.Min + (float64(i) * distance) // 500 + (0 * 100) = 500 => (500,600,700,800,900,1000)
		}

		for i := 0; i < dl; i++ {
			if i == vPoint {
				if i == 0 {
					// Rumus mencari z jika posisi rule berada pada titik paling awal
					result = points[i+1] - (apValue * (points[i+1] - points[i]))
				} else if i == dl-1 {
					// Rumus mencari z jika posisi rule berada pada titik paling akhir
					result = points[i-1] + (apValue * (points[i] - points[i-1]))
				} else {
					// Rumus mencari z jika posisi rule berada pada titik selain dari keduanya
					result = points[i]
				}
				break
			}
		}

	}
	return result
}
