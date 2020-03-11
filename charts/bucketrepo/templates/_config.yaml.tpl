http:
    addr: ":{{ .Values.service.internalPort }}"
    username: "{{ .Values.config.auth.username | default .Values.secrets.adminUser.username }}"
    password: "{{ .Values.config.auth.password | default .Values.secrets.adminUser.password }}"
    chartPath: "{{ .Values.config.charts.path}}"

storage:
{{- if .Values.config.storage.bucketUrl .Values.jxRequirements.storage.repository.url }}
    enabled: true
{{- end }}
    bucket_url: "{{ .Values.config.storage.bucketUrl | default .Values.jxRequirements.storage.repository.url }}"

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
