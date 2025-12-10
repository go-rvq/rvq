module github.com/go-rvq/rvq

go 1.24

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/a8m/envsubst v1.4.3
	github.com/ahmetb/go-linq/v3 v3.2.0
	github.com/aws/aws-sdk-go v1.42.34
	github.com/aws/aws-sdk-go-v2 v1.41.0
	github.com/aws/aws-sdk-go-v2/config v1.32.4
	github.com/aws/aws-sdk-go-v2/credentials v1.19.4
	github.com/aws/aws-sdk-go-v2/service/s3 v1.93.1
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.4
	github.com/disintegration/imaging v1.6.2
	github.com/dustin/go-humanize v1.0.1
	github.com/fatih/color v1.17.0
	github.com/gad-lang/gad v0.0.0-20251118140158-818a37fe57e6
	github.com/go-chi/chi/v5 v5.0.12
	github.com/go-playground/form v3.1.4+incompatible
	github.com/go-playground/form/v4 v4.2.1
	github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-rvq/htmlgo v1.0.2
	github.com/gobwas/glob v0.2.3
	github.com/gocarina/gocsv v0.0.0-20240520201108-78e41c74b4b1
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/google/go-cmp v0.7.0
	github.com/google/uuid v1.6.0
	github.com/gosimple/slug v1.14.0
	github.com/gosimple/unidecode v1.0.1
	github.com/hack-pad/hackpadfs v0.2.4
	github.com/hashicorp/go-multierror v1.1.1
	github.com/iancoleman/strcase v0.3.0
	github.com/jinzhu/configor v1.2.1
	github.com/jinzhu/inflection v1.0.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/manifoldco/promptui v0.9.0
	github.com/markbates/goth v1.80.0
	github.com/mattn/go-tty v0.0.7
	github.com/mholt/archiver/v4 v4.0.0-alpha.8
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/moisespsena-go/tree2html v0.0.0-20240923160936-f6c5cb6a48f1
	github.com/mpvl/unique v0.0.0-20150818121801-cbe035fff7de
	github.com/ory/ladon v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/pquerna/otp v1.4.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/rs/xid v1.5.0
	github.com/samber/lo v1.39.0
	github.com/shopspring/decimal v1.4.0
	github.com/shurcooL/sanitized_anchor_name v1.0.0
	github.com/spf13/cast v1.6.0
	github.com/stretchr/testify v1.10.0
	github.com/sunfmin/reflectutils v1.0.4
	github.com/sunfmin/snippetgo v0.0.2
	github.com/tailscale/hujson v0.0.0-20250605163823-992244df8c5a
	github.com/theplant/bimg v1.1.1
	github.com/theplant/docgo v0.0.16
	github.com/theplant/gofixtures v1.1.2
	github.com/theplant/htmltestingutils v0.0.0-20190423050759-0e06de7b6967
	github.com/theplant/osenv v0.0.1
	github.com/theplant/sliceutils v0.0.0-20200406042209-89153d988eb1
	github.com/theplant/testenv v0.0.0-20240513012518-1c94c8c84239
	github.com/theplant/testingutils v0.0.2
	github.com/tnclong/go-que v0.0.0-20240226030728-4e1f3c8ec781
	github.com/ua-parser/uap-go v0.0.0-20240113215029-33f8e6d47f38
	github.com/wcharczuk/go-chart/v2 v2.1.1
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
	goji.io/v3 v3.0.0
	golang.org/x/crypto v0.38.0
	golang.org/x/text v0.25.0
	gorm.io/driver/mysql v1.5.7
	gorm.io/driver/postgres v1.5.7
	gorm.io/driver/sqlite v1.5.5
	gorm.io/gorm v1.25.10
	gorm.io/hints v1.1.2
)

require (
	cloud.google.com/go/compute/metadata v0.7.0 // indirect
	dario.cat/mergo v1.0.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Microsoft/hcsshim v0.12.3 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.4 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.12 // indirect
	github.com/aws/smithy-go v1.24.0 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/blend/go-sdk v1.20220411.3 // indirect
	github.com/bodgit/plumbing v1.3.0 // indirect
	github.com/bodgit/sevenzip v1.5.1 // indirect
	github.com/bodgit/windows v1.0.1 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/containerd/containerd v1.7.16 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/cpuguy83/dockercfg v0.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/dlclark/regexp2 v1.11.0 // indirect
	github.com/docker/docker v26.1.2+incompatible // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.2.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/klauspost/pgzip v1.2.6 // indirect
	github.com/lufia/plan9stats v0.0.0-20240408141607-282e7b5d6b74 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/markbates/going v1.0.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/microcosm-cc/bluemonday v1.0.26 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/patternmatcher v0.6.0 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/sys/user v0.1.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mrjones/oauth v0.0.0-20190623134757-126b35219450 // indirect
	github.com/nwaples/rardecode/v2 v2.0.0-beta.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/ory/pagination v0.0.1 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/shirou/gopsutil/v3 v3.24.4 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/shurcooL/github_flavored_markdown v0.0.0-20210228213109-c3a9aa474629 // indirect
	github.com/shurcooL/highlight_diff v0.0.0-20230708024848-22f825814995 // indirect
	github.com/shurcooL/highlight_go v0.0.0-20230708025100-33e05792540a // indirect
	github.com/shurcooL/octicon v0.0.0-20230705024016-66bff059edb8 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/sourcegraph/annotate v0.0.0-20160123013949-f4cad6c6324d // indirect
	github.com/sourcegraph/syntaxhighlight v0.0.0-20170531221838-bd320f5d308e // indirect
	github.com/testcontainers/testcontainers-go v0.31.0 // indirect
	github.com/therootcompany/xz v1.0.1 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	github.com/yosssi/gohtml v0.0.0-20201013000340-ee4748c638f4
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.60.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go4.org v0.0.0-20230225012048-214862532bf5 // indirect
	golang.org/x/exp v0.0.0-20240604190554-fc45aab8b7f8 // indirect
	golang.org/x/image v0.17.0 // indirect
	golang.org/x/net v0.40.0
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250512202823-5a2f75b736a9 // indirect
	google.golang.org/grpc v1.72.1 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.4 // indirect
	github.com/qor5/web v1.2.3 // indirect
	github.com/theplant/htmlgo v1.0.3 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250505200425-f936aa4a68b2 // indirect
)

//replace github.com/theplant/docgo => ../../docgo/
