package main

import (
	"strconv"
	"strings"
)

func validateXADDKey(streams map[string][]Stream, key, id string) (bool, error) {
	streamList := streams[key]
	currSplitParts := strings.Split(id, "-")
	currMiliSeconds, err := strconv.Atoi(currSplitParts[0])
	if err != nil {
		return false, err
	}
	currSequence, err := strconv.Atoi(currSplitParts[1])
	if err != nil {
		return false, err
	}
	if len(streamList) > 0 {
		prevStream := streamList[len(streamList)-1]
		prevSplitParts := strings.Split(prevStream.StreamID, "-")
		prevMiliSeconds, err := strconv.Atoi(prevSplitParts[0])
		if err != nil {
			return false, err
		}
		prevSequence, err := strconv.Atoi(prevSplitParts[1])
		if err != nil {
			return false, err
		}

		if currMiliSeconds > prevMiliSeconds {
			return true, nil
		} else if currMiliSeconds == prevMiliSeconds {
			return currSequence > prevSequence, nil
		} else {
			return false, nil
		}
	}
	miliSecondsNotZero := currMiliSeconds >= 0
	currSeqNotZero := currSequence > 0
	return miliSecondsNotZero && currSeqNotZero, nil
}
