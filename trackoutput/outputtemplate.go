package trackoutput

const outputTemplate = `#amk 2

#spc
{
  #comment "Generated with smw2midi"
}

w255
{{range $i, $t := .ChannelOutputs}}
; {{.Name}}
#{{$i}} t{{$.Bpm}}
@{{.DefaultSample}} v150 o{{.StartOctave}}
{{.NoteOutput}}
{{end}}`
