package somesql

import "text/template"

var (
	selectTplStr = `SELECT{{ " " -}}
	{{- range $i, $v := .MetaFields -}}
		"{{ $v }}"
		{{- if ne (len $.MetaFields) (plus $i) }}, {{ end }}
	{{- end -}}
	{{- if and (ne (len .MetaFields) 0) (ne (len .DataFields) 0) }}, {{ end -}}
	{{- range $i, $v := .DataFields -}}
		{{ if $.Inner -}}
			"{{ $.FieldDataLang }}"->>'{{ $v }}' "{{ $v }}"
		{{- else -}}
			{{ if eq $i 0 }}json_build_object({{ end -}}
			'{{ $v }}', "{{ $.FieldDataLang }}"->'{{ $v }}'
			{{- if eq (len $.DataFields) (plus $i) }}) "{{ $.FieldData }}"{{ end }}
		{{- end }}
		{{- if ne (len $.DataFields) (plus $i) }}, {{end}}
	{{- end }} FROM repo
	{{- if ne (len .Conditions) 0 }} WHERE {{ .Conditions }}{{ end }}
	{{- if ne (.Query.GetLimit) 0 }} LIMIT {{ .Query.GetLimit }}{{ end }}
	{{- if ne (.Query.GetOffset) 0 }} OFFSET {{ .Query.GetOffset }}{{ end -}}
	`

	mergeTplStr = `INSERT INTO 
	
	`

	deleteTplStr = `DELETE FROM repo
	{{- if ne (len .Conditions) 0 }} WHERE {{ .Conditions }}{{ end }}
	{{- if ne (.Query.GetLimit) 0 }} LIMIT {{ .Query.GetLimit }}{{ end }}
	{{- if ne (.Query.GetOffset) 0 }} OFFSET {{ .Query.GetOffset }}{{ end -}}
	`

	funcMap = template.FuncMap{
		// The name "plus" is what the function will be called in the template text.
		"plus": func(i int) int {
			return i + 1
		},
	}
)
