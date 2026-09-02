module github.com/Azure/karpenter-provider-azure

go 1.26.4

require (
	dario.cat/mergo v1.0.2
	github.com/Azure/aks-middleware v0.0.42
	github.com/Azure/azure-kusto-go v0.16.1
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible
	github.com/Azure/azure-sdk-for-go-extensions v0.6.0
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.22.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.14.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2 v2.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7 v7.3.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v9 v9.4.0-beta.1
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault v1.5.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi v1.3.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork v1.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph v0.10.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions v1.3.0
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys v1.5.0
	github.com/Azure/go-autorest/autorest v0.11.30
	github.com/Azure/skewer v0.0.24
	github.com/Pallinder/go-randomdata v1.2.0
	github.com/alecthomas/kong v1.15.0
	github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b
	github.com/awslabs/operatorpkg v0.0.0-20260501204335-c49b4ca8d58d
	github.com/blang/semver/v4 v4.0.0
	github.com/go-logr/logr v1.4.3
	github.com/go-logr/zapr v1.3.0
	github.com/go-openapi/errors v0.22.8
	github.com/go-openapi/runtime v0.32.4
	github.com/go-openapi/strfmt v0.26.3
	github.com/go-openapi/swag v0.26.1
	github.com/go-openapi/validate v0.26.0
	github.com/go-playground/validator/v10 v10.30.3
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/golangci/golangci-lint/v2 v2.13.2
	github.com/google/go-cmp v0.7.0
	github.com/google/uuid v1.6.0
	github.com/jongio/azidext/go/azidext v0.5.0
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/onsi/ginkgo/v2 v2.32.0
	github.com/onsi/gomega v1.42.1
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v1.23.2
	github.com/prometheus/client_model v0.6.2
	github.com/samber/lo v1.53.0
	github.com/stretchr/testify v1.12.1
	go.uber.org/multierr v1.11.0
	go.yaml.in/yaml/v2 v2.4.4
	golang.org/x/sync v0.22.0
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.35.1
	k8s.io/apiextensions-apiserver v0.35.0
	k8s.io/apimachinery v0.35.1
	k8s.io/client-go v0.35.1
	k8s.io/component-base v0.35.0
	k8s.io/utils v0.0.0-20251222233032-718f0e51e6d2
	sigs.k8s.io/cloud-provider-azure/pkg/azclient v0.14.3
	sigs.k8s.io/controller-runtime v0.23.1
	sigs.k8s.io/karpenter v1.13.0
)

