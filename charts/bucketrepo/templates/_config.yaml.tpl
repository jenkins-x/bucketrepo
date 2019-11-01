http:
    addr: ":{{ .Values.service.internalPort }}"

storage:
    enabled: {{ .Values.config.storage.enabled }}
    bucket_url: "{{ .Values.config.storage.bucketUrl }}"

cache:
    base_dir: "{{ .Values.config.cache.dir }}"

repositories:
    - url: "https://repo.maven.org/maven2"
    - url: "https://repo1.maven.org/maven2"
    - url: "http://uk.maven.org/maven2/"
    - url: "https://repo.spring.io/release/"
    - url: "https://repo.spring.io/milestone/"
    - url: "https://services.gradle.org/distributions/"
    - url: "https://repo.jenkins-ci.org/public/"
    - url: "https://repo.jenkins-ci.org/releases/"
    - url: "https://jitpack.io/"
    - url: "https://repo.jenkins-ci.org/releases/"
    - url: "https://registry.npmjs.org/"
    - url: "https://plugins.gradle.org/m2/"
