package models

import (
	"time"
)

type Monitor struct {
	ID                         uint           `gorm:"primaryKey" json:"id"`
	Name                       *string        `json:"name"`
	Active                     bool           `gorm:"default:true" json:"active"`
	UserID                     *uint          `json:"user_id"`
	Interval                   int            `gorm:"default:20" json:"interval"`
	URL                        *string        `json:"url"`
	Type                       *string        `json:"type"`
	Weight                     *int           `json:"weight"`
	Hostname                   *string        `json:"hostname"`
	Port                       *int           `json:"port"`
	CreatedDate                time.Time      `json:"created_date"`
	Keyword                    *string        `json:"keyword"`
	MaxRetries                 int            `gorm:"default:0" json:"maxretries"`
	IgnoreTLS                  bool           `gorm:"default:false" json:"ignore_tls"`
	UpsideDown                 bool           `gorm:"default:false" json:"upside_down"`
	MaxRedirects               int            `gorm:"default:10" json:"maxredirects"`
	AcceptedStatusCodesJSON    string         `gorm:"default:'[\"200-299\"]'" json:"accepted_statuscodes_json"`
	DNSResolveType             *string        `json:"dns_resolve_type"`
	DNSResolveServer           *string        `json:"dns_resolve_server"`
	DNSLastResult              *string        `json:"dns_last_result"`
	RetryInterval              int            `gorm:"default:0" json:"retry_interval"`
	PushToken                  *string        `json:"pushToken"`
	Method                     string         `gorm:"default:'GET'" json:"method"`
	Body                       *string        `json:"body"`
	Headers                    *string        `json:"headers"`
	BasicAuthUser              *string        `json:"basic_auth_user"`
	BasicAuthPass              *string        `json:"basic_auth_pass"`
	DockerHost                 *uint          `json:"docker_host"`
	DockerContainer            *string        `json:"docker_container"`
	ProxyID                    *uint          `json:"proxy_id"`
	ExpiryNotification         *bool          `json:"expiry_notification"`
	MQTTTopic                  *string        `json:"mqtt_topic"`
	MQTTSuccessMessage         *string        `json:"mqtt_success_message"`
	MQTTUsername               *string        `json:"mqtt_username"`
	MQTTPassword               *string        `json:"mqtt_password"`
	DatabaseConnectionString   *string        `json:"database_connection_string"`
	DatabaseQuery              *string        `json:"database_query"`
	AuthMethod                 *string        `json:"auth_method"`
	AuthDomain                 *string        `json:"auth_domain"`
	AuthWorkstation            *string        `json:"auth_workstation"`
	GRPCURL                    *string        `json:"grpc_url"`
	GRPCProtobuf               *string        `json:"grpc_protobuf"`
	GRPCBody                   *string        `json:"grpc_body"`
	GRPCMetadata               *string        `json:"grpc_metadata"`
	GRPCMethod                 *string        `json:"grpc_method"`
	GRPCServiceName            *string        `json:"grpc_service_name"`
	GRPCEnableTLS              bool           `gorm:"default:false" json:"grpc_enable_tls"`
	RadiusUsername             *string        `json:"radius_username"`
	RadiusPassword             *string        `json:"radius_password"`
	RadiusCallingStationID     *string        `json:"radius_calling_station_id"`
	RadiusCalledStationID      *string        `json:"radius_called_station_id"`
	RadiusSecret               *string        `json:"radius_secret"`
	ResendInterval             int            `gorm:"default:0" json:"resend_interval"`
	PacketSize                 int            `gorm:"default:56" json:"packet_size"`
	Game                       *string        `json:"game"`
	HTTPBodyEncoding           *string        `json:"http_body_encoding"`
	Description                *string        `json:"description"`
	TLSCA                      *string        `json:"tls_ca"`
	TLSCert                    *string        `json:"tls_cert"`
	TLSKey                     *string        `json:"tls_key"`
	Parent                     *uint          `json:"parent"`
	InvertKeyword              bool           `gorm:"default:false" json:"invert_keyword"`
	JSONPath                   *string        `json:"json_path"`
	ExpectedValue              *string        `json:"expected_value"`
	KafkaProducerTopic         *string        `json:"kafka_producer_topic"`
	KafkaProducerBrokers       *string        `json:"kafka_producer_brokers"`
	KafkaProducerSSL           bool           `gorm:"default:false" json:"kafka_producer_ssl"`
	KafkaProducerAllowAutoTopicCreation bool   `gorm:"default:false" json:"kafka_producer_allow_auto_topic_creation"`
	KafkaProducerSASLOptions   *string        `json:"kafka_producer_sasl_options"`
	KafkaProducerMessage       *string        `json:"kafka_producer_message"`
	OAuthClientID              *string        `json:"oauth_client_id"`
	OAuthClientSecret          *string        `json:"oauth_client_secret"`
	OAuthTokenURL              *string        `json:"oauth_token_url"`
	OAuthScopes                *string        `json:"oauth_scopes"`
	OAuthAuthMethod            *string        `json:"oauth_auth_method"`
	Timeout                    float64        `gorm:"default:0" json:"timeout"`
	GameDigGivenPortOnly       bool           `gorm:"default:true" json:"gamedig_given_port_only"`
	MQTTCheckType              string         `gorm:"default:'keyword'" json:"mqtt_check_type"`
	RemoteBrowser              *uint          `json:"remote_browser"`
	SNMPOID                    *string        `json:"snmp_oid"`
	SNMPVersion                *string        `gorm:"default:'2c'" json:"snmp_version"`
	JSONPathOperator           *string        `json:"json_path_operator"`
	CacheBust                  bool           `gorm:"default:false" json:"cache_bust"`
	Conditions                 string         `gorm:"default:'[]'" json:"conditions"`
	RabbitMQNodes              *string        `json:"rabbitmq_nodes"`
	RabbitMQUsername           *string        `json:"rabbitmq_username"`
	RabbitMQPassword           *string        `json:"rabbitmq_password"`
	SMTPSecurity               *string        `json:"smtp_security"`
	WSIgnoreSecWebsocketAcceptHeader bool     `gorm:"default:false" json:"ws_ignore_sec_websocket_accept_header"`
	WSSubprotocol              string         `gorm:"default:''" json:"ws_subprotocol"`
	PingCount                  int            `gorm:"default:1" json:"ping_count"`
	PingNumeric                bool           `gorm:"default:true" json:"ping_numeric"`
	PingPerRequestTimeout      int            `gorm:"default:2" json:"ping_per_request_timeout"`
	IPFamily                   *string        `json:"ip_family"`
	ManualStatus               *int           `json:"manual_status"`
	OAuthAudience              *string        `json:"oauth_audience"`
	MQTTWebsocketPath          *string        `json:"mqtt_websocket_path"`
	DomainExpiryNotification   *bool          `json:"domain_expiry_notification"`
	SaveResponse               bool           `gorm:"default:false" json:"save_response"`
	SaveErrorResponse          bool           `gorm:"default:true" json:"save_error_response"`
	ResponseMaxLength          int            `gorm:"default:1024" json:"response_max_length"`
	SystemServiceName          *string        `json:"system_service_name"`
	Subtype                    *string        `json:"subtype"`
	Location                   *string        `json:"location"`
	Protocol                   *string        `json:"protocol"`
	SNMPV3Username             *string        `json:"snmp_v3_username"`
	ExpectedTLSAlert           *string        `json:"expected_tls_alert"`
	RetryOnlyOnStatusCodeFailure bool         `gorm:"default:false" json:"retry_only_on_status_code_failure"`
	ScreenshotDelay            uint           `gorm:"default:0" json:"screenshot_delay"`
	NTPStratumThreshold        *int           `json:"ntp_stratum_threshold"`
	NTPTimeOffsetThreshold     *int           `json:"ntp_time_offset_threshold"`
	NTPRootDispersionThreshold *int           `json:"ntp_root_dispersion_threshold"`
	BearerToken                *string        `json:"bearer_token"`
	GameDigToken               *string        `json:"gamedig_token"`
	SSHUsername                *string        `json:"ssh_username"`
	SSHPassword                *string        `json:"ssh_password"`
	SFTPPath                   *string        `json:"sftp_path"`
	SSHPrivateKey              *string        `json:"ssh_private_key"`
	SSHPassphrase              *string        `json:"ssh_passphrase"`
	SSHAuthMethod              string         `gorm:"default:'password'" json:"ssh_auth_method"`
	Tags                       []MonitorTag `gorm:"foreignKey:MonitorID" json:"tags,omitempty"`
}

