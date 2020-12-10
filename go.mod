module github.com/xelalexv/dregsy

go 1.13

require (
	github.com/aws/aws-sdk-go v1.32.5
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v20.10.0+incompatible
	github.com/goharbor/harbor/src v0.0.0-20201208094003-e92674a42ad7
	github.com/google/go-containerregistry v0.2.1
	github.com/heroku/docker-registry-client v0.0.0-20190909225348-afc9e1acc3d5
	github.com/moby/term v0.0.0-20201110203204-bea5bbe245bf // indirect
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58
	gopkg.in/yaml.v2 v2.3.0
)

// replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.7.0
