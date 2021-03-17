module github.com/hellofresh/janus

go 1.15

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20180108190415-b31f603f5e1e
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/prometheus v0.3.0
	github.com/DataDog/datadog-go v0.0.0-20180330214955-e67964b4021a // indirect
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/asaskevich/govalidator v0.0.0-20171111151018-521b25f4b05f
	github.com/bshuster-repo/logrus-logstash-hook v0.4.1 // indirect
	github.com/cactus/go-statsd-client v3.1.1+incompatible // indirect
	github.com/cucumber/godog v0.10.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/felixge/httpsnoop v1.0.0
	github.com/fiam/gounidecode v0.0.0-20150629112515-8deddbd03fec // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-chi/chi v3.3.2+incompatible
	github.com/go-redis/redis/v7 v7.4.0
	github.com/gocql/gocql v0.0.0-20200624222514-34081eda590e
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v0.0.0-20170111101155-53e6ce116135 // indirect
	github.com/hellofresh/health-go/v3 v3.2.0
	github.com/hellofresh/logging-go v0.1.6
	github.com/hellofresh/opencensus-go-extras v0.0.0-20191004131501-7bd94f603dcf
	github.com/hellofresh/stats-go v0.8.0
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/klauspost/compress v1.10.10 // indirect
	github.com/magiconair/properties v1.8.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/rafaeljesus/retry-go v0.0.0-20171214204623-5981a380a879
	github.com/rs/cors v1.4.0
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tidwall/gjson v1.1.0
	github.com/tidwall/match v1.0.0 // indirect
	github.com/ulule/limiter/v3 v3.5.0
	go.mongodb.org/mongo-driver v1.4.1
	go.opencensus.io v0.23.0
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/alexcesaro/statsd.v2 v2.0.0 // indirect
	gopkg.in/gemnasium/logrus-graylog-hook.v2 v2.0.6 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.12.0