func (Monitor) TableName() string {
	return "monitor"
}

type MonitorCreateRequest struct {
	Name                       *string `json:"name" validate:"required"`
	Active                     *bool   `json:"active"`
	Interval                   *int    `json:"interval" validate:"required,min=20"`
	URL                        *string `json:"url"`
	Type                       *string `json:"type" validate:"required"`
	Weight                     *int    `json:"weight"`
	Hostname                   *string `json:"hostname"`
	Port                       *int    `json:"port" validate:"omitempty,min=0,max=65535"`
	Keyword                    *string `json:"keyword"`
	MaxRetries                 *int    `json:"maxretries" validate:"min=0"`
	IgnoreTLS                  *bool   `json:"ignore_tls"`
	UpsideDown                 *bool   `json:"upside_down"`
	MaxRedirects               *int    `json:"maxredirects" validate:"min=0"`
	AcceptedStatusCodes        []string `json:"accepted_statuscodes"`
	DNSResolveType             *string `json:"dns_resolve_type"`
	DNSResolveServer           *string `json:"dns_resolve_server"`
	RetryInterval              *int    `json:"retry_interval" validate:"min=20"`
	PushToken                  *string `json:"pushToken"`
	Method                     *string `json:"method"`
	Body                       *string `json:"body"`
	Headers                    *string `json:"headers"`
	BasicAuthUser              *string `json:"basic_auth_user"`
	BasicAuthPass              *string `json:"basic_auth_pass"`
	DockerHost                 *uint   `json:"docker_host"`
	DockerContainer            *string `json:"docker_container"`
	ProxyID                    *uint   `json:"proxy_id"`
	ExpiryNotification         *bool   `json:"expiry_notification"`
	MQTTTopic                  *string `json:"mqtt_topic"`
	MQTTSuccessMessage         *string `json:"mqtt_success_message"`
	MQTTUsername               *string `json:"mqtt_username"`
	MQTTPassword               *string `json:"mqtt_password"`
	DatabaseConnectionString   *string `json:"database_connection_string"`
	DatabaseQuery              *string `json:"database_query"`
	AuthMethod                 *string `json:"auth_method"`
	AuthDomain                 *string `json:"auth_domain"`
	AuthWorkstation            *string `json:"auth_workstation"`
	GRPCURL                    *string `json:"grpc_url"`
	GRPCProtobuf               *string `json:"grpc_protobuf"`
	GRPCBody                   *string `json:"grpc_body"`
	GRPCMetadata               *string `json:"grpc_metadata"`
	GRPCMethod                 *string `json:"grpc_method"`
	GRPCServiceName            *string `json:"grpc_service_name"`
	GRPCEnableTLS              *bool   `json:"grpc_enable_tls"`
	RadiusUsername             *string `json:"radius_username"`
	RadiusPassword             *string `json:"radius_password"`
	RadiusCallingStationID     *string `json:"radius_calling_station_id"`
	RadiusCalledStationID      *string `json:"radius_called_station_id"`
	RadiusSecret               *string `json:"radius_secret"`
	ResendInterval             *int    `json:"resend_interval"`
	PacketSize                 *int    `json:"packet_size"`
	Game                       *string `json:"game"`
	HTTPBodyEncoding           *string `json:"http_body_encoding"`
	Description                *string `json:"description"`
	TLSCA                      *string `json:"tls_ca"`
	TLSCert                    *string `json:"tls_cert"`
	TLSKey                     *string `json:"tls_key"`
	Parent                     *uint   `json:"parent"`
	InvertKeyword              *bool   `json:"invert_keyword"`
	JSONPath                   *string `json:"json_path"`
	ExpectedValue              *string `json:"expected_value"`
	KafkaProducerTopic         *string `json:"kafka_producer_topic"`
	KafkaProducerBrokers       *string `json:"kafka_producer_brokers"`
	KafkaProducerSSL           *bool   `json:"kafka_producer_ssl"`
	KafkaProducerAllowAutoTopicCreation *bool `json:"kafka_producer_allow_auto_topic_creation"`
	KafkaProducerSASLOptions   *string `json:"kafka_producer_sasl_options"`
	KafkaProducerMessage       *string `json:"kafka_producer_message"`
	OAuthClientID              *string `json:"oauth_client_id"`
	OAuthClientSecret          *string `json:"oauth_client_secret"`
	OAuthTokenURL              *string `json:"oauth_token_url"`
	OAuthScopes                *string `json:"oauth_scopes"`
	OAuthAuthMethod            *string `json:"oauth_auth_method"`
	Timeout                    *float64 `json:"timeout"`
	GameDigGivenPortOnly       *bool   `json:"gamedig_given_port_only"`
	MQTTCheckType              *string `json:"mqtt_check_type"`
	RemoteBrowser              *uint   `json:"remote_browser"`
	SNMPOID                    *string `json:"snmp_oid"`
	SNMPVersion                *string `json:"snmp_version"`
	JSONPathOperator           *string `json:"json_path_operator"`
	CacheBust                  *bool   `json:"cache_bust"`
	Conditions                 *string `json:"conditions"`
	RabbitMQNodes              *string `json:"rabbitmq_nodes"`
	RabbitMQUsername           *string `json:"rabbitmq_username"`
	RabbitMQPassword           *string `json:"rabbitmq_password"`
	SMTPSecurity               *string `json:"smtp_security"`
	WSIgnoreSecWebsocketAcceptHeader *bool `json:"ws_ignore_sec_websocket_accept_header"`
	WSSubprotocol              *string `json:"ws_subprotocol"`
	PingCount                  *int    `json:"ping_count"`
	PingNumeric                *bool   `json:"ping_numeric"`
	PingPerRequestTimeout      *int    `json:"ping_per_request_timeout"`
	IPFamily                   *string `json:"ip_family"`
	ManualStatus               *int    `json:"manual_status"`
	OAuthAudience              *string `json:"oauth_audience"`
	MQTTWebsocketPath          *string `json:"mqtt_websocket_path"`
	DomainExpiryNotification   *bool   `json:"domain_expiry_notification"`
	SaveResponse               *bool   `json:"save_response"`
	SaveErrorResponse          *bool   `json:"save_error_response"`
	ResponseMaxLength          *int    `json:"response_max_length"`
	SystemServiceName          *string `json:"system_service_name"`
	Subtype                    *string `json:"subtype"`
	Location                   *string `json:"location"`
	Protocol                   *string `json:"protocol"`
	SNMPV3Username             *string `json:"snmp_v3_username"`
	ExpectedTLSAlert           *string `json:"expected_tls_alert"`
	RetryOnlyOnStatusCodeFailure *bool `json:"retry_only_on_status_code_failure"`
	ScreenshotDelay            *uint   `json:"screenshot_delay"`
	NTPStratumThreshold        *int    `json:"ntp_stratum_threshold"`
	NTPTimeOffsetThreshold     *int    `json:"ntp_time_offset_threshold"`
	NTPRootDispersionThreshold *int    `json:"ntp_root_dispersion_threshold"`
	BearerToken                *string `json:"bearer_token"`
	GameDigToken               *string `json:"gamedig_token"`
	SSHUsername                *string `json:"ssh_username"`
	SSHPassword                *string `json:"ssh_password"`
	SFTPPath                   *string `json:"sftp_path"`
	SSHPrivateKey              *string `json:"ssh_private_key"`
	SSHPassphrase              *string `json:"ssh_passphrase"`
	SSHAuthMethod              *string `json:"ssh_auth_method"`
}

