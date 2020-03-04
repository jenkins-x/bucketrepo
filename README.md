[![MIT License][license-image]][license-url]

# Bucketrepo

This is a lightweight artifacts repository with low memory footprint which can be used as a minimal replacement for 
[Nexus](https://www.sonatype.com/nexus-repository-sonatype). It is able to cache artifacts from a remote repository
on a local filesystem volume and also to store them on a cloud storage bucket via [go-cloud](https://github.com/google/go-cloud/).

It can be deployed either as a side-car container to a Kubernetes build pod or as a standalone service.

## Getting Started

### Configuration

The [default configuration](config/config.yaml) enables only the local file cache:
```YAML
http:
    addr: ":8080"

storage:
    enabled: false
    bucket_url: "gs://bucketrepo"

cache:
    base_dir: "/tmp/bucketrepo"

repositories:
    - url: "https://repo1.maven.org/maven2"
    - url: "http://uk.maven.org/maven2/"
```

The header (for example Bearer token authentication) and timeout can be modified for every remote repository:
```YAML
...

repositories:
    - url: "https://repo1.maven.org/maven2"
      timeout: "30s"
    - url: "http://uk.maven.org/maven2/"
    - url: "http://my.private.maven.repo"
      timeout: "10s"
      header:
          Authorization: "Bearer <Token>"
```

The cloud storage can be enabled by providing a bucket URL:
```YAML
storage:
    enabled: true
    bucket_url: "gs://mybucket"
    # if necessary is possible to use Path prefix
    prefix: "myfolder" 
```

Also the TLS and basic authentication can be configured with:
```YAML
http:
    addr: ":8080"
    https: true
    crt: "/certs/domain.crt"
    key: "/certs/domain.key"
    username: "myuser"
    password: "mypassword"
```
Note that the basic authentication is turned off when HTTPS is disabled.

### Supported Artifacts

This repository has been tested with `maven` and `helm` tools, but it can also store other artifacts.

### Installing

#### Kubernetes

The repository service can be installed in a Kubernetes cluster using helm. First, you need to add the jenkins-x chart repository to your helm repositories:

```sh
helm repo add jenkins-x http://chartmuseum.jenkins-x.io
helm repo update
```
You can now install the chart with:
```
helm install jenkins-x/bucketrepo --name bucketrepo
```

#### Locally
The repository can be started in a docker container usinged the [latest released](https://github.com/jenkins-x/bucketrepo/releases) image:
```bash
docker run -v $(pwd)/config:/config -p 8080:8080 gcr.io/jenkinsxio/bucketrepo:0.1.12 -config-path=/config
```

Or it can be built and run with:
```bash
make build
bin/bucketrepo -config-path=config -log-level=debug
```

### Maven Configuration

`bucketrepo`can be configured as a mirror by adding the following in `~/.m2/settings.xml` file:
```XML
<settings>
<mirrors>
    <mirror>
      <id>bucketrepo</id>
      <name>bucketrepo- mirror</name>
      <url>http://localhost:8080/bucketrepo/</url>
      <mirrorOf>*</mirrorOf>
    </mirror>
  </mirrors>
</settings>
```

And as a repository by adding the following in the `pom.xml` file:
```XML
<repositories>
    <repository>
        <id>bucketrepo</id>
        <url>http://localhost:8080/bucketrepo/</url>
    </repository>
</repositories>

<distributionManagement>
    <snapshotRepository>
        <id>snapshots</id>
        <url>http://localhost:8080/bucketrepo/deploy/maven-snapshots/</url>
    </snapshotRepository>
    <repository>
        <id>releases</id>
        <url>http://localhost:8080/bucketrepo/deploy/maven-releases/</url>
    </repository>
</distributionManagement>
```

## Acknowledgments

This project is originally based on [nexus-minimal](https://github.com/atsman/nexus-minimal). Thank you [atsman](https://github.com/atsman) for creating that project.

## License

[MIT](LICENSE)

[license-url]: LICENSE

[license-image]: https://img.shields.io/github/license/mashape/apistatus.svg

[capture]: capture.png
