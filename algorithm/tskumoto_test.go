package algorithm

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFuzzyTsukamoto(t *testing.T) {
	// Exmple with 2 variables and @3 domains | 9 rules
	f := Tsukamoto{

		Inputs: []Variable{
			Variable{
				Name:    "Kilometer",
				Mode:    1,
				Min:     6000,
				Max:     82000,
				Domains: []string{"Dekat", "Sedang", "Jauh"},
			},
			Variable{
				Name:    "Tahun",
				Mode:    1,
				Min:     1,
				Max:     8,
				Domains: []string{"Baru", "Sedang", "Lama"},
			},
		},
		Rules: []string{
			"Kilometer_Jauh {{And}} Tahun_Baru == Harga_Murah ",
			"Kilometer_Jauh {{And}} Tahun_Sedang == Harga_Murah",
			"Kilometer_Jauh {{And}} Tahun_Lama == Harga_Murah",
			"Kilometer_Sedang {{And}} Tahun_Baru == Harga_Sedang",
			"Kilometer_Sedang {{And}} Tahun_Sedang == Harga_Sedang",
			"Kilometer_Sedang {{And}} Tahun_Lama == Harga_Murah",
			"Kilometer_Dekat {{And}} Tahun_Baru == Harga_Mahal",
			"Kilometer_Dekat {{And}} Tahun_Sedang == Harga_Sedang",
			"Kilometer_Dekat {{And}} Tahun_Lama == Harga_Murah",
		},
		Output: Variable{
			Name:    "Harga",
			Mode:    1,
			Min:     6250000,
			Max:     14500000,
			Domains: []string{"Murah", "Sedang", "Mahal"},
		},
	}
	data := []map[string]float64{
		{"Kilometer": 10000, "Tahun": 8},
		{"Kilometer": 15000, "Tahun": 7},
		{"Kilometer": 20000, "Tahun": 6},
		{"Kilometer": 25000, "Tahun": 5},
		{"Kilometer": 30000, "Tahun": 4},
		{"Kilometer": 35000, "Tahun": 3},
		{"Kilometer": 40000, "Tahun": 2},
		{"Kilometer": 60000, "Tahun": 1},
	}

	zs := make([]float64, len(data))

	for i, d := range data {
		f.Calc(d)

		zs[i] = f.Output.Data
	}

	for _, d := range zs {
		fmt.Printf("Estimasi harga: %.2f\n", d)
	}
	j, _ := json.MarshalIndent(f, "", "  ")
	fmt.Printf("\"Data\": %v", string(j))
}