require (
	4d63.com/gocheckcompilerdirectives v1.4.0 // indirect
	4d63.com/gochecknoglobals v0.2.2 // indirect
	charm.land/lipgloss/v2 v2.0.6 // indirect
	codeberg.org/chavacava/garif v0.2.0 // indirect
	codeberg.org/polyfloyd/go-errorlint v1.9.0 // indirect
	dev.gaijin.team/go/exhaustruct/v4 v4.0.0 // indirect
	dev.gaijin.team/go/exhaustruct/v5 v5.0.3 // indirect
	dev.gaijin.team/go/golib v0.8.1 // indirect
	github.com/4meepo/tagalign v1.4.3 // indirect
	github.com/Abirdcfly/dupword v0.1.8 // indirect
	github.com/AdminBenni/iota-mixing v1.0.0 // indirect
	github.com/AlwxSin/noinlineerr v1.0.6 // indirect
	github.com/Antonboom/errname v1.1.2 // indirect
	github.com/Antonboom/nilnil v1.1.2 // indirect
	github.com/Antonboom/testifylint v1.6.4 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.12.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v6 v6.6.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6 v6.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns v1.3.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage/v2 v2.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets v1.4.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal v1.2.0 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.24 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.1 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/Azure/msi-dataplane v0.4.3 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.7.2 // indirect
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/ClickHouse/clickhouse-go-linter v1.2.1 // indirect
	github.com/Djarvur/go-err113 v0.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/MirrexOne/unqueryvet v1.5.4 // indirect
	github.com/OpenPeeDeeP/depguard/v2 v2.2.1 // indirect
	github.com/alecthomas/chroma/v2 v2.27.0 // indirect
	github.com/alecthomas/go-check-sumtype v0.3.1 // indirect
	github.com/alexkohler/nakedret/v2 v2.0.6 // indirect
	github.com/alexkohler/prealloc v1.1.0 // indirect
	github.com/alfatraining/structtag v1.0.0 // indirect
	github.com/alingse/asasalint v0.0.11 // indirect
	github.com/alingse/nilnesserr v0.2.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/ashanbrown/forbidigo/v2 v2.3.1 // indirect
	github.com/ashanbrown/makezero/v2 v2.2.1 // indirect
	github.com/avast/retry-go v3.0.0+incompatible // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bkielbasa/cyclop v1.2.3 // indirect
	github.com/blizzy78/varnamelen v0.8.0 // indirect
	github.com/bombsimon/wsl/v4 v4.7.0 // indirect
	github.com/bombsimon/wsl/v5 v5.9.0 // indirect
	github.com/breml/bidichk v0.3.3 // indirect
	github.com/breml/errchkjson v0.4.1 // indirect
	github.com/butuzov/ireturn v0.4.1 // indirect
	github.com/butuzov/mirror v1.3.3 // indirect
	github.com/catenacyber/perfsprint v0.10.1 // indirect
	github.com/ccojocar/zxcvbn-go v1.0.4 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/charithe/durationcheck v0.0.11 // indirect
	github.com/charmbracelet/colorprofile v0.4.3 // indirect
	github.com/charmbracelet/ultraviolet v0.0.0-20260811164956-006e29f97886 // indirect
	github.com/charmbracelet/x/ansi v0.11.8 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/ckaznocha/intrange v0.3.1 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/curioswitch/go-reassign v0.3.0 // indirect
	github.com/daixiang0/gci v0.13.7 // indirect
	github.com/dave/dst v0.27.3 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/denis-tingaikin/go-header v0.5.0 // indirect
	github.com/dlclark/regexp2/v2 v2.2.1 // indirect
	github.com/emicklei/go-restful/v3 v3.13.0 // indirect
	github.com/ettle/strcase v0.2.0 // indirect
	github.com/evanphx/json-patch v5.9.11+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/felixge/httpsnoop v1.1.0 // indirect
	github.com/firefart/nonamedreturns v1.0.8 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/fzipp/gocyclo v0.6.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/ghostiam/protogetter v0.3.21 // indirect
	github.com/go-critic/go-critic v0.14.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/analysis v0.25.2 // indirect
	github.com/go-openapi/jsonpointer v0.23.1 // indirect
	github.com/go-openapi/jsonreference v0.21.6 // indirect
	github.com/go-openapi/loads v0.24.0 // indirect
	github.com/go-openapi/runtime/server-middleware v0.30.0 // indirect
	github.com/go-openapi/spec v0.22.6 // indirect
	github.com/go-openapi/swag/cmdutils v0.26.1 // indirect
	github.com/go-openapi/swag/conv v0.26.1 // indirect
	github.com/go-openapi/swag/fileutils v0.26.1 // indirect
	github.com/go-openapi/swag/jsonname v0.26.1 // indirect
	github.com/go-openapi/swag/jsonutils v0.26.1 // indirect
	github.com/go-openapi/swag/loading v0.26.1 // indirect
	github.com/go-openapi/swag/mangling v0.26.1 // indirect
	github.com/go-openapi/swag/netutils v0.26.1 // indirect
	github.com/go-openapi/swag/stringutils v0.26.1 // indirect
	github.com/go-openapi/swag/typeutils v0.26.1 // indirect
	github.com/go-openapi/swag/yamlutils v0.26.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/go-toolsmith/astcast v1.1.0 // indirect
	github.com/go-toolsmith/astcopy v1.1.0 // indirect
	github.com/go-toolsmith/astequal v1.2.0 // indirect
	github.com/go-toolsmith/astfmt v1.1.0 // indirect
	github.com/go-toolsmith/astp v1.1.0 // indirect
	github.com/go-toolsmith/strparse v1.1.0 // indirect
	github.com/go-toolsmith/typep v1.1.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/go-xmlfmt/xmlfmt v1.1.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/godoc-lint/godoc-lint v0.11.2 // indirect
	github.com/gofrs/flock v0.13.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/golangci/asciicheck v0.5.0 // indirect
	github.com/golangci/canonicalheader v0.0.0-20260827115959-a25c71c521f6 // indirect
	github.com/golangci/dupl v0.0.0-20260401084720-c99c5cf5c202 // indirect
	github.com/golangci/go-printf-func-name v0.1.1 // indirect
	github.com/golangci/gofmt v0.0.0-20260820135601-e84e05053792 // indirect
	github.com/golangci/golines v0.15.0 // indirect
	github.com/golangci/misspell v0.8.0 // indirect
	github.com/golangci/plugin-module-register v0.1.2 // indirect
	github.com/golangci/revgrep v0.8.0 // indirect
	github.com/golangci/rowserrcheck v0.0.0-20260419091836-c5f79b8a11ba // indirect
	github.com/golangci/swaggoswag v0.0.0-20250504205917-77f2aca3143e // indirect
	github.com/golangci/unconvert v0.0.0-20250410112200-a129a6e6413e // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/gnostic-models v0.7.1 // indirect
	github.com/google/pprof v0.0.0-20260709232956-b9395ee17fa0 // indirect
	github.com/gordonklaus/ineffassign v0.2.0 // indirect
	github.com/gostaticanalysis/analysisutil v0.7.1 // indirect
	github.com/gostaticanalysis/comment v1.5.0 // indirect
	github.com/gostaticanalysis/forcetypeassert v0.2.0 // indirect
	github.com/gostaticanalysis/nilerr v0.1.2 // indirect
	github.com/hashicorp/go-immutable-radix/v2 v2.1.0 // indirect
	github.com/hashicorp/go-version v1.9.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hexops/gotextdiff v1.0.3 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jgautheron/goconst v1.11.0 // indirect
	github.com/jjti/go-spancheck v0.6.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/julz/importas v0.2.0 // indirect
	github.com/karamaru-alpha/copyloopvar v1.2.2 // indirect
	github.com/kisielk/errcheck v1.20.0 // indirect
	github.com/kkHAIKE/contextcheck v1.1.6 // indirect
	github.com/kulti/thelper v0.7.1 // indirect
	github.com/kunwardeep/paralleltest v1.0.15 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/ldez/exptostd v0.4.5 // indirect
	github.com/ldez/gomoddirectives v0.9.0 // indirect
	github.com/ldez/grignotin v0.10.1 // indirect
	github.com/ldez/structtags v0.6.1 // indirect
	github.com/ldez/tagliatelle v0.7.2 // indirect
	github.com/ldez/usetesting v0.5.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/leonklingele/grouper v1.1.2 // indirect
	github.com/lucasb-eyer/go-colorful v1.4.1 // indirect
	github.com/macabu/inamedparam v0.2.0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/manuelarte/embeddedstructfieldcheck v0.4.0 // indirect
	github.com/manuelarte/funcorder v0.6.0 // indirect
	github.com/maratori/testableexamples v1.0.1 // indirect
	github.com/maratori/testpackage v1.1.2 // indirect
	github.com/matoous/godox v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.24 // indirect
	github.com/mgechev/revive v1.15.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/moricho/tparallel v0.3.2 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nakabonne/nestif v0.3.1 // indirect
	github.com/nishanths/exhaustive v0.12.0 // indirect
	github.com/nishanths/predeclared v0.2.2 // indirect
	github.com/nunnatsa/ginkgolinter v0.24.0 // indirect
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.4.3 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/common v0.67.4 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/quasilyte/go-ruleguard v0.4.5 // indirect
	github.com/quasilyte/go-ruleguard/dsl v0.3.23 // indirect
	github.com/quasilyte/gogrep v0.5.0 // indirect
	github.com/quasilyte/regex/syntax v0.0.0-20210819130434-b3f0c404a727 // indirect
	github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567 // indirect
	github.com/raeperd/recvcheck v0.3.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rogpeppe/go-internal v1.16.0 // indirect
	github.com/ryancurrah/gomodguard v1.4.1 // indirect
	github.com/ryancurrah/gomodguard/v2 v2.1.3 // indirect
	github.com/ryanrolds/sqlclosecheck v0.6.0 // indirect
	github.com/sanposhiho/wastedassign/v2 v2.1.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.3 // indirect
	github.com/sashamelentyev/interfacebloat v1.1.0 // indirect
	github.com/sashamelentyev/usestdlibvars v1.29.0 // indirect
	github.com/securego/gosec/v2 v2.28.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/sirupsen/logrus v1.10.1 // indirect
	github.com/sivchari/containedctx v1.0.3 // indirect
	github.com/sonatard/noctx v0.5.1 // indirect
	github.com/sourcegraph/go-diff v0.8.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spf13/viper v1.12.0 // indirect
	github.com/ssgreg/nlreturn/v2 v2.2.1 // indirect
	github.com/stbenjam/no-sprintf-host-port v0.3.1 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	github.com/tetafro/godot v1.5.6 // indirect
	github.com/timakin/bodyclose v0.0.0-20260129054331-73d1f95b84b4 // indirect
	github.com/timonwong/loggercheck v0.11.0 // indirect
	github.com/tomarrell/wrapcheck/v2 v2.12.0 // indirect
	github.com/tommy-muehle/go-mnd/v2 v2.5.1 // indirect
	github.com/ultraware/funlen v0.2.0 // indirect
	github.com/ultraware/whitespace v0.2.0 // indirect
	github.com/uudashr/gocognit v1.2.1 // indirect
	github.com/uudashr/iface v1.5.1 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xen0n/gosmopolitan v1.3.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	github.com/yagipy/maintidx v1.0.0 // indirect
	github.com/yeya24/promlinter v0.3.0 // indirect
	github.com/ykadowak/zerologlint v0.1.5 // indirect
	gitlab.com/bosi/decorder v0.4.2 // indirect
	go-simpler.org/musttag v0.14.0 // indirect
	go-simpler.org/sloglint v0.12.0 // indirect
	go.augendre.info/arangolint v0.4.0 // indirect
	go.augendre.info/fatcontext v0.10.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.69.0 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/metric v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
	go.uber.org/mock v0.6.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b // indirect
	golang.org/x/exp/typeparams v0.0.0-20260811152304-ee035b5b010f // indirect
	golang.org/x/mod v0.40.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/term v0.45.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	golang.org/x/tools v0.49.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.5.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/grpc v1.82.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/dnaeon/go-vcr.v3 v3.2.0 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.13.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	honnef.co/go/tools v0.8.1 // indirect
	k8s.io/cloud-provider v0.35.0 // indirect
	k8s.io/component-helpers v0.35.0 // indirect
	k8s.io/csi-translation-lib v0.35.0 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20251125145642-4e65d59e963e // indirect
	mvdan.cc/gofumpt v0.11.0 // indirect
	mvdan.cc/unparam v0.0.0-20260823230713-2fa3d841b0c8 // indirect
	sigs.k8s.io/json v0.0.0-20250730193827-2d320260d730 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v6 v6.3.2-0.20260122202528-d9cc6641c482 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace sigs.k8s.io/karpenter => github.com/openshift/kubernetes-sigs-karpenter v0.0.0-20260721214330-2b8ed744bf1c
