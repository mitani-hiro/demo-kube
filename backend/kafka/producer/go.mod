module producer

go 1.24.2

require (
	common v0.0.0-00010101000000-000000000000
	github.com/segmentio/kafka-go v0.4.47
	google.golang.org/grpc v1.72.1
	proto v0.0.0-00010101000000-000000000000
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250519155744-55703ea1f237 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace (
	common => ../../common
	proto => ../../proto
)