type MonitorUpdateRequest struct {
	Name                       *string `json:"name"`
	Active                     *bool   `json:"active"`
	Interval                   *int    `json:"interval" validate:"omitempty,min=20"`
	URL                        *string `json:"url"`
	Type                       *string `json:"type"`
	Weight                     *int    `json:"weight"`
	Hostname                   *string `json:"hostname"`
	Port                       *int    `json:"port" validate:"omitempty,min=0,max=65535"`
	Keyword                    *string `json:"keyword"`
	MaxRetries                 *int    `json:"maxretries" validate:"omitempty,min=0"`
	IgnoreTLS                  *bool   `json:"ignore_tls"`
	UpsideDown                 *bool   `json:"upside_down"`
	MaxRedirects               *int    `json:"maxredirects" validate:"omitempty,min=0"`
	AcceptedStatusCodes        []string `json:"accepted_statuscodes"`
	DNSResolveType             *string `json:"dns_resolve_type"`
	DNSResolveServer           *string `json:"dns_resolve_server"`
	RetryInterval              *int    `json:"retry_interval" validate:"omitempty,min=20"`
	PushToken                  *string `json:"pushToken"`
	Method                     *string `json:"method"`
	Body                       *string `json:"body"`
	Headers                    *string `json:"headers"`
	BasicAuthUser              *string `json:"basic_auth_user"`
	BasicAuthPass              *string `json:"basic_auth_pass"`
	DockerHost                 *uint   `json:"docker_host"`
	DockerContainer            *string `json:"docker_container"`
	ProxyID                    *uint   `json:"proxy_id"`
	ExpiryNotification         *bool   `json:"expiry_notification"`
	MQTTTopic                  *string `json:"mqtt_topic"`
	MQTTSuccessMessage         *string `json:"mqtt_success_message"`
	MQTTUsername               *string `json:"mqtt_username"`
	MQTTPassword               *string `json:"mqtt_password"`
	DatabaseConnectionString   *string `json:"database_connection_string"`
	DatabaseQuery              *string `json:"database_query"`
	AuthMethod                 *string `json:"auth_method"`
	AuthDomain                 *string `json:"auth_domain"`
	AuthWorkstation            *string `json:"auth_workstation"`
	GRPCURL                    *string `json:"grpc_url"`
	GRPCProtobuf               *string `json:"grpc_protobuf"`
	GRPCBody                   *string `json:"grpc_body"`
	GRPCMetadata               *string `json:"grpc_metadata"`
	GRPCMethod                 *string `json:"grpc_method"`
	GRPCServiceName            *string `json:"grpc_service_name"`
	GRPCEnableTLS              *bool   `json:"grpc_enable_tls"`
	RadiusUsername             *string `json:"radius_username"`
	RadiusPassword             *string `json:"radius_password"`
	RadiusCallingStationID     *string `json:"radius_calling_station_id"`
	RadiusCalledStationID      *string `json:"radius_called_station_id"`
	RadiusSecret               *string `json:"radius_secret"`
	ResendInterval             *int    `json:"resend_interval"`
	PacketSize                 *int    `json:"packet_size"`
	Game                       *string `json:"game"`
	HTTPBodyEncoding           *string `json:"http_body_encoding"`
	Description                *string `json:"description"`
	TLSCA                      *string `json:"tls_ca"`
	TLSCert                    *string `json:"tls_cert"`
	TLSKey                     *string `json:"tls_key"`
	Parent                     *uint   `json:"parent"`
	InvertKeyword              *bool   `json:"invert_keyword"`
	JSONPath                   *string `json:"json_path"`
	ExpectedValue              *string `json:"expected_value"`
	KafkaProducerTopic         *string `json:"kafka_producer_topic"`
	KafkaProducerBrokers       *string `json:"kafka_producer_brokers"`
	KafkaProducerSSL           *bool   `json:"kafka_producer_ssl"`
	KafkaProducerAllowAutoTopicCreation *bool `json:"kafka_producer_allow_auto_topic_creation"`
	KafkaProducerSASLOptions   *string `json:"kafka_producer_sasl_options"`
	KafkaProducerMessage       *string `json:"kafka_producer_message"`
	OAuthClientID              *string `json:"oauth_client_id"`
	OAuthClientSecret          *string `json:"oauth_client_secret"`
	OAuthTokenURL              *string `json:"oauth_token_url"`
	OAuthScopes                *string `json:"oauth_scopes"`
	OAuthAuthMethod            *string `json:"oauth_auth_method"`
	Timeout                    *float64 `json:"timeout"`
	GameDigGivenPortOnly       *bool   `json:"gamedig_given_port_only"`
	MQTTCheckType              *string `json:"mqtt_check_type"`
	RemoteBrowser              *uint   `json:"remote_browser"`
	SNMPOID                    *string `json:"snmp_oid"`
	SNMPVersion                *string `json:"snmp_version"`
	JSONPathOperator           *string `json:"json_path_operator"`
	CacheBust                  *bool   `json:"cache_bust"`
	Conditions                 *string `json:"conditions"`
	RabbitMQNodes              *string `json:"rabbitmq_nodes"`
	RabbitMQUsername           *string `json:"rabbitmq_username"`
	RabbitMQPassword           *string `json:"rabbitmq_password"`
	SMTPSecurity               *string `json:"smtp_security"`
	WSIgnoreSecWebsocketAcceptHeader *bool `json:"ws_ignore_sec_websocket_accept_header"`
	WSSubprotocol              *string `json:"ws_subprotocol"`
	PingCount                  *int    `json:"ping_count"`
	PingNumeric                *bool   `json:"ping_numeric"`
	PingPerRequestTimeout      *int    `json:"ping_per_request_timeout"`
	IPFamily                   *string `json:"ip_family"`
	ManualStatus               *int    `json:"manual_status"`
	OAuthAudience              *string `json:"oauth_audience"`
	MQTTWebsocketPath          *string `json:"mqtt_websocket_path"`
	DomainExpiryNotification   *bool   `json:"domain_expiry_notification"`
	SaveResponse               *bool   `json:"save_response"`
	SaveErrorResponse          *bool   `json:"save_error_response"`
	ResponseMaxLength          *int    `json:"response_max_length"`
	SystemServiceName          *string `json:"system_service_name"`
	Subtype                    *string `json:"subtype"`
	Location                   *string `json:"location"`
	Protocol                   *string `json:"protocol"`
	SNMPV3Username             *string `json:"snmp_v3_username"`
	ExpectedTLSAlert           *string `json:"expected_tls_alert"`
	RetryOnlyOnStatusCodeFailure *bool `json:"retry_only_on_status_code_failure"`
	ScreenshotDelay            *uint   `json:"screenshot_delay"`
	NTPStratumThreshold        *int    `json:"ntp_stratum_threshold"`
	NTPTimeOffsetThreshold     *int    `json:"ntp_time_offset_threshold"`
	NTPRootDispersionThreshold *int    `json:"ntp_root_dispersion_threshold"`
	BearerToken                *string `json:"bearer_token"`
	GameDigToken               *string `json:"gamedig_token"`
	SSHUsername                *string `json:"ssh_username"`
	SSHPassword                *string `json:"ssh_password"`
	SFTPPath                   *string `json:"sftp_path"`
	SSHPrivateKey              *string `json:"ssh_private_key"`
	SSHPassphrase              *string `json:"ssh_passphrase"`
	SSHAuthMethod              *string `json:"ssh_auth_method"`
}

