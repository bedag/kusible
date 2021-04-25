module github.com/bedag/kusible

go 1.15

require (
	github.com/Luzifer/go-openssl/v3 v3.1.0
	github.com/Shopify/ejson v1.2.2
	github.com/aws/aws-sdk-go v1.36.29
	github.com/gabriel-vasile/mimetype v1.1.2
	github.com/geofffranks/simpleyaml v0.0.0-20161109204137-c9320f076de5
	github.com/geofffranks/spruce v1.27.0
	github.com/go-test/deep v1.0.7
	github.com/gofrs/flock v0.8.0
	github.com/google/uuid v1.1.2
	github.com/imdario/mergo v0.3.11
	github.com/kjk/lzmadec v0.0.0-20200118223809-980b947af806
	github.com/kr/pretty v0.2.1 // indirect
	github.com/mitchellh/mapstructure v1.3.1
	github.com/olekukonko/tablewriter v0.0.2
	github.com/pborman/ansi v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	go.hein.dev/go-version v0.1.0
	gotest.tools v2.2.0+incompatible
	helm.sh/helm/v3 v3.5.3
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/cli-runtime v0.20.2
	k8s.io/client-go v0.20.2
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
)
