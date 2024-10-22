package util

func GetMinioInstanceFromId(id string, instanceCount int) int {
	hash := int32(0)
	for _, c := range id {
		hash += int32(c)
	}
	//To ensure that index is non-negative
	index := (int(hash)%instanceCount + instanceCount) % instanceCount
	return index
}
