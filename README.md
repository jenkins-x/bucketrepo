[![MIT License][license-image]][license-url]

# Bucketrepo

This is a lightweight artifacts repository with low memory footprint which can be used as a minimal replacement for 
[Nexus](https://www.sonatype.com/nexus-repository-sonatype). It is able to cache artifacts from a remote repository
on a local filesystem volume and also to store them on a cloud storage bucket via [go-cloud](https://github.com/google/go-cloud/).

It can be deployed either as a side-car container to a Kubernetes build pod or as a standalone service.

## Getting Started

### Installation 

### Maven Configuration

## Acknowledgments

This project is originally based on [nexus-minimal](https://github.com/atsman/nexus-minimal). Thank you [assman](https://github.com/atsman) for creating that project.


## License

[MIT](LICENSE)

[license-url]: LICENSE

[license-image]: https://img.shields.io/github/license/mashape/apistatus.svg

[capture]: capture.png
