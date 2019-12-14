package resourceAllocation

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type resourceAllocationTask struct {
	// N
	countSuppliers int
	// M
	countCastomers int
	// T
	countTacts int

	// a1, a2, ..., aN
	volumeGoodsFromSuppliers []int

	// b11, b12, ..., b1T
	// ...
	// bN1, bN2, ..., bNT
	volumeGoodsFromSuppliersOnTact [][]int

	// c11, c12, ..., c1T
	// ...
	// cM1, cM2, ..., cMT
	volumeUsedGoodsByCastomersOnTact [][]int

	// D (D1, D2, ..., DM)
	setsSuppliersForCastomers [][]int
}

func MakeResourceAllocationTaskFromFile(taskFile string) *resourceAllocationTask {
	content, err := ioutil.ReadFile(taskFile)
	if err != nil {
		log.Fatal(err)
	}
	contentLines := strings.Split(string(content), "\r\n")

	countSuppliers, err := strconv.Atoi(contentLines[0])
	if err != nil {
		log.Fatal(err)
	}

	countCastomers, err := strconv.Atoi(contentLines[1])
	if err != nil {
		log.Fatal(err)
	}

	countTacts, err := strconv.Atoi(contentLines[2])
	if err != nil {
		log.Fatal(err)
	}

	volumeGoodsFromSuppliersString := strings.Trim(contentLines[4], " ")
	volumeGoodsFromSuppliersStrings := strings.Split(volumeGoodsFromSuppliersString, " ")
	volumeGoodsFromSuppliers := make([]int, countSuppliers)
	for i, volumeGood := range volumeGoodsFromSuppliersStrings {
		volumeGoodsFromSuppliers[i], err = strconv.Atoi(volumeGood)
		if err != nil {
			log.Fatal(err)
		}
	}

	// allocation B
	volumeGoodsFromSuppliersOnTact := make([][]int, countSuppliers)
	for i := 0; i < countSuppliers; i++ {
		volumeGoodsFromSuppliersOnTact[i] = make([]int, countTacts)
	}

	// initialization B
	for i, contentLine := range contentLines[6 : 6+countSuppliers] {
		trimmedLine := strings.Trim(contentLine, " ")
		splittedVolumesStrings := strings.Split(trimmedLine, " ")
		for j, volumeString := range splittedVolumesStrings {
			volumeGoodsFromSuppliersOnTact[i][j], err = strconv.Atoi(volumeString)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// allocation C
	volumeUsedGoodsByCastomersOnTact := make([][]int, countCastomers)
	for i := 0; i < countCastomers; i++ {
		volumeUsedGoodsByCastomersOnTact[i] = make([]int, countTacts)
	}

	// initialization C
	startCastomersLine := 6 + countSuppliers + 1
	for i, contentLine := range contentLines[startCastomersLine : startCastomersLine+countCastomers] {
		trimmedLine := strings.Trim(contentLine, " ")
		splittedVolumesStrings := strings.Split(trimmedLine, " ")
		for j, volumeString := range splittedVolumesStrings {
			volumeUsedGoodsByCastomersOnTact[i][j], err = strconv.Atoi(volumeString)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// allocation D
	startSetsSuppliersLine := startCastomersLine + countCastomers + 1
	setsSuppliersForCastomers := make([][]int, countCastomers)
	for i := 0; i < countCastomers; i++ {
		setsSuppliersForCastomers[i] = []int{}
	}

	// initialization D
	for i, contentLine := range contentLines[startSetsSuppliersLine : startSetsSuppliersLine+countCastomers] {
		trimmedLine := strings.Trim(contentLine, " ")
		splittedSuppliersStrings := strings.Split(trimmedLine, " ")
		for _, supplierString := range splittedSuppliersStrings {
			supplier, err := strconv.Atoi(supplierString)
			if err != nil {
				log.Fatal(err)
			}
			setsSuppliersForCastomers[i] = append(setsSuppliersForCastomers[i], supplier)
		}
	}

	return &resourceAllocationTask{
		countSuppliers:                   countSuppliers,
		countCastomers:                   countCastomers,
		countTacts:                       countTacts,
		volumeGoodsFromSuppliers:         volumeGoodsFromSuppliers,
		volumeGoodsFromSuppliersOnTact:   volumeGoodsFromSuppliersOnTact,
		volumeUsedGoodsByCastomersOnTact: volumeUsedGoodsByCastomersOnTact,
		setsSuppliersForCastomers:        setsSuppliersForCastomers,
	}

}
