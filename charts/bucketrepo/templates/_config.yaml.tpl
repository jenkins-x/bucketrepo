http:
    addr: ":{{ .Values.service.internalPort }}"
    username: "{{ .Values.config.auth.username }}"
    password: "{{ .Values.config.auth.password }}"

storage:
    enabled: {{ .Values.config.storage.enabled }}
    bucket_url: "{{ .Values.config.storage.bucketUrl }}"

cache:
    base_dir: "{{ .Values.config.cache.dir }}"

repositories:
{{- if .Values.config.repositories }}
{{- range $key, $value := .Values.config.repositories }}
  {{- if $value }}
    - url: {{ $value | quote }}
  {{- end }}
{{- end }}
{{- end }}
