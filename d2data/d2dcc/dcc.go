package d2dcc

import (
	"log"

	"github.com/OpenDiablo2/D2Shared/d2common/d2interface"

	"github.com/OpenDiablo2/D2Shared/d2common"
)

type DCC struct {
	Signature          int
	Version            int
	NumberOfDirections int
	FramesPerDirection int
	Directions         []DCCDirection
	valid              bool
}

func (v DCC) IsValid() bool {
	return v.valid
}

func LoadDCC(path string, fileProvider d2interface.FileProvider) DCC {
	result := DCC{}
	fileData := fileProvider.LoadFile(path)
	if len(fileData) == 0 {
		ret := DCC{}
		ret.valid = false
		return ret
	}
	var bm = d2common.CreateBitMuncher(fileData, 0)
	result.Signature = int(bm.GetByte())
	if result.Signature != 0x74 {
		log.Fatal("Signature expected to be 0x74 but it is not.")
	}
	result.Version = int(bm.GetByte())
	result.NumberOfDirections = int(bm.GetByte())
	result.FramesPerDirection = int(bm.GetInt32())
	if bm.GetInt32() != 1 {
		log.Fatal("This value isn't 1. It has to be 1.")
	}
	bm.GetInt32() // TotalSizeCoded
	directionOffsets := make([]int, result.NumberOfDirections)
	for i := 0; i < result.NumberOfDirections; i++ {
		directionOffsets[i] = int(bm.GetInt32())
	}
	result.Directions = make([]DCCDirection, result.NumberOfDirections)
	for i := 0; i < result.NumberOfDirections; i++ {
		result.Directions[i] = CreateDCCDirection(d2common.CreateBitMuncher(fileData, directionOffsets[i]*8), result)
	}
	result.valid = true
	return result
}
