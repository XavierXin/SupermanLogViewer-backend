package helpers

func getMinFromNumArray(a []int) (min int, minIndex int) {
	if len(a) == 0 {
		return 0, 0
	}
	min = a[0]
	minIndex = 0
	for i, num := range a {
		if num < min {
			min = num
			minIndex = i
		}
	}
	return
}

// GenerateSplitViewList will generate an array of SplitViewEntry for front end
// to display
func GenerateSplitViewList(daemonsLog []*LogEntryList) []*SplitViewEntry {
	numOfDaemonsToSplit := len(daemonsLog)
	totalEntryCount := 0
	for _, eachDaemonEntryList := range daemonsLog {
		totalEntryCount += len(eachDaemonEntryList.List)
	}

	currentDaemonEntryIndex := make([]int, numOfDaemonsToSplit)

	res := make([]*SplitViewEntry, totalEntryCount)

	indexToCompare := make([]int, numOfDaemonsToSplit)
	for i := range res {
		for j := 0; j < numOfDaemonsToSplit; j++ {
			if currentDaemonEntryIndex[j] >= len(daemonsLog[j].List) {
				indexToCompare[j] = 2147483647
			} else {
				indexToCompare[j] = daemonsLog[j].List[currentDaemonEntryIndex[j]].AbsoluteIndex
			}
		}
		_, daemonIndexForThisRow := getMinFromNumArray(indexToCompare)
		res[i] = &SplitViewEntry{
			LogEntry:    daemonsLog[daemonIndexForThisRow].List[currentDaemonEntryIndex[daemonIndexForThisRow]],
			ColumnIndex: daemonIndexForThisRow,
		}
		// then update currentDaemonEntryIndex
		if len(daemonsLog[daemonIndexForThisRow].List) >= currentDaemonEntryIndex[daemonIndexForThisRow]+1 {
			currentDaemonEntryIndex[daemonIndexForThisRow]++
		}
	}

	return res
}
