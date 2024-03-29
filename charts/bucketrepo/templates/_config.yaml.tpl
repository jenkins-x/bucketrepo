http:
    addr: ":{{ .Values.service.internalPort }}"
    username: "{{ .Values.config.auth.username | default .Values.secrets.adminUser.username }}"
    password: "{{ .Values.config.auth.password | default .Values.secrets.adminUser.password }}"
    chartPath: "{{ .Values.config.charts.path}}"

storage:
{{- if .Values.config.storage.bucketUrl }}
    enabled: true
    bucket_url: "{{ .Values.config.storage.bucketUrl }}"
{{- else if and (hasKey .Values.jxRequirements "storage") ( .Values.jxRequirements.storage)  }}
{{- range $key, $val := .Values.jxRequirements.storage }}
{{- if eq "repository" $val.name }}
    enabled: true
    bucket_url: "{{ $val.url }}"
{{- end }}
{{- end }}
{{- else }}
    bucket_url: ""
{{- end }}

cache:
    base_dir: "{{ .Values.config.cache.dir }}"

repositories:
{{- if .Values.config.repositories }}
{{- range $key, $value := .Values.config.repositories }}
{{- if $value }}
{{- if kindIs "string" $value }}
    - url: {{ $value | quote }}
{{- else }}
    - {{ toYaml $value | nindent 7 }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