type MonitorResponse struct {
	ID                         uint    `json:"id"`
	Name                       *string `json:"name"`
	Active                     bool    `json:"active"`
	UserID                     *uint   `json:"user_id"`
	Interval                   int     `json:"interval"`
	URL                        *string `json:"url"`
	Type                       *string `json:"type"`
	Weight                     *int    `json:"weight"`
	Hostname                   *string `json:"hostname"`
	Port                       *int    `json:"port"`
	CreatedDate                string  `json:"created_date"`
	Keyword                    *string `json:"keyword"`
	MaxRetries                 int     `json:"maxretries"`
	IgnoreTLS                  bool    `json:"ignore_tls"`
	UpsideDown                 bool    `json:"upside_down"`
	MaxRedirects               int     `json:"maxredirects"`
	AcceptedStatusCodesJSON    string  `json:"accepted_statuscodes_json"`
	DNSResolveType             *string `json:"dns_resolve_type"`
	DNSResolveServer           *string `json:"dns_resolve_server"`
	DNSLastResult              *string `json:"dns_last_result"`
	RetryInterval              int     `json:"retry_interval"`
	PushToken                  *string `json:"pushToken"`
	Method                     string  `json:"method"`
	Body                       *string `json:"body"`
	Headers                    *string `json:"headers"`
	BasicAuthUser              *string `json:"basic_auth_user"`
	BasicAuthPass              *string `json:"basic_auth_pass"`
	DockerHost                 *uint   `json:"docker_host"`
	DockerContainer            *string `json:"docker_container"`
	ProxyID                    *uint   `json:"proxy_id"`
	ExpiryNotification         *bool   `json:"expiry_notification"`
	MQTTTopic                  *string `json:"mqtt_topic"`
	MQTTSuccessMessage         *string `json:"mqtt_success_message"`
	MQTTUsername               *string `json:"mqtt_username"`
	MQTTPassword               *string `json:"mqtt_password"`
	DatabaseConnectionString   *string `json:"database_connection_string"`
	DatabaseQuery              *string `json:"database_query"`
	AuthMethod                 *string `json:"auth_method"`
	AuthDomain                 *string `json:"auth_domain"`
	AuthWorkstation            *string `json:"auth_workstation"`
	GRPCURL                    *string `json:"grpc_url"`
	GRPCProtobuf               *string `json:"grpc_protobuf"`
	GRPCBody                   *string `json:"grpc_body"`
	GRPCMetadata               *string `json:"grpc_metadata"`
	GRPCMethod                 *string `json:"grpc_method"`
	GRPCServiceName            *string `json:"grpc_service_name"`
	GRPCEnableTLS              bool    `json:"grpc_enable_tls"`
	RadiusUsername             *string `json:"radius_username"`
	RadiusPassword             *string `json:"radius_password"`
	RadiusCallingStationID     *string `json:"radius_calling_station_id"`
	RadiusCalledStationID      *string `json:"radius_called_station_id"`
	RadiusSecret               *string `json:"radius_secret"`
	ResendInterval             int     `json:"resend_interval"`
	PacketSize                 int     `json:"packet_size"`
	Game                       *string `json:"game"`
	HTTPBodyEncoding           *string `json:"http_body_encoding"`
	Description                *string `json:"description"`
	TLSCA                      *string `json:"tls_ca"`
	TLSCert                    *string `json:"tls_cert"`
	TLSKey                     *string `json:"tls_key"`
	Parent                     *uint   `json:"parent"`
	InvertKeyword              bool    `json:"invert_keyword"`
	JSONPath                   *string `json:"json_path"`
	ExpectedValue              *string `json:"expected_value"`
	KafkaProducerTopic         *string `json:"kafka_producer_topic"`
	KafkaProducerBrokers       *string `json:"kafka_producer_brokers"`
	KafkaProducerSSL           bool    `json:"kafka_producer_ssl"`
	KafkaProducerAllowAutoTopicCreation bool `json:"kafka_producer_allow_auto_topic_creation"`
	KafkaProducerSASLOptions   *string `json:"kafka_producer_sasl_options"`
	KafkaProducerMessage       *string `json:"kafka_producer_message"`
	OAuthClientID              *string `json:"oauth_client_id"`
	OAuthClientSecret          *string `json:"oauth_client_secret"`
	OAuthTokenURL              *string `json:"oauth_token_url"`
	OAuthScopes                *string `json:"oauth_scopes"`
	OAuthAuthMethod            *string `json:"oauth_auth_method"`
	Timeout                    float64 `json:"timeout"`
	GameDigGivenPortOnly       bool    `json:"gamedig_given_port_only"`
	MQTTCheckType              string  `json:"mqtt_check_type"`
	RemoteBrowser              *uint   `json:"remote_browser"`
	SNMPOID                    *string `json:"snmp_oid"`
	SNMPVersion                *string `json:"snmp_version"`
	JSONPathOperator           *string `json:"json_path_operator"`
	CacheBust                  bool    `json:"cache_bust"`
	Conditions                 string  `json:"conditions"`
	RabbitMQNodes              *string `json:"rabbitmq_nodes"`
	RabbitMQUsername           *string `json:"rabbitmq_username"`
	RabbitMQPassword           *string `json:"rabbitmq_password"`
	SMTPSecurity               *string `json:"smtp_security"`
	WSIgnoreSecWebsocketAcceptHeader bool `json:"ws_ignore_sec_websocket_accept_header"`
	WSSubprotocol              string  `json:"ws_subprotocol"`
	PingCount                  int     `json:"ping_count"`
	PingNumeric                bool    `json:"ping_numeric"`
	PingPerRequestTimeout      int     `json:"ping_per_request_timeout"`
	IPFamily                   *string `json:"ip_family"`
	ManualStatus               *int    `json:"manual_status"`
	OAuthAudience              *string `json:"oauth_audience"`
	MQTTWebsocketPath          *string `json:"mqtt_websocket_path"`
	DomainExpiryNotification   *bool   `json:"domain_expiry_notification"`
	SaveResponse               bool    `json:"save_response"`
	SaveErrorResponse          bool    `json:"save_error_response"`
	ResponseMaxLength          int     `json:"response_max_length"`
	SystemServiceName          *string `json:"system_service_name"`
	Subtype                    *string `json:"subtype"`
	Location                   *string `json:"location"`
	Protocol                   *string `json:"protocol"`
	SNMPV3Username             *string `json:"snmp_v3_username"`
	ExpectedTLSAlert           *string `json:"expected_tls_alert"`
	RetryOnlyOnStatusCodeFailure bool   `json:"retry_only_on_status_code_failure"`
	ScreenshotDelay            uint    `json:"screenshot_delay"`
	NTPStratumThreshold        *int    `json:"ntp_stratum_threshold"`
	NTPTimeOffsetThreshold     *int    `json:"ntp_time_offset_threshold"`
	NTPRootDispersionThreshold *int    `json:"ntp_root_dispersion_threshold"`
	BearerToken                *string `json:"bearer_token"`
	GameDigToken               *string `json:"gamedig_token"`
	SSHUsername                *string `json:"ssh_username"`
	SSHPassword                *string `json:"ssh_password"`
	SFTPPath                   *string `json:"sftp_path"`
	SSHPrivateKey              *string `json:"ssh_private_key"`
	SSHPassphrase              *string `json:"ssh_passphrase"`
	SSHAuthMethod              string                  `json:"ssh_auth_method"`
	Tags                       []MonitorTagResponse    `json:"tags"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
}

type Heartbeat struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Important  bool      `json:"important"`
	MonitorID  uint      `json:"monitor_id"`
	Status     int       `json:"status"`
	Msg        *string   `json:"msg"`
	Time       time.Time `json:"time"`
	Ping       *int64    `json:"ping"`
	Duration   int       `json:"duration"`
	DownCount  int       `json:"down_count"`
	EndTime    *time.Time `json:"end_time"`
	Retries    int       `json:"retries"`
	Response   *string   `json:"response"`
}

func (Heartbeat) TableName() string {
	return "heartbeat"
}

type HeartbeatResponse struct {
	ID         uint   `json:"id"`
	Important  bool   `json:"important"`
	MonitorID  uint   `json:"monitor_id"`
	Status     int    `json:"status"`
	Msg        *string `json:"msg"`
	Time       string `json:"time"`
	Ping       *int64 `json:"ping"`
	Duration   int    `json:"duration"`
	DownCount  int    `json:"down_count"`
	EndTime    *string `json:"end_time"`
	Retries    int    `json:"retries"`
	Response   *string `json:"response"`
}

type Tag struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	CreatedDate time.Time `json:"created_date"`
}

func (Tag) TableName() string {
	return "tag"
}

type TagResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	CreatedDate string `json:"created_date"`
}

type TagCreateRequest struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color" validate:"required"`
}

