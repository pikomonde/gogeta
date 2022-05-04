package gogetautil

// SliceIntOrderedFindIdx returns the index of match the searched val, if not found it will returns the index of
// the val closest and greater to the searched val. The input array MUST be SORTED.
// TODO: create unit test
func SliceIntOrderedFindIdx(arr []int, val int) (int, bool) {
	start := 0
	end := len(arr) - 1

	for start <= end {
		mid := (start + end) / 2

		if val == arr[mid] {
			return mid, true
		} else if val < arr[mid] {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}
	return end, false
}

func SliceStringFindIdx(arr []string, val string) int {
	for i, v := range arr {
		if val == v {
			return i
		}
	}
	return -1
}

func SliceInsert(arr []int, idx int, val int) []int {
	return append(arr[:idx+1], append([]int{val}, arr[idx+1:]...)...)
}

func SliceCut(arr []int, idx int) []int {
	return append(arr[:idx], arr[idx+1:]...)
}

func SliceCutZeros(arr []int) []int {
	return sliceCutZerosWithoutAppend(arr)
}

// sliceCutZerosWithAppend slightly faster, but can be exponentially slower in some cases. It might be
// because of the calling of append function repeatedly.
func sliceCutZerosWithAppend(arr []int) []int {
	for idx := len(arr) - 1; idx >= 0; idx-- {
		if arr[idx] == 0 {
			arr = SliceCut(arr, idx)
		}
	}
	return arr
}

// sliceCutZerosWithoutAppend should be faster. TODO: unit test
func sliceCutZerosWithoutAppend(arr []int) []int {
out:
	for idx, idxToSwitch := 0, 1; idx < len(arr); idx++ {
		if arr[idx] == 0 {
			for {
				if idxToSwitch > len(arr)-1 {
					arr = arr[:idx]
					break out
				}
				if arr[idxToSwitch] == 0 {
					idxToSwitch++
					continue
				}
				break
			}

			// switch
			arr[idxToSwitch], arr[idx] = arr[idx], arr[idxToSwitch]
			idxToSwitch++
		}
	}
	return arr
}
