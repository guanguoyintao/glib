module git.umu.work/AI/uglib

go 1.16

replace (
	cloud.google.com/go/storage => cloud.google.com/go/storage v1.10.0
	github.com/micro/go-micro/v2 => git.umu.work/github/go-micro/v2 v2.9.8-0.20211221035607-ae5bb0075016
	google.golang.org/api => google.golang.org/api v0.56.0
	google.golang.org/grpc => google.golang.org/grpc v1.38.0
)

require (
	cloud.google.com/go/storage v1.10.0
	git.umu.work/be/goframework v1.1.22
	git.umu.work/eng/proto/tech/ai/asr_poly/go v1.2.0
	github.com/Shopify/sarama v1.29.1
	github.com/asticode/go-astisub v0.26.0
	github.com/aws/aws-sdk-go-v2 v1.13.0
	github.com/aws/aws-sdk-go-v2/credentials v1.8.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.24.1
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/bytedance/sonic v1.9.2
	github.com/emirpasic/gods v1.12.0
	github.com/go-playground/assert/v2 v2.0.1
	github.com/go-redis/redis/v8 v8.4.4
	github.com/go-resty/resty/v2 v2.7.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.2.0
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/liushuochen/gotable v0.0.0-20221119160816-1113793e7092
	github.com/micro/go-micro/v2 v2.9.1
	github.com/moznion/go-unicode-east-asian-width v0.0.0-20140622124307-0231aeb79f9b
	github.com/pkg/errors v0.9.1
	github.com/spaolacci/murmur3 v0.0.0-20180118202830-f09979ecbc72
	github.com/stretchr/testify v1.8.1
	github.com/tencentyun/cos-go-sdk-v5 v0.7.40
	golang.org/x/arch v0.4.0 // indirect
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e // indirect
	golang.org/x/net v0.0.0-20221014081412-f15817d10f9b
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	golang.org/x/sys v0.10.0 // indirect
	google.golang.org/api v0.54.0
	google.golang.org/grpc v1.40.0
	gopkg.in/vansante/go-ffprobe.v2 v2.1.1
	gorm.io/driver/mysql v1.4.0
	gorm.io/gen v0.3.18
	gorm.io/gorm v1.24.0
	gorm.io/plugin/dbresolver v1.3.0
)
