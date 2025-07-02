package xtrace

// The following keys for span tags and logs come from the OpenTracing Semantic Conventions,
// visit https://github.com/opentracing/specification/blob/master/semantic_conventions.md for
// more details.
// NOTE: Names that contain Palfish are customized tag names, which will not be respected by
// tools from the ecosystem.

// NOTE: Trace semantic conventions in OpenTelemetry
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/README.md
const (
	// OpenTracing
	TagComponent                     = "component"
	TagDBInstance                    = "db.instance"
	TagDBStatement                   = "db.statement"
	TagDBType                        = "db.type"
	TagDBUser                        = "db.user"
	TagPalfishDBCluster              = "db.cluster"
	TagPalfishDBTable                = "db.table"
	TagError                         = "error"
	TagHTTPMethod                    = "http.method"
	TagHTTPStatusCode                = "http.status_code"
	TagHTTPURL                       = "http.url"
	TagMessageBusDestination         = "message_bus.destination"
	TagPalfishMessageBusType         = "message_bus.type"
	TagPalfishKafkaConsumerBrokers   = "kafka.consumer.brokers"
	TagPalfishKafkaConsumerGroupID   = "kafka.consumer.group_id"
	TagPalfishKafkaConsumerPartition = "kafka.consumer.partition"
	TagPalfishCacheType              = "cache.type"
	TagPalfishCacheOp                = "cache.op"
	TagPalfishCacheKey               = "cache.key"
	TagPeerAddress                   = "peer.address"
	TagPeerHostname                  = "peer.hostname"
	TagPeerIPv4                      = "peer.ipv4"
	TagPeerIPv6                      = "peer.ipv6"
	TagPeerPort                      = "peer.port"
	TagPeerService                   = "peer.service"
	TagPeerSamplingPriority          = "sampling.priority"
	TagSpanKind                      = "span.kind"

	LogErrorKind    = "error.kind"
	LogErrorObject  = "error.object"
	LogErrorEvent   = "event"
	LogErrorMessage = "message"
	LogErrorStack   = "stack"

	// OpenTelemetry
	TagNetPeerIP         = "net.peer.ip"
	TagDBSystem          = "db.system"
	TagDBOperation       = "db.operation"
	TagDBSQLTable        = "db.sql.table"
	TagMongodbCollection = "db.mongodb.collection"

	TagMessagingSystem          = "messaging.system"
	TagMessagingDestination     = "messaging.destination"
	TagMessagingDestinationKind = "messaging.destination_kind"

	TagMessagingKafkaConsumerGroup = "messaging.kafka.consumer_group"
	TagMessagingKafkaPartition     = "messaging.kafka.partition"

	TagMessagingPulsarConsumerGroup = "messaging.pulsar.consumer_group"
)

const (
	SpanKindServer   = "server"
	SpanKindClient   = "client"
	SpanKindProducer = "producer"
	SpanKindConsumer = "consumer"

	DBTypeSQL     = "sql"
	DBTypeMongoDB = "mongodb"
	DBTypeRedis   = "redis"

	MessagingSystemKafka          = "kafka"
	MessagingSystemPulsar         = "pulsar"
	MessagingDestinationKindTopic = "topic"
	MessagingDestinationKindQueue = "queue"
)
