package d2cof

import (
	"strings"

	"github.com/OpenDiablo2/D2Shared/d2common/d2interface"

	"github.com/OpenDiablo2/D2Shared/d2common/d2enum"

	"github.com/OpenDiablo2/D2Shared/d2common"
)

type COF struct {
	NumberOfDirections int
	FramesPerDirection int
	NumberOfLayers     int
	CofLayers          []CofLayer
	CompositeLayers    map[d2enum.CompositeType]int
	AnimationFrames    []d2enum.AnimationFrame
	Priority           [][][]d2enum.CompositeType
}

func LoadCOF(fileName string, fileProvider d2interface.FileProvider) *COF {
	result := &COF{}
	fileData := fileProvider.LoadFile(fileName)
	if len(fileData) == 0 {
		return result
	}
	streamReader := d2common.CreateStreamReader(fileData)
	result.NumberOfLayers = int(streamReader.GetByte())
	result.FramesPerDirection = int(streamReader.GetByte())
	result.NumberOfDirections = int(streamReader.GetByte())
	streamReader.SkipBytes(25) // Skip 25 unknown bytes...
	result.CofLayers = make([]CofLayer, result.NumberOfLayers)
	result.CompositeLayers = make(map[d2enum.CompositeType]int, 0)
	for i := 0; i < result.NumberOfLayers; i++ {
		layer := CofLayer{}
		layer.Type = d2enum.CompositeType(streamReader.GetByte())
		layer.Shadow = streamReader.GetByte()
		streamReader.SkipBytes(1) // Unknown
		layer.Transparent = streamReader.GetByte() != 0
		layer.DrawEffect = d2enum.DrawEffect(streamReader.GetByte())
		weaponClassStr, _ := streamReader.ReadBytes(4)
		layer.WeaponClass = d2enum.WeaponClassFromString(strings.TrimSpace(strings.ReplaceAll(string(weaponClassStr), string(0), "")))
		result.CofLayers[i] = layer
		result.CompositeLayers[layer.Type] = i
	}
	animationFrameBytes, _ := streamReader.ReadBytes(result.FramesPerDirection)
	result.AnimationFrames = make([]d2enum.AnimationFrame, result.FramesPerDirection)
	for i := range animationFrameBytes {
		result.AnimationFrames[i] = d2enum.AnimationFrame(animationFrameBytes[i])
	}
	priorityLen := result.FramesPerDirection * result.NumberOfDirections * result.NumberOfLayers
	result.Priority = make([][][]d2enum.CompositeType, result.NumberOfDirections)
	priorityBytes, _ := streamReader.ReadBytes(priorityLen)
	priorityIndex := 0
	for direction := 0; direction < result.NumberOfDirections; direction++ {
		result.Priority[direction] = make([][]d2enum.CompositeType, result.FramesPerDirection)
		for frame := 0; frame < result.FramesPerDirection; frame++ {
			result.Priority[direction][frame] = make([]d2enum.CompositeType, result.NumberOfLayers)
			for i := 0; i < result.NumberOfLayers; i++ {
				result.Priority[direction][frame][i] = d2enum.CompositeType(priorityBytes[priorityIndex])
				priorityIndex++
			}
		}
	}
	return result
}
