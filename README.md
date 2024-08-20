# metronome

It is a metronome.

There are tons of web/etc metronomes out there but for some reason it's very common to have a lower bound on bpm, like 40 or so.

One way I like to practice is with a click/beep every beat one which is hard for slower tempos with that lower bound.

`metronome -bpm 15` lets me do that for a 60 bpm tempo.

## Install

``` shell
go install github.com/danp/metronome@latest
```

This program uses [oto/v2](https://pkg.go.dev/github.com/hajimehoshi/oto/v2) which does not requrire cgo for macOS or Windows.
See docs for what's necessary on other platforms.
