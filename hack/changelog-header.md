### Linux

```shell
curl -L https://github.com/jenkins-x/bucketrepo/releases/download/v{{.Version}}/bucketrepo-linux-amd64.tar.gz | tar xzv 
sudo mv bucketrepo /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x/bucketrepo/releases/download/v{{.Version}}/bucketrepo-darwin-amd64.tar.gz | tar xzv
sudo mv bucketrepo /usr/local/bin
```

