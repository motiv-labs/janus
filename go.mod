module github.com/hellofresh/janus

go 1.15

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20180108190415-b31f603f5e1e
	github.com/DataDog/datadog-go v0.0.0-20180330214955-e67964b4021a // indirect
	github.com/Knetic/govaluate v3.0.0+incompatible
	github.com/afex/hystrix-go v0.0.0-20180406012432-f86abeeb9f72
	github.com/asaskevich/govalidator v0.0.0-20171111151018-521b25f4b05f
	github.com/bshuster-repo/logrus-logstash-hook v0.4.1 // indirect
	github.com/cactus/go-statsd-client v3.1.1+incompatible // indirect
	github.com/cucumber/godog v0.10.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/felixge/httpsnoop v1.0.0
	github.com/fiam/gounidecode v0.0.0-20150629112515-8deddbd03fec // indirect
	github.com/fsnotify/fsnotify v1.4.7
	github.com/globalsign/mgo v0.0.0-20180615134936-113d3961e731
	github.com/go-chi/chi v3.3.2+incompatible
	github.com/go-redis/redis v6.12.0+incompatible
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v0.0.0-20170111101155-53e6ce116135 // indirect
	github.com/hashicorp/hcl v0.0.0-20171017181929-23c074d0eceb // indirect
	github.com/hellofresh/health-go/v3 v3.2.0
	github.com/hellofresh/logging-go v0.1.6
	github.com/hellofresh/opencensus-go-extras v0.0.0-20191004131501-7bd94f603dcf
	github.com/hellofresh/stats-go v0.8.0
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/magiconair/properties v1.7.4
	github.com/mitchellh/go-homedir v0.0.0-20161203194507-b8bc1bf76747
	github.com/mitchellh/mapstructure v0.0.0-20180111000720-b4575eea38cc
	github.com/prometheus/client_golang v0.9.1 // indirect
	github.com/prometheus/common v0.0.0-20181109100915-0b1957f9d949 // indirect
	github.com/prometheus/procfs v0.0.0-20181005140218-185b4288413d // indirect
	github.com/rafaeljesus/retry-go v0.0.0-20171214204623-5981a380a879
	github.com/rcrowley/go-metrics v0.0.0-20180406234716-d932a24a8ccb // indirect
	github.com/rs/cors v1.4.0
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337 // indirect
	github.com/spf13/afero v1.0.2 // indirect
	github.com/spf13/cast v1.1.0 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/jwalterweatherman v0.0.0-20180109140146-7c0cea34c8ec // indirect
	github.com/spf13/viper v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/tidwall/gjson v1.1.0
	github.com/tidwall/match v1.0.0 // indirect
	github.com/ulule/limiter v2.2.2+incompatible
	go.opencensus.io v0.18.0
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be
	gopkg.in/alexcesaro/statsd.v2 v2.0.0 // indirect
	gopkg.in/gemnasium/logrus-graylog-hook.v2 v2.0.6 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.12.0
