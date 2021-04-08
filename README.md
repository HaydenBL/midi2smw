# midi2smw

[![ci-test](https://github.com/HaydenBL/midi2smw/actions/workflows/ci.yml/badge.svg)](https://github.com/HaydenBL/midi2smw/actions/workflows/ci.yml)

### Overview

This project's aim is to simplify the process of porting a song into Super Mario World with AddMusicK.
It reads in a midi file, creates track timelines of note events, quantizes them down to the nearest
64th note (smallest note value SMW handles), and then prints them out in a format readable by AddMusicK
(including track numbers and BPM, converted to SMW tempo). The only things you might have to modify
yourself are the channel volumes, and removing tracks if it exports more than the 8 max.

It also keeps in mind the start/end time of other tracks by padding out empty space on either side
of the track with rests so all channels will be properly synced up when played together.

The project has no third-party dependencies and is run as a plain executable.

### Features

* Overlapping notes/chords in a track will be output to separate channels, so no notes are lost. If 
  multiples notes are playing in a track at the same time, priority is resolved by whichever note
  started playing first, and if the notes started playing at the same time, the note with the higher 
  key will be exported to the track first.
* Samples can be specified for certain notes in a midi track. This is particularly useful for drum tracks,
  where different notes are used for different drum sounds.
* Tracks can be split by specifying specific notes in a track to split into a separate track. Again, this
  would most likely be useful for a drum track where you want to remove a certain sound from the channel.
* By default, output is be compressed with loops for repeating sections (no guarantee it's 100% optimal)
  but this can be disabled. It will also put long rests into loops so they're not strung together with
  a bunch of ties.

### Usage
`midi2smw.exe <flags> inputMidi.mid`  
To specify an output file:  
`midi2smw.exe <flags> inputMidi.mid outputFile.txt`

Flags must be put before the input/output file names.  
None of the flags have arguments, simply put them in the command, and the cli will walk you through using them.

| Flag            | Description                                               |
| --------------- | --------------------------------------------------------- |
|`-specifyTracks` | Manually specify which tracks to insert into the output   |
|`-split`         | Specify tracks to split with note groupings               |
|`-samples`       | Specify samples for notes                                 |
|`-noLoop`        | Print output without loops                                |

### Limitations

The program doesn't currently support triplets. I haven't even looked into how those work in AMK,
and it's possible it would take a somewhat large refactor to implement that.

Overall, it's not a super intelligent algorithm, and since it's just rounding to the nearest 64th note,
it's possible the output might be somewhat ugly. In cases like that I'd recommend modifying the midi file
itself to fit better into SMW's limitations. Generally, output may just not be optimal for all of AddMusicK's
features. I've never read all its documentation, and in general I don't know a ton about music porting.

Hopefully it's useful to someone anyway!

### Special Thanks
* [OneLoneCoder's video on midi parsing](https://youtu.be/040BKtnDdg0)
* [Wakana's Music Porting Tutorial](https://www.smwcentral.net/?p=viewthread&t=89606)
