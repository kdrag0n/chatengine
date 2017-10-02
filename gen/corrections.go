package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Generating correction data...")

	lineSeg1 := []byte("\n\t\"")
	lineSeg2 := []byte(`": []string{`)
	lineSeg3 := []byte(`",`)
	lineSeg4 := []byte(`},`)

	byteIn, err := ioutil.ReadFile("data/corrections.txt")
	check(err)

	output := bytes.NewBufferString(`package main

var correctionData = map[string][]string{`)
	lines := bytes.FieldsFunc(byteIn, func(r rune) bool { return r == '\n' })

	for _, line := range lines {
		arrowDashIndex := bytes.LastIndexByte(line, '-')
		from := line[:arrowDashIndex]
		toList := line[arrowDashIndex+2:] // ? need : len()?

		output.Write(lineSeg1)
		output.Write(from)
		output.Write(lineSeg2)

		for _, to := range bytes.FieldsFunc(toList, func(r rune) bool { return r == ',' }) {
			output.WriteByte('"')
			output.Write(to)
			output.Write(lineSeg3)
		}

		output.Write(lineSeg4)
	}
	output.WriteString("\n}")

	err = ioutil.WriteFile("gen_correction_data.go", output.Bytes(), 644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Correction data generated.")
}
