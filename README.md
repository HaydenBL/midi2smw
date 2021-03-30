# midi2smw

[![ci-test](https://github.com/HaydenBL/midi2smw/actions/workflows/ci.yml/badge.svg)](https://github.com/HaydenBL/midi2smw/actions/workflows/ci.yml)

### Overview

This project's aim is to simplify the process of porting a song into Super Mario World via AddMusicK.
It reads in a midi file, parses the events, gets out the important ones (note on/off), converts them
into track timelines of note events, quantizes them down to the nearest 64th note (smallest note
value SMW handles), and then prints them out in a format readable by AddMusicK. It also has some
handling for chords/overlapping notes, by passing over a track multiple times if necessary,
creating separate channel outputs so every note should be exported. The priority is resolved by
whichever note started playing first, and if the notes started playing at the same time, the note
with the higher key will be exported to the track first.

It also keeps in mind the start/end time of other tracks by padding out empty space on either side
of the track with rests so all channels will be properly synced up when played together.

Currently, the project has no third-party dependencies, so it's pretty easily to get up and running
locally if you want to tweak something yourself.

### TODO
* Optimize outputs (loops, dotted notes)

### Special Thanks
* [OneLoneCoder's video on midi parsing](https://youtu.be/040BKtnDdg0)
* [Wakana's Music Porting Tutorial](https://www.smwcentral.net/?p=viewthread&t=89606)