type TagUpdateRequest struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

type MonitorTag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MonitorID uint      `json:"monitor_id"`
	TagID     uint      `json:"tag_id"`
	Value     *string   `json:"value"`
}

func (MonitorTag) TableName() string {
	return "monitor_tag"
}

type MonitorTagResponse struct {
	MonitorID uint    `json:"monitor_id"`
	TagID     uint    `json:"tag_id"`
	Value     *string `json:"value"`
	Name      string  `json:"name"`
	Color     string  `json:"color"`
}

type MonitorTagCreateRequest struct {
	MonitorID uint   `json:"monitor_id" validate:"required"`
	TagID     uint   `json:"tag_id" validate:"required"`
	Value     *string `json:"value"`
}

type Maintenance struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	UserID          *uint      `json:"user_id"`
	Active          bool       `gorm:"default:true" json:"active"`
	Strategy        string     `gorm:"default:'single'" json:"strategy"`
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	StartTime       *string    `json:"start_time"`
	EndTime         *string    `json:"end_time"`
	Weekdays        string     `gorm:"default:'[]'" json:"weekdays"`
	DaysOfMonth     string     `gorm:"default:'[]'" json:"days_of_month"`
	IntervalDay     *int       `json:"interval_day"`
	Cron            *string    `json:"cron"`
	Timezone        *string    `json:"timezone"`
	Duration        *int       `json:"duration"`
	LastStartDate   *time.Time `json:"last_start_date"`
}

