http:
    addr: ":{{ .Values.service.internalPort }}"

storage:
    enabled: {{ .Values.config.storage.enabled }}
    bucket_url: "{{ .Values.config.storage.bucketUrl }}"

cache:
    base_dir: "{{ .Values.config.cache.dir }}"

repositories:
    - url: "https://repo1.maven.org/maven2"
    - url: "http://uk.maven.org/maven2/"
