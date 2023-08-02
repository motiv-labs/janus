module github.com/hellofresh/janus

go 1.20

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20180108190415-b31f603f5e1e
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/prometheus v0.3.0
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/asaskevich/govalidator v0.0.0-20171111151018-521b25f4b05f
	github.com/cucumber/godog v0.10.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/felixge/httpsnoop v1.0.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-chi/chi v3.3.2+incompatible
	github.com/go-redis/redis/v7 v7.4.0
	github.com/gocql/gocql v0.0.0-20200624222514-34081eda590e
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/google/go-github v17.0.0+incompatible
	github.com/hellofresh/health-go/v3 v3.2.0
	github.com/hellofresh/logging-go v0.1.6
	github.com/hellofresh/opencensus-go-extras v0.0.0-20191004131501-7bd94f603dcf
	github.com/hellofresh/stats-go v0.8.0
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/magiconair/properties v1.8.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/rafaeljesus/retry-go v0.0.0-20171214204623-5981a380a879
	github.com/rs/cors v1.4.0
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tidwall/gjson v1.1.0
	github.com/ulule/limiter/v3 v3.5.0
	go.mongodb.org/mongo-driver v1.4.1
	go.opencensus.io v0.23.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)

require (
	github.com/DataDog/datadog-go v0.0.0-20180330214955-e67964b4021a // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/aws/aws-sdk-go v1.29.15 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bshuster-repo/logrus-logstash-hook v0.4.1 // indirect
	github.com/cactus/go-statsd-client v3.1.1+incompatible // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/cucumber/gherkin-go/v11 v11.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fiam/gounidecode v0.0.0-20150629112515-8deddbd03fec // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/go-querystring v0.0.0-20170111101155-53e6ce116135 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hashicorp/go-immutable-radix v1.2.0 // indirect
	github.com/hashicorp/go-memdb v1.2.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/klauspost/compress v1.10.10 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/onsi/gomega v1.10.1 // indirect
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.15.0 // indirect
	github.com/prometheus/procfs v0.2.0 // indirect
	github.com/prometheus/statsd_exporter v0.20.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a // indirect
	github.com/spf13/afero v1.1.2 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tidwall/match v1.0.0 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v0.0.0-20180714160509-73f8eece6fdc // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208 // indirect
	golang.org/x/sys v0.0.0-20201214210602-f9fddec55a1e // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/api v0.29.0 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6 // indirect
	gopkg.in/alexcesaro/statsd.v2 v2.0.0 // indirect
	gopkg.in/gemnasium/logrus-graylog-hook.v2 v2.0.6 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.51.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.12.0
