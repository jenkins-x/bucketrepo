[![MIT License][license-image]][license-url]
[![Build Status](https://travis-ci.org/Atsman/nexus-minimal.svg?branch=master)](https://travis-ci.org/Atsman/nexus-minimal)
[![Go Report Card](https://goreportcard.com/badge/github.com/Atsman/nexus-minimal)](https://goreportcard.com/report/github.com/Atsman/nexus-minimal)

# Nexus Minimal (aka Lil Nexus)

Nexus minimal is an implementation of nexus repo in golang. Decision to implement new tiny repo come to me when i realized that nexus requires 4gb RAM machine at minimum. (Java world :)) Proof here: [Nexus memory requrements](https://help.sonatype.com/display/NXRM3/System+Requirements#SystemRequirements-Memory). That is too much, especially when you need it only for small personal projects. This project gives you minimal, but complete nexus functionality, with ability to save artifacts on filesystem or to s3 and basic auth. For most of usecases that is more than enought.

## How to use it ?

```
docker run -d -v /etc/nexus-minimal:/etc/nexus-minimal -p 8080:8080 astma/nexus-minimal
```

## Configuration 

Create config.yml in /etc/nexus-minimal or in the same directory where you run binary.

For s3:
```yml
---
http:
  addr: ":443"
  username: "myuser"
  password: "mypassword"
  https: true
  crt: "/certs/domain.crt"
  key: "/certs/domain.key"

storage:
  type: "s3"
  bucket_name: "my-super-nexus-bucket"
  access_key: "*******************"
  secret_key: "**************************************"
```

And for file system:
```yml
---
http:
  addr: ":8080"
  username: "myuser"
  password: "mypassword"

storage:
  type: "fs"
  base_dir: "/tmp/nexus-minimal"
```

## How to build it ?

Make sure to install golang, set all env variables etc.
Clone project to your go-workspace.
Cd to the project folder and run:

```
make build
```

And it will compile app.

Run:

```
make run
```

And it will run app locally on port 8080 by default.

## License

[MIT](LICENSE)

[license-url]: LICENSE

[license-image]: https://img.shields.io/github/license/mashape/apistatus.svg

[capture]: capture.png
