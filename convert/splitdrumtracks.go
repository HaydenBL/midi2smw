package convert

import (
	"fmt"
	"midi2smw/drumtrack"
	"sort"
)

func splitDrumTracks(tracks []noteTrack, drumTrackGroups []drumtrack.Group) []noteTrack {
	// maps a track number to a list of newly split note tracks
	noteSplitMap := make(map[uint8][]noteTrack)
	for i, dtg := range drumTrackGroups {
		trackToSplit := tracks[i]
		noteSplitMap[dtg.TrackNumber] = make([]noteTrack, len(dtg.NoteGroups)+1)
		for _, note := range trackToSplit.Notes {
			index := getIndexForDrumTrackGroup(note, dtg)
			noteSplitMap[dtg.TrackNumber][index].Notes = append(noteSplitMap[dtg.TrackNumber][index].Notes, note)
		}
	}

	splitTrackNums := make([]uint8, 0)
	for trackNum := range noteSplitMap {
		splitTrackNums = append(splitTrackNums, trackNum)
	}
	sort.Slice(splitTrackNums, func(i, j int) bool {
		return splitTrackNums[i] > splitTrackNums[j]
	})
	// loop over tracks backwards and insert split tracks in their place
	for _, trackNum := range splitTrackNums {
		tracksToInsert := noteSplitMap[trackNum]
		tracksToInsert = setSplitTrackNames(tracks[trackNum].Name, tracksToInsert)
		tracks = insertTracksAt(tracks, tracksToInsert, int(trackNum))
	}
	return tracks
}

func insertTracksAt(tracks []noteTrack, tracksToInsert []noteTrack, index int) []noteTrack {
	tail := make([]noteTrack, len(tracks[index+1:]))
	copy(tail, tracks[index+1:]) // save tail
	endPadding := make([]noteTrack, len(tracksToInsert)-1)
	tracks = append(tracks, endPadding...)         // pad out end with of slice
	copy(tracks[index:], tracksToInsert)           // copy new tracks from index onward
	copy(tracks[index+len(tracksToInsert):], tail) // reinsert tail at the end
	return tracks
}

func setSplitTrackNames(oldName string, tracks []noteTrack) []noteTrack {
	for i := range tracks {
		tracks[i].Name = fmt.Sprintf("%s - Split %d", oldName, i+1)
	}
	return tracks
}