func (Maintenance) TableName() string {
	return "maintenance"
}

type MaintenanceResponse struct {
	ID              uint   `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	UserID          *uint  `json:"user_id"`
	Active          bool   `json:"active"`
	Strategy        string `json:"strategy"`
	StartDate       *string `json:"start_date"`
	EndDate         *string `json:"end_date"`
	StartTime       *string `json:"start_time"`
	EndTime         *string `json:"end_time"`
	Weekdays        string `json:"weekdays"`
	DaysOfMonth     string `json:"days_of_month"`
	IntervalDay     *int   `json:"interval_day"`
	Cron            *string `json:"cron"`
	Timezone        *string `json:"timezone"`
	Duration        *int   `json:"duration"`
	LastStartDate   *string `json:"last_start_date"`
}

type MaintenanceCreateRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	UserID      *uint   `json:"user_id"`
	Active      *bool   `json:"active"`
	Strategy    string  `json:"strategy" validate:"required"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	StartTime   *string `json:"start_time"`
	EndTime     *string `json:"end_time"`
	Weekdays    *string `json:"weekdays"`
	DaysOfMonth *string `json:"days_of_month"`
	IntervalDay *int    `json:"interval_day"`
	Cron        *string `json:"cron"`
	Timezone    *string `json:"timezone"`
	Duration    *int    `json:"duration"`
}

type MaintenanceUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	UserID      *uint   `json:"user_id"`
	Active      *bool   `json:"active"`
	Strategy    *string `json:"strategy"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	StartTime   *string `json:"start_time"`
	EndTime     *string `json:"end_time"`
	Weekdays    *string `json:"weekdays"`
	DaysOfMonth *string `json:"days_of_month"`
	IntervalDay *int    `json:"interval_day"`
	Cron        *string `json:"cron"`
	Timezone    *string `json:"timezone"`
	Duration    *int    `json:"duration"`
}

type MonitorMaintenance struct {
	ID            uint `gorm:"primaryKey" json:"id"`
	MonitorID     uint `json:"monitor_id"`
	MaintenanceID uint `json:"maintenance_id"`
}

func (MonitorMaintenance) TableName() string {
	return "monitor_maintenance"
}

type MonitorMaintenanceResponse struct {
	ID            uint `json:"id"`
	MonitorID     uint `json:"monitor_id"`
	MaintenanceID uint `json:"maintenance_id"`
}

type MonitorMaintenanceCreateRequest struct {
	MonitorID     uint `json:"monitor_id" validate:"required"`
	MaintenanceID uint `json:"maintenance_id" validate:"required"`
}

type StatMinutely struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MonitorID uint      `json:"monitor_id"`
	Timestamp int64     `json:"timestamp"`
	Ping      float64   `json:"ping"`
	Up        int       `json:"up"`
	Down      int       `json:"down"`
	PingMin   float64   `json:"ping_min"`
	PingMax   float64   `json:"ping_max"`
	Extras    *string   `json:"extras"`
}

func (StatMinutely) TableName() string {
	return "stat_minutely"
}

type StatHourly struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MonitorID uint      `json:"monitor_id"`
	Timestamp int64     `json:"timestamp"`
	Ping      float64   `json:"ping"`
	PingMin   float64   `json:"ping_min"`
	PingMax   float64   `json:"ping_max"`
	Up        int       `json:"up"`
	Down      int       `json:"down"`
	Extras    *string   `json:"extras"`
}

func (StatHourly) TableName() string {
	return "stat_hourly"
}

type StatDaily struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MonitorID uint      `json:"monitor_id"`
	Timestamp int64     `json:"timestamp"`
	Ping      float64   `json:"ping"`
	Up        int       `json:"up"`
	Down      int       `json:"down"`
	PingMin   float64   `json:"ping_min"`
	PingMax   float64   `json:"ping_max"`
	Extras    *string   `json:"extras"`
}

func (StatDaily) TableName() string {
	return "stat_daily"
}

type MonitorStatsResponse struct {
	CurrentPing      *float64 `json:"current_ping"`
	AvgPing24h       *float64 `json:"avg_ping_24h"`
	Uptime24h        *float64 `json:"uptime_24h"`
	Uptime30d        *float64 `json:"uptime_30d"`
	Uptime1y         *float64 `json:"uptime_1y"`
}