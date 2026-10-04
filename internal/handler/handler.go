package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"zimbres/uptime-kuma-api/internal/models"
	"zimbres/uptime-kuma-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MonitorHandler struct {
	repo      repository.MonitorRepositoryInterface
	tagRepo   repository.TagRepositoryInterface
	validate  *validator.Validate
}

func NewMonitorHandler(repo repository.MonitorRepositoryInterface, tagRepo repository.TagRepositoryInterface) *MonitorHandler {
	return &MonitorHandler{
		repo:      repo,
		tagRepo:   tagRepo,
		validate:  validator.New(),
	}
}

// CreateMonitor godoc
// @Summary Create a new monitor
// @Description Create a new Uptime Kuma monitor
// @Tags monitors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param monitor body models.MonitorCreateRequest true "Monitor data"
// @Success 201 {object} models.APIResponse{data=models.MonitorResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /monitors [post]
func (h *MonitorHandler) CreateMonitor(c *gin.Context) {
	var req models.MonitorCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Convert request to model
	monitor := h.createRequestToMonitor(&req)

	if err := h.repo.Create(monitor); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create monitor: " + err.Error(),
		})
		return
	}

	response := h.monitorToResponse(monitor)
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Monitor created successfully",
		Data:    response,
	})
}

// GetMonitor godoc
// @Summary Get a monitor by ID
// @Description Get a specific Uptime Kuma monitor by its ID
// @Tags monitors
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Success 200 {object} models.APIResponse{data=models.MonitorResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /monitors/{id} [get]
func (h *MonitorHandler) GetMonitor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	monitor, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response := h.monitorToResponse(monitor)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

// GetMonitors godoc
// @Summary List all monitors
// @Description Get a paginated list of Uptime Kuma monitors
// @Tags monitors
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} models.APIResponse
// @Router /monitors [get]
func (h *MonitorHandler) GetMonitors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	monitors, total, err := h.repo.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch monitors: " + err.Error(),
		})
		return
	}

	var responses []models.MonitorResponse
	for _, m := range monitors {
		responses = append(responses, *h.monitorToResponse(&m))
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Data:    responses,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}

// UpdateMonitor godoc
// @Summary Update a monitor
// @Description Update an existing Uptime Kuma monitor
// @Tags monitors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Param monitor body models.MonitorUpdateRequest true "Monitor data"
// @Success 200 {object} models.APIResponse{data=models.MonitorResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /monitors/{id} [put]
func (h *MonitorHandler) UpdateMonitor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	exists, err := h.repo.Exists(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to check monitor: " + err.Error(),
		})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Monitor not found",
		})
		return
	}

	var req models.MonitorUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Convert to map for updates
	updates := h.updateRequestToMap(&req)

	if err := h.repo.Update(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update monitor: " + err.Error(),
		})
		return
	}

	// Fetch updated monitor
	monitor, _ := h.repo.GetByID(uint(id))
	response := h.monitorToResponse(monitor)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Monitor updated successfully",
		Data:    response,
	})
}

// DeleteMonitor godoc
// @Summary Delete a monitor
// @Description Delete a Uptime Kuma monitor by its ID
// @Tags monitors
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /monitors/{id} [delete]
func (h *MonitorHandler) DeleteMonitor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	exists, err := h.repo.Exists(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to check monitor: " + err.Error(),
		})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Monitor not found",
		})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete monitor: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Monitor deleted successfully",
	})
}

// PauseMonitor godoc
// @Summary Pause a monitor
// @Description Pause (deactivate) a Uptime Kuma monitor
// @Tags monitors
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Success 200 {object} models.APIResponse{data=models.MonitorResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /monitors/{id}/pause [post]
func (h *MonitorHandler) PauseMonitor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	monitor, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	active := false
	if err := h.repo.Update(uint(id), map[string]interface{}{"active": active}); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to pause monitor: " + err.Error(),
		})
		return
	}

	monitor.Active = false
	response := h.monitorToResponse(monitor)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Monitor paused successfully",
		Data:    response,
	})
}

// ResumeMonitor godoc
// @Summary Resume a monitor
// @Description Resume (activate) a Uptime Kuma monitor
// @Tags monitors
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Success 200 {object} models.APIResponse{data=models.MonitorResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /monitors/{id}/resume [post]
func (h *MonitorHandler) ResumeMonitor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	monitor, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	active := true
	if err := h.repo.Update(uint(id), map[string]interface{}{"active": active}); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to resume monitor: " + err.Error(),
		})
		return
	}

	monitor.Active = true
	response := h.monitorToResponse(monitor)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Monitor resumed successfully",
		Data:    response,
	})
}

func (h *MonitorHandler) createRequestToMonitor(req *models.MonitorCreateRequest) *models.Monitor {
	monitor := &models.Monitor{}

	if req.Name != nil {
		monitor.Name = req.Name
	}
	if req.Active != nil {
		monitor.Active = *req.Active
	} else {
		monitor.Active = true
	}
	if req.Interval != nil {
		monitor.Interval = *req.Interval
	} else {
		monitor.Interval = 20
	}
	monitor.URL = req.URL
	monitor.Type = req.Type
	monitor.Weight = req.Weight
	monitor.Hostname = req.Hostname
	monitor.Port = req.Port
	monitor.Keyword = req.Keyword
	if req.MaxRetries != nil {
		monitor.MaxRetries = *req.MaxRetries
	}
	if req.IgnoreTLS != nil {
		monitor.IgnoreTLS = *req.IgnoreTLS
	}
	if req.UpsideDown != nil {
		monitor.UpsideDown = *req.UpsideDown
	}
	if req.MaxRedirects != nil {
		monitor.MaxRedirects = *req.MaxRedirects
	} else {
		monitor.MaxRedirects = 10
	}
	if req.AcceptedStatusCodes != nil {
		codes, _ := json.Marshal(req.AcceptedStatusCodes)
		monitor.AcceptedStatusCodesJSON = string(codes)
	} else {
		monitor.AcceptedStatusCodesJSON = `["200-299"]`
	}
	monitor.DNSResolveType = req.DNSResolveType
	monitor.DNSResolveServer = req.DNSResolveServer
	if req.RetryInterval != nil {
		monitor.RetryInterval = *req.RetryInterval
	} else {
		monitor.RetryInterval = 0
	}
	monitor.PushToken = req.PushToken
	if req.Method != nil {
		monitor.Method = *req.Method
	} else {
		monitor.Method = "GET"
	}
	monitor.Body = req.Body
	monitor.Headers = req.Headers
	monitor.BasicAuthUser = req.BasicAuthUser
	monitor.BasicAuthPass = req.BasicAuthPass
	monitor.DockerHost = req.DockerHost
	monitor.DockerContainer = req.DockerContainer
	monitor.ProxyID = req.ProxyID
	monitor.ExpiryNotification = req.ExpiryNotification
	monitor.MQTTTopic = req.MQTTTopic
	monitor.MQTTSuccessMessage = req.MQTTSuccessMessage
	monitor.MQTTUsername = req.MQTTUsername
	monitor.MQTTPassword = req.MQTTPassword
	monitor.DatabaseConnectionString = req.DatabaseConnectionString
	monitor.DatabaseQuery = req.DatabaseQuery
	monitor.AuthMethod = req.AuthMethod
	monitor.AuthDomain = req.AuthDomain
	monitor.AuthWorkstation = req.AuthWorkstation
	monitor.GRPCURL = req.GRPCURL
	monitor.GRPCProtobuf = req.GRPCProtobuf
	monitor.GRPCBody = req.GRPCBody
	monitor.GRPCMetadata = req.GRPCMetadata
	monitor.GRPCMethod = req.GRPCMethod
	monitor.GRPCServiceName = req.GRPCServiceName
	if req.GRPCEnableTLS != nil {
		monitor.GRPCEnableTLS = *req.GRPCEnableTLS
	}
	monitor.RadiusUsername = req.RadiusUsername
	monitor.RadiusPassword = req.RadiusPassword
	monitor.RadiusCallingStationID = req.RadiusCallingStationID
	monitor.RadiusCalledStationID = req.RadiusCalledStationID
	monitor.RadiusSecret = req.RadiusSecret
	if req.ResendInterval != nil {
		monitor.ResendInterval = *req.ResendInterval
	}
	if req.PacketSize != nil {
		monitor.PacketSize = *req.PacketSize
	} else {
		monitor.PacketSize = 56
	}
	monitor.Game = req.Game
	monitor.HTTPBodyEncoding = req.HTTPBodyEncoding
	monitor.Description = req.Description
	monitor.TLSCA = req.TLSCA
	monitor.TLSCert = req.TLSCert
	monitor.TLSKey = req.TLSKey
	monitor.Parent = req.Parent
	if req.InvertKeyword != nil {
		monitor.InvertKeyword = *req.InvertKeyword
	}
	monitor.JSONPath = req.JSONPath
	monitor.ExpectedValue = req.ExpectedValue
	monitor.KafkaProducerTopic = req.KafkaProducerTopic
	monitor.KafkaProducerBrokers = req.KafkaProducerBrokers
	if req.KafkaProducerSSL != nil {
		monitor.KafkaProducerSSL = *req.KafkaProducerSSL
	}
	if req.KafkaProducerAllowAutoTopicCreation != nil {
		monitor.KafkaProducerAllowAutoTopicCreation = *req.KafkaProducerAllowAutoTopicCreation
	}
	monitor.KafkaProducerSASLOptions = req.KafkaProducerSASLOptions
	monitor.KafkaProducerMessage = req.KafkaProducerMessage
	monitor.OAuthClientID = req.OAuthClientID
	monitor.OAuthClientSecret = req.OAuthClientSecret
	monitor.OAuthTokenURL = req.OAuthTokenURL
	monitor.OAuthScopes = req.OAuthScopes
	monitor.OAuthAuthMethod = req.OAuthAuthMethod
	if req.Timeout != nil {
		monitor.Timeout = *req.Timeout
	}
	if req.GameDigGivenPortOnly != nil {
		monitor.GameDigGivenPortOnly = *req.GameDigGivenPortOnly
	} else {
		monitor.GameDigGivenPortOnly = true
	}
	if req.MQTTCheckType != nil {
		monitor.MQTTCheckType = *req.MQTTCheckType
	} else {
		monitor.MQTTCheckType = "keyword"
	}
	monitor.RemoteBrowser = req.RemoteBrowser
	monitor.SNMPOID = req.SNMPOID
	monitor.SNMPVersion = req.SNMPVersion
	monitor.JSONPathOperator = req.JSONPathOperator
	if req.CacheBust != nil {
		monitor.CacheBust = *req.CacheBust
	}
	if req.Conditions != nil {
		monitor.Conditions = *req.Conditions
	} else {
		monitor.Conditions = "[]"
	}
	monitor.RabbitMQNodes = req.RabbitMQNodes
	monitor.RabbitMQUsername = req.RabbitMQUsername
	monitor.RabbitMQPassword = req.RabbitMQPassword
	monitor.SMTPSecurity = req.SMTPSecurity
	if req.WSIgnoreSecWebsocketAcceptHeader != nil {
		monitor.WSIgnoreSecWebsocketAcceptHeader = *req.WSIgnoreSecWebsocketAcceptHeader
	}
	if req.WSSubprotocol != nil {
		monitor.WSSubprotocol = *req.WSSubprotocol
	}
	if req.PingCount != nil {
		monitor.PingCount = *req.PingCount
	} else {
		monitor.PingCount = 1
	}
	if req.PingNumeric != nil {
		monitor.PingNumeric = *req.PingNumeric
	} else {
		monitor.PingNumeric = true
	}
	if req.PingPerRequestTimeout != nil {
		monitor.PingPerRequestTimeout = *req.PingPerRequestTimeout
	} else {
		monitor.PingPerRequestTimeout = 2
	}
	monitor.IPFamily = req.IPFamily
	monitor.ManualStatus = req.ManualStatus
	monitor.OAuthAudience = req.OAuthAudience
	monitor.MQTTWebsocketPath = req.MQTTWebsocketPath
	monitor.DomainExpiryNotification = req.DomainExpiryNotification
	if req.SaveResponse != nil {
		monitor.SaveResponse = *req.SaveResponse
	}
	if req.SaveErrorResponse != nil {
		monitor.SaveErrorResponse = *req.SaveErrorResponse
	} else {
		monitor.SaveErrorResponse = true
	}
	if req.ResponseMaxLength != nil {
		monitor.ResponseMaxLength = *req.ResponseMaxLength
	} else {
		monitor.ResponseMaxLength = 1024
	}
	monitor.SystemServiceName = req.SystemServiceName
	monitor.Subtype = req.Subtype
	monitor.Location = req.Location
	monitor.Protocol = req.Protocol
	monitor.SNMPV3Username = req.SNMPV3Username
	monitor.ExpectedTLSAlert = req.ExpectedTLSAlert
	if req.RetryOnlyOnStatusCodeFailure != nil {
		monitor.RetryOnlyOnStatusCodeFailure = *req.RetryOnlyOnStatusCodeFailure
	}
	if req.ScreenshotDelay != nil {
		monitor.ScreenshotDelay = *req.ScreenshotDelay
	}
	monitor.NTPStratumThreshold = req.NTPStratumThreshold
	monitor.NTPTimeOffsetThreshold = req.NTPTimeOffsetThreshold
	monitor.NTPRootDispersionThreshold = req.NTPRootDispersionThreshold
	monitor.BearerToken = req.BearerToken
	monitor.GameDigToken = req.GameDigToken
	monitor.SSHUsername = req.SSHUsername
	monitor.SSHPassword = req.SSHPassword
	monitor.SFTPPath = req.SFTPPath
	monitor.SSHPrivateKey = req.SSHPrivateKey
	monitor.SSHPassphrase = req.SSHPassphrase
	if req.SSHAuthMethod != nil {
		monitor.SSHAuthMethod = *req.SSHAuthMethod
	} else {
		monitor.SSHAuthMethod = "password"
	}

	return monitor
}

func (h *MonitorHandler) updateRequestToMap(req *models.MonitorUpdateRequest) map[string]interface{} {
	updates := make(map[string]interface{})

	// Use reflection-like approach to add non-nil fields
	if req.Name != nil {
		updates["name"] = req.Name
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}
	if req.Interval != nil {
		updates["interval"] = *req.Interval
	}
	if req.URL != nil {
		updates["url"] = req.URL
	}
	if req.Type != nil {
		updates["type"] = req.Type
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}
	if req.Hostname != nil {
		updates["hostname"] = req.Hostname
	}
	if req.Port != nil {
		updates["port"] = *req.Port
	}
	if req.Keyword != nil {
		updates["keyword"] = req.Keyword
	}
	if req.MaxRetries != nil {
		updates["maxretries"] = *req.MaxRetries
	}
	if req.IgnoreTLS != nil {
		updates["ignore_tls"] = *req.IgnoreTLS
	}
	if req.UpsideDown != nil {
		updates["upside_down"] = *req.UpsideDown
	}
	if req.MaxRedirects != nil {
		updates["maxredirects"] = *req.MaxRedirects
	}
	if req.AcceptedStatusCodes != nil {
		updates["accepted_statuscodes"] = req.AcceptedStatusCodes
	}
	if req.DNSResolveType != nil {
		updates["dns_resolve_type"] = req.DNSResolveType
	}
	if req.DNSResolveServer != nil {
		updates["dns_resolve_server"] = req.DNSResolveServer
	}
	if req.RetryInterval != nil {
		updates["retry_interval"] = *req.RetryInterval
	}
	if req.PushToken != nil {
		updates["push_token"] = req.PushToken
	}
	if req.Method != nil {
		updates["method"] = req.Method
	}
	if req.Body != nil {
		updates["body"] = req.Body
	}
	if req.Headers != nil {
		updates["headers"] = req.Headers
	}
	if req.BasicAuthUser != nil {
		updates["basic_auth_user"] = req.BasicAuthUser
	}
	if req.BasicAuthPass != nil {
		updates["basic_auth_pass"] = req.BasicAuthPass
	}
	if req.DockerHost != nil {
		updates["docker_host"] = *req.DockerHost
	}
	if req.DockerContainer != nil {
		updates["docker_container"] = req.DockerContainer
	}
	if req.ProxyID != nil {
		updates["proxy_id"] = *req.ProxyID
	}
	if req.ExpiryNotification != nil {
		updates["expiry_notification"] = *req.ExpiryNotification
	}
	if req.MQTTTopic != nil {
		updates["mqtt_topic"] = req.MQTTTopic
	}
	if req.MQTTSuccessMessage != nil {
		updates["mqtt_success_message"] = req.MQTTSuccessMessage
	}
	if req.MQTTUsername != nil {
		updates["mqtt_username"] = req.MQTTUsername
	}
	if req.MQTTPassword != nil {
		updates["mqtt_password"] = req.MQTTPassword
	}
	if req.DatabaseConnectionString != nil {
		updates["database_connection_string"] = req.DatabaseConnectionString
	}
	if req.DatabaseQuery != nil {
		updates["database_query"] = req.DatabaseQuery
	}
	if req.AuthMethod != nil {
		updates["auth_method"] = req.AuthMethod
	}
	if req.AuthDomain != nil {
		updates["auth_domain"] = req.AuthDomain
	}
	if req.AuthWorkstation != nil {
		updates["auth_workstation"] = req.AuthWorkstation
	}
	if req.GRPCURL != nil {
		updates["grpc_url"] = req.GRPCURL
	}
	if req.GRPCProtobuf != nil {
		updates["grpc_protobuf"] = req.GRPCProtobuf
	}
	if req.GRPCBody != nil {
		updates["grpc_body"] = req.GRPCBody
	}
	if req.GRPCMetadata != nil {
		updates["grpc_metadata"] = req.GRPCMetadata
	}
	if req.GRPCMethod != nil {
		updates["grpc_method"] = req.GRPCMethod
	}
	if req.GRPCServiceName != nil {
		updates["grpc_service_name"] = req.GRPCServiceName
	}
	if req.GRPCEnableTLS != nil {
		updates["grpc_enable_tls"] = *req.GRPCEnableTLS
	}
	if req.RadiusUsername != nil {
		updates["radius_username"] = req.RadiusUsername
	}
	if req.RadiusPassword != nil {
		updates["radius_password"] = req.RadiusPassword
	}
	if req.RadiusCallingStationID != nil {
		updates["radius_calling_station_id"] = req.RadiusCallingStationID
	}
	if req.RadiusCalledStationID != nil {
		updates["radius_called_station_id"] = req.RadiusCalledStationID
	}
	if req.RadiusSecret != nil {
		updates["radius_secret"] = req.RadiusSecret
	}
	if req.ResendInterval != nil {
		updates["resend_interval"] = *req.ResendInterval
	}
	if req.PacketSize != nil {
		updates["packet_size"] = *req.PacketSize
	}
	if req.Game != nil {
		updates["game"] = req.Game
	}
	if req.HTTPBodyEncoding != nil {
		updates["http_body_encoding"] = req.HTTPBodyEncoding
	}
	if req.Description != nil {
		updates["description"] = req.Description
	}
	if req.TLSCA != nil {
		updates["tls_ca"] = req.TLSCA
	}
	if req.TLSCert != nil {
		updates["tls_cert"] = req.TLSCert
	}
	if req.TLSKey != nil {
		updates["tls_key"] = req.TLSKey
	}
	if req.Parent != nil {
		updates["parent"] = *req.Parent
	}
	if req.InvertKeyword != nil {
		updates["invert_keyword"] = *req.InvertKeyword
	}
	if req.JSONPath != nil {
		updates["json_path"] = req.JSONPath
	}
	if req.ExpectedValue != nil {
		updates["expected_value"] = req.ExpectedValue
	}
	if req.KafkaProducerTopic != nil {
		updates["kafka_producer_topic"] = req.KafkaProducerTopic
	}
	if req.KafkaProducerBrokers != nil {
		updates["kafka_producer_brokers"] = req.KafkaProducerBrokers
	}
	if req.KafkaProducerSSL != nil {
		updates["kafka_producer_ssl"] = *req.KafkaProducerSSL
	}
	if req.KafkaProducerAllowAutoTopicCreation != nil {
		updates["kafka_producer_allow_auto_topic_creation"] = *req.KafkaProducerAllowAutoTopicCreation
	}
	if req.KafkaProducerSASLOptions != nil {
		updates["kafka_producer_sasl_options"] = req.KafkaProducerSASLOptions
	}
	if req.KafkaProducerMessage != nil {
		updates["kafka_producer_message"] = req.KafkaProducerMessage
	}
	if req.OAuthClientID != nil {
		updates["oauth_client_id"] = req.OAuthClientID
	}
	if req.OAuthClientSecret != nil {
		updates["oauth_client_secret"] = req.OAuthClientSecret
	}
	if req.OAuthTokenURL != nil {
		updates["oauth_token_url"] = req.OAuthTokenURL
	}
	if req.OAuthScopes != nil {
		updates["oauth_scopes"] = req.OAuthScopes
	}
	if req.OAuthAuthMethod != nil {
		updates["oauth_auth_method"] = req.OAuthAuthMethod
	}
	if req.Timeout != nil {
		updates["timeout"] = *req.Timeout
	}
	if req.GameDigGivenPortOnly != nil {
		updates["gamedig_given_port_only"] = *req.GameDigGivenPortOnly
	}
	if req.MQTTCheckType != nil {
		updates["mqtt_check_type"] = req.MQTTCheckType
	}
	if req.RemoteBrowser != nil {
		updates["remote_browser"] = *req.RemoteBrowser
	}
	if req.SNMPOID != nil {
		updates["snmp_oid"] = req.SNMPOID
	}
	if req.SNMPVersion != nil {
		updates["snmp_version"] = req.SNMPVersion
	}
	if req.JSONPathOperator != nil {
		updates["json_path_operator"] = req.JSONPathOperator
	}
	if req.CacheBust != nil {
		updates["cache_bust"] = *req.CacheBust
	}
	if req.Conditions != nil {
		updates["conditions"] = req.Conditions
	}
	if req.RabbitMQNodes != nil {
		updates["rabbitmq_nodes"] = req.RabbitMQNodes
	}
	if req.RabbitMQUsername != nil {
		updates["rabbitmq_username"] = req.RabbitMQUsername
	}
	if req.RabbitMQPassword != nil {
		updates["rabbitmq_password"] = req.RabbitMQPassword
	}
	if req.SMTPSecurity != nil {
		updates["smtp_security"] = req.SMTPSecurity
	}
	if req.WSIgnoreSecWebsocketAcceptHeader != nil {
		updates["ws_ignore_sec_websocket_accept_header"] = *req.WSIgnoreSecWebsocketAcceptHeader
	}
	if req.WSSubprotocol != nil {
		updates["ws_subprotocol"] = req.WSSubprotocol
	}
	if req.PingCount != nil {
		updates["ping_count"] = *req.PingCount
	}
	if req.PingNumeric != nil {
		updates["ping_numeric"] = *req.PingNumeric
	}
	if req.PingPerRequestTimeout != nil {
		updates["ping_per_request_timeout"] = *req.PingPerRequestTimeout
	}
	if req.IPFamily != nil {
		updates["ip_family"] = req.IPFamily
	}
	if req.ManualStatus != nil {
		updates["manual_status"] = *req.ManualStatus
	}
	if req.OAuthAudience != nil {
		updates["oauth_audience"] = req.OAuthAudience
	}
	if req.MQTTWebsocketPath != nil {
		updates["mqtt_websocket_path"] = req.MQTTWebsocketPath
	}
	if req.DomainExpiryNotification != nil {
		updates["domain_expiry_notification"] = *req.DomainExpiryNotification
	}
	if req.SaveResponse != nil {
		updates["save_response"] = *req.SaveResponse
	}
	if req.SaveErrorResponse != nil {
		updates["save_error_response"] = *req.SaveErrorResponse
	}
	if req.ResponseMaxLength != nil {
		updates["response_max_length"] = *req.ResponseMaxLength
	}
	if req.SystemServiceName != nil {
		updates["system_service_name"] = req.SystemServiceName
	}
	if req.Subtype != nil {
		updates["subtype"] = req.Subtype
	}
	if req.Location != nil {
		updates["location"] = req.Location
	}
	if req.Protocol != nil {
		updates["protocol"] = req.Protocol
	}
	if req.SNMPV3Username != nil {
		updates["snmp_v3_username"] = req.SNMPV3Username
	}
	if req.ExpectedTLSAlert != nil {
		updates["expected_tls_alert"] = req.ExpectedTLSAlert
	}
	if req.RetryOnlyOnStatusCodeFailure != nil {
		updates["retry_only_on_status_code_failure"] = *req.RetryOnlyOnStatusCodeFailure
	}
	if req.ScreenshotDelay != nil {
		updates["screenshot_delay"] = *req.ScreenshotDelay
	}
	if req.NTPStratumThreshold != nil {
		updates["ntp_stratum_threshold"] = *req.NTPStratumThreshold
	}
	if req.NTPTimeOffsetThreshold != nil {
		updates["ntp_time_offset_threshold"] = *req.NTPTimeOffsetThreshold
	}
	if req.NTPRootDispersionThreshold != nil {
		updates["ntp_root_dispersion_threshold"] = *req.NTPRootDispersionThreshold
	}
	if req.BearerToken != nil {
		updates["bearer_token"] = req.BearerToken
	}
	if req.GameDigToken != nil {
		updates["gamedig_token"] = req.GameDigToken
	}
	if req.SSHUsername != nil {
		updates["ssh_username"] = req.SSHUsername
	}
	if req.SSHPassword != nil {
		updates["ssh_password"] = req.SSHPassword
	}
	if req.SFTPPath != nil {
		updates["sftp_path"] = req.SFTPPath
	}
	if req.SSHPrivateKey != nil {
		updates["ssh_private_key"] = req.SSHPrivateKey
	}
	if req.SSHPassphrase != nil {
		updates["ssh_passphrase"] = req.SSHPassphrase
	}
	if req.SSHAuthMethod != nil {
		updates["ssh_auth_method"] = req.SSHAuthMethod
	}

	return updates
}

func (h *MonitorHandler) monitorToResponse(m *models.Monitor) *models.MonitorResponse {
	hasMaintenance := false
	if h.repo != nil {
		hasMaintenance, _ = h.repo.HasActiveMaintenance(m.ID)
	}
	return &models.MonitorResponse{
		ID:                         m.ID,
		Name:                       m.Name,
		Active:                     m.Active,
		UserID:                     m.UserID,
		Interval:                   m.Interval,
		URL:                        m.URL,
		Type:                       m.Type,
		Weight:                     m.Weight,
		Hostname:                   m.Hostname,
		Port:                       m.Port,
		CreatedDate:                m.CreatedDate.Format("2006-01-02 15:04:05"),
		Keyword:                    m.Keyword,
		MaxRetries:                 m.MaxRetries,
		IgnoreTLS:                  m.IgnoreTLS,
		UpsideDown:                 m.UpsideDown,
		MaxRedirects:               m.MaxRedirects,
		AcceptedStatusCodesJSON:    m.AcceptedStatusCodesJSON,
		DNSResolveType:             m.DNSResolveType,
		DNSResolveServer:           m.DNSResolveServer,
		DNSLastResult:              m.DNSLastResult,
		RetryInterval:              m.RetryInterval,
		PushToken:                  m.PushToken,
		Method:                     m.Method,
		Body:                       m.Body,
		Headers:                    m.Headers,
		BasicAuthUser:              m.BasicAuthUser,
		BasicAuthPass:              m.BasicAuthPass,
		DockerHost:                 m.DockerHost,
		DockerContainer:            m.DockerContainer,
		ProxyID:                    m.ProxyID,
		ExpiryNotification:         m.ExpiryNotification,
		MQTTTopic:                  m.MQTTTopic,
		MQTTSuccessMessage:         m.MQTTSuccessMessage,
		MQTTUsername:               m.MQTTUsername,
		MQTTPassword:               m.MQTTPassword,
		DatabaseConnectionString:   m.DatabaseConnectionString,
		DatabaseQuery:              m.DatabaseQuery,
		AuthMethod:                 m.AuthMethod,
		AuthDomain:                 m.AuthDomain,
		AuthWorkstation:            m.AuthWorkstation,
		GRPCURL:                    m.GRPCURL,
		GRPCProtobuf:               m.GRPCProtobuf,
		GRPCBody:                   m.GRPCBody,
		GRPCMetadata:               m.GRPCMetadata,
		GRPCMethod:                 m.GRPCMethod,
		GRPCServiceName:            m.GRPCServiceName,
		GRPCEnableTLS:              m.GRPCEnableTLS,
		RadiusUsername:             m.RadiusUsername,
		RadiusPassword:             m.RadiusPassword,
		RadiusCallingStationID:     m.RadiusCallingStationID,
		RadiusCalledStationID:      m.RadiusCalledStationID,
		RadiusSecret:               m.RadiusSecret,
		ResendInterval:             m.ResendInterval,
		PacketSize:                 m.PacketSize,
		Game:                       m.Game,
		HTTPBodyEncoding:           m.HTTPBodyEncoding,
		Description:                m.Description,
		TLSCA:                      m.TLSCA,
		TLSCert:                    m.TLSCert,
		TLSKey:                     m.TLSKey,
		Parent:                     m.Parent,
		InvertKeyword:              m.InvertKeyword,
		JSONPath:                   m.JSONPath,
		ExpectedValue:              m.ExpectedValue,
		KafkaProducerTopic:         m.KafkaProducerTopic,
		KafkaProducerBrokers:       m.KafkaProducerBrokers,
		KafkaProducerSSL:           m.KafkaProducerSSL,
		KafkaProducerAllowAutoTopicCreation: m.KafkaProducerAllowAutoTopicCreation,
		KafkaProducerSASLOptions:   m.KafkaProducerSASLOptions,
		KafkaProducerMessage:       m.KafkaProducerMessage,
		OAuthClientID:              m.OAuthClientID,
		OAuthClientSecret:          m.OAuthClientSecret,
		OAuthTokenURL:              m.OAuthTokenURL,
		OAuthScopes:                m.OAuthScopes,
		OAuthAuthMethod:            m.OAuthAuthMethod,
		Timeout:                    m.Timeout,
		GameDigGivenPortOnly:       m.GameDigGivenPortOnly,
		MQTTCheckType:              m.MQTTCheckType,
		RemoteBrowser:              m.RemoteBrowser,
		SNMPOID:                    m.SNMPOID,
		SNMPVersion:                m.SNMPVersion,
		JSONPathOperator:           m.JSONPathOperator,
		CacheBust:                  m.CacheBust,
		Conditions:                 m.Conditions,
		RabbitMQNodes:              m.RabbitMQNodes,
		RabbitMQUsername:           m.RabbitMQUsername,
		RabbitMQPassword:           m.RabbitMQPassword,
		SMTPSecurity:               m.SMTPSecurity,
		WSIgnoreSecWebsocketAcceptHeader: m.WSIgnoreSecWebsocketAcceptHeader,
		WSSubprotocol:              m.WSSubprotocol,
		PingCount:                  m.PingCount,
		PingNumeric:                m.PingNumeric,
		PingPerRequestTimeout:      m.PingPerRequestTimeout,
		IPFamily:                   m.IPFamily,
		ManualStatus:               m.ManualStatus,
		OAuthAudience:              m.OAuthAudience,
		MQTTWebsocketPath:          m.MQTTWebsocketPath,
		DomainExpiryNotification:   m.DomainExpiryNotification,
		SaveResponse:               m.SaveResponse,
		SaveErrorResponse:          m.SaveErrorResponse,
		ResponseMaxLength:          m.ResponseMaxLength,
		SystemServiceName:          m.SystemServiceName,
		Subtype:                    m.Subtype,
		Location:                   m.Location,
		Protocol:                   m.Protocol,
		SNMPV3Username:             m.SNMPV3Username,
		ExpectedTLSAlert:           m.ExpectedTLSAlert,
		RetryOnlyOnStatusCodeFailure: m.RetryOnlyOnStatusCodeFailure,
		ScreenshotDelay:            m.ScreenshotDelay,
		NTPStratumThreshold:        m.NTPStratumThreshold,
		NTPTimeOffsetThreshold:     m.NTPTimeOffsetThreshold,
		NTPRootDispersionThreshold: m.NTPRootDispersionThreshold,
		BearerToken:                m.BearerToken,
		GameDigToken:               m.GameDigToken,
		SSHUsername:                m.SSHUsername,
		SSHPassword:                m.SSHPassword,
		SFTPPath:                   m.SFTPPath,
		SSHPrivateKey:              m.SSHPrivateKey,
		SSHPassphrase:              m.SSHPassphrase,
		SSHAuthMethod:              m.SSHAuthMethod,
		Tags:                       h.monitorTagsToResponse(m.Tags),
		Maintenance:                hasMaintenance,
	}
}

type HeartbeatHandler struct {
	repo repository.HeartbeatRepositoryInterface
}

func NewHeartbeatHandler(repo repository.HeartbeatRepositoryInterface) *HeartbeatHandler {
	return &HeartbeatHandler{repo: repo}
}

// GetMonitorLastHeartbeat godoc
// @Summary Get monitor's last heartbeat
// @Description Get the last heartbeat status for a specific monitor
// @Tags monitors
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Success 200 {object} models.APIResponse{data=models.HeartbeatResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /monitors/{id}/heartbeat [get]
func (h *HeartbeatHandler) GetMonitorLastHeartbeat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	heartbeat, err := h.repo.GetLastByMonitorID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch heartbeat: " + err.Error(),
		})
		return
	}

	if heartbeat == nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "No heartbeat found for this monitor",
		})
		return
	}

	response := h.heartbeatToResponse(heartbeat)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

// GetMonitorHeartbeats godoc
// @Summary Get monitor's heartbeats
// @Description Get all heartbeats for a specific monitor with pagination
// @Tags monitors
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /monitors/{id}/heartbeats [get]
func (h *HeartbeatHandler) GetMonitorHeartbeats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	heartbeats, err := h.repo.GetByMonitorID(uint(id), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch heartbeats: " + err.Error(),
		})
		return
	}

	// Apply pagination manually since repo returns all
	start := (page - 1) * limit
	end := start + limit
	if start >= len(heartbeats) {
		heartbeats = []models.Heartbeat{}
	} else {
		if end > len(heartbeats) {
			end = len(heartbeats)
		}
		heartbeats = heartbeats[start:end]
	}

	var responses []models.HeartbeatResponse
	for _, hb := range heartbeats {
		responses = append(responses, *h.heartbeatToResponse(&hb))
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Data:    responses,
		Total:   int64(len(heartbeats)),
		Page:    page,
		Limit:   limit,
	})
}

func (h *HeartbeatHandler) heartbeatToResponse(hb *models.Heartbeat) *models.HeartbeatResponse {
	var endTime *string
	if hb.EndTime != nil {
		s := hb.EndTime.Format("2006-01-02 15:04:05")
		endTime = &s
	}
	return &models.HeartbeatResponse{
		ID:        hb.ID,
		Important: hb.Important,
		MonitorID: hb.MonitorID,
		Status:    hb.Status,
		Msg:       hb.Msg,
		Time:      hb.Time.Format("2006-01-02 15:04:05"),
		Ping:      hb.Ping,
		Duration:  hb.Duration,
		DownCount: hb.DownCount,
		EndTime:   endTime,
		Retries:   hb.Retries,
		Response:  hb.Response,
	}
}

type StatsHandler struct {
	repo repository.StatsRepositoryInterface
}

func NewStatsHandler(repo repository.StatsRepositoryInterface) *StatsHandler {
	return &StatsHandler{repo: repo}
}

// GetMonitorStats godoc
// @Summary Get monitor statistics
// @Description Get monitor statistics including current ping, average ping 24h, and uptime percentages
// @Tags monitors
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Success 200 {object} models.APIResponse{data=models.MonitorStatsResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /monitors/{id}/stats [get]
func (h *StatsHandler) GetMonitorStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	stats := models.MonitorStatsResponse{}

	currentPing, err := h.repo.GetCurrentPing(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch current ping: " + err.Error(),
		})
		return
	}
	stats.CurrentPing = currentPing

	avgPing24h, err := h.repo.GetAvgPing24h(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch avg ping 24h: " + err.Error(),
		})
		return
	}
	stats.AvgPing24h = avgPing24h

	uptime24h, err := h.repo.GetUptime24h(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch uptime 24h: " + err.Error(),
		})
		return
	}
	stats.Uptime24h = uptime24h

	uptime30d, err := h.repo.GetUptime30d(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch uptime 30d: " + err.Error(),
		})
		return
	}
	stats.Uptime30d = uptime30d

	uptime1y, err := h.repo.GetUptime1y(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch uptime 1y: " + err.Error(),
		})
		return
	}
	stats.Uptime1y = uptime1y

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    stats,
	})
}

type TagHandler struct {
	repo     repository.TagRepositoryInterface
	validate *validator.Validate
}

func NewTagHandler(repo repository.TagRepositoryInterface) *TagHandler {
	return &TagHandler{
		repo:     repo,
		validate: validator.New(),
	}
}

// CreateTag godoc
// @Summary Create a new tag
// @Description Create a new tag
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tag body models.TagCreateRequest true "Tag data"
// @Success 201 {object} models.APIResponse{data=models.TagResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req models.TagCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	tag := &models.Tag{
		Name:  req.Name,
		Color: req.Color,
	}

	if err := h.repo.Create(tag); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create tag: " + err.Error(),
		})
		return
	}

	response := h.tagToResponse(tag)
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Tag created successfully",
		Data:    response,
	})
}

// GetTags godoc
// @Summary List all tags
// @Description Get all tags
// @Tags tags
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.TagResponse}
// @Failure 500 {object} models.APIResponse
// @Router /tags [get]
func (h *TagHandler) GetTags(c *gin.Context) {
	tags, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch tags: " + err.Error(),
		})
		return
	}

	var responses []models.TagResponse
	for _, t := range tags {
		responses = append(responses, *h.tagToResponse(&t))
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    responses,
	})
}

// GetTag godoc
// @Summary Get a tag by ID
// @Description Get a specific tag by its ID
// @Tags tags
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Success 200 {object} models.APIResponse{data=models.TagResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /tags/{id} [get]
func (h *TagHandler) GetTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid tag ID",
		})
		return
	}

	tag, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response := h.tagToResponse(tag)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

// UpdateTag godoc
// @Summary Update a tag
// @Description Update an existing tag
// @Tags tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Param tag body models.TagUpdateRequest true "Tag data"
// @Success 200 {object} models.APIResponse{data=models.TagResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid tag ID",
		})
		return
	}

	var req models.TagUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Color != nil {
		updates["color"] = *req.Color
	}

	if err := h.repo.Update(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update tag: " + err.Error(),
		})
		return
	}

	tag, _ := h.repo.GetByID(uint(id))
	response := h.tagToResponse(tag)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Tag updated successfully",
		Data:    response,
	})
}

// DeleteTag godoc
// @Summary Delete a tag
// @Description Delete a tag by its ID
// @Tags tags
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid tag ID",
		})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete tag: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Tag deleted successfully",
	})
}

func (h *TagHandler) tagToResponse(t *models.Tag) *models.TagResponse {
	return &models.TagResponse{
		ID:          t.ID,
		Name:        t.Name,
		Color:       t.Color,
		CreatedDate: t.CreatedDate.Format("2006-01-02 15:04:05"),
	}
}

type MonitorTagHandler struct {
	repo repository.MonitorTagRepositoryInterface
}

func NewMonitorTagHandler(repo repository.MonitorTagRepositoryInterface) *MonitorTagHandler {
	return &MonitorTagHandler{repo: repo}
}

// AddMonitorTag godoc
// @Summary Add tag to monitor
// @Description Add a tag to a monitor
// @Tags monitor-tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Param tag body models.MonitorTagCreateRequest true "Monitor tag data"
// @Success 201 {object} models.APIResponse{data=models.MonitorTagResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /monitors/{id}/tags [post]
func (h *MonitorTagHandler) AddMonitorTag(c *gin.Context) {
	monitorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	var req models.MonitorTagCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if req.MonitorID != uint(monitorID) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Monitor ID mismatch",
		})
		return
	}

	// Check if already exists
	existing, _ := h.repo.GetByMonitorAndTag(req.MonitorID, req.TagID)
	if existing != nil {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Error:   "Tag already assigned to this monitor",
		})
		return
	}

	mt := &models.MonitorTag{
		MonitorID: req.MonitorID,
		TagID:     req.TagID,
		Value:     req.Value,
	}

	if err := h.repo.Create(mt); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to add tag to monitor: " + err.Error(),
		})
		return
	}

	response := h.monitorTagToResponse(mt)
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Tag added to monitor successfully",
		Data:    response,
	})
}

// GetMonitorTags godoc
// @Summary Get monitor tags
// @Description Get all tags for a specific monitor
// @Tags monitor-tags
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Success 200 {object} models.APIResponse{data=[]models.MonitorTagResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /monitors/{id}/tags [get]
func (h *MonitorTagHandler) GetMonitorTags(c *gin.Context) {
	monitorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	monitorTags, err := h.repo.GetByMonitorID(uint(monitorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch monitor tags: " + err.Error(),
		})
		return
	}

	var responses []models.MonitorTagResponse
	for _, mt := range monitorTags {
		responses = append(responses, *h.monitorTagToResponse(&mt))
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    responses,
	})
}

// DeleteMonitorTag godoc
// @Summary Remove tag from monitor
// @Description Remove a tag from a monitor
// @Tags monitor-tags
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Param tagId path int true "Tag ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /monitors/{id}/tags/{tagId} [delete]
func (h *MonitorTagHandler) DeleteMonitorTag(c *gin.Context) {
	monitorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	tagID, err := strconv.ParseUint(c.Param("tagId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid tag ID",
		})
		return
	}

	if err := h.repo.Delete(uint(monitorID), uint(tagID)); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to remove tag from monitor: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Tag removed from monitor successfully",
	})
}

func (h *MonitorTagHandler) monitorTagToResponse(mt *models.MonitorTag) *models.MonitorTagResponse {
	return &models.MonitorTagResponse{
		MonitorID: mt.MonitorID,
		TagID:     mt.TagID,
		Value:     mt.Value,
	}
}

func (h *MonitorHandler) monitorTagsToResponse(tags []models.MonitorTag) []models.MonitorTagResponse {
	if tags == nil || len(tags) == 0 {
		return []models.MonitorTagResponse{}
	}
	responses := make([]models.MonitorTagResponse, len(tags))
	for i, t := range tags {
		name := ""
		color := ""
		if h.tagRepo != nil {
			tag, err := h.tagRepo.GetByID(t.TagID)
			if err == nil && tag != nil {
				name = tag.Name
				color = tag.Color
			}
		}
		responses[i] = models.MonitorTagResponse{
			MonitorID: t.MonitorID,
			TagID:     t.TagID,
			Value:     t.Value,
			Name:      name,
			Color:     color,
		}
	}
	return responses
}

type MaintenanceHandler struct {
	repo     repository.MaintenanceRepositoryInterface
	validate *validator.Validate
}

func NewMaintenanceHandler(repo repository.MaintenanceRepositoryInterface) *MaintenanceHandler {
	return &MaintenanceHandler{
		repo:     repo,
		validate: validator.New(),
	}
}

// CreateMaintenance godoc
// @Summary Create a new maintenance
// @Description Create a new maintenance window
// @Tags maintenances
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param maintenance body models.MaintenanceCreateRequest true "Maintenance data"
// @Success 201 {object} models.APIResponse{data=models.MaintenanceResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /maintenances [post]
func (h *MaintenanceHandler) CreateMaintenance(c *gin.Context) {
	var req models.MaintenanceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	m := &models.Maintenance{
		Title:       req.Title,
		Description: req.Description,
		UserID:      req.UserID,
		Strategy:    req.Strategy,
		Weekdays:    "[]",
		DaysOfMonth: "[]",
	}

	if req.Active != nil {
		m.Active = *req.Active
	} else {
		m.Active = true
	}
	if req.StartDate != nil {
		if t, err := time.Parse("2006-01-02 15:04:05", *req.StartDate); err == nil {
			m.StartDate = &t
		}
	}
	if req.EndDate != nil {
		if t, err := time.Parse("2006-01-02 15:04:05", *req.EndDate); err == nil {
			m.EndDate = &t
		}
	}
	m.StartTime = req.StartTime
	m.EndTime = req.EndTime
	if req.Weekdays != nil {
		m.Weekdays = *req.Weekdays
	}
	if req.DaysOfMonth != nil {
		m.DaysOfMonth = *req.DaysOfMonth
	}
	m.IntervalDay = req.IntervalDay
	m.Cron = req.Cron
	m.Timezone = req.Timezone
	m.Duration = req.Duration

	if err := h.repo.Create(m); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create maintenance: " + err.Error(),
		})
		return
	}

	response := h.maintenanceToResponse(m)
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Maintenance created successfully",
		Data:    response,
	})
}

// GetMaintenances godoc
// @Summary List all maintenances
// @Description Get a paginated list of maintenances
// @Tags maintenances
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} models.APIResponse
// @Router /maintenances [get]
func (h *MaintenanceHandler) GetMaintenances(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	maintenances, total, err := h.repo.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch maintenances: " + err.Error(),
		})
		return
	}

	var responses []models.MaintenanceResponse
	for _, m := range maintenances {
		responses = append(responses, *h.maintenanceToResponse(&m))
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    map[string]interface{}{
			"data":  responses,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetMaintenance godoc
// @Summary Get a maintenance by ID
// @Description Get a specific maintenance by its ID
// @Tags maintenances
// @Produce json
// @Security BearerAuth
// @Param id path int true "Maintenance ID"
// @Success 200 {object} models.APIResponse{data=models.MaintenanceResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /maintenances/{id} [get]
func (h *MaintenanceHandler) GetMaintenance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid maintenance ID",
		})
		return
	}

	m, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response := h.maintenanceToResponse(m)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

// UpdateMaintenance godoc
// @Summary Update a maintenance
// @Description Update an existing maintenance
// @Tags maintenances
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Maintenance ID"
// @Param maintenance body models.MaintenanceUpdateRequest true "Maintenance data"
// @Success 200 {object} models.APIResponse{data=models.MaintenanceResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /maintenances/{id} [put]
func (h *MaintenanceHandler) UpdateMaintenance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid maintenance ID",
		})
		return
	}

	var req models.MaintenanceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.UserID != nil {
		updates["user_id"] = *req.UserID
	}
	if req.Active != nil {
		updates["active"] = *req.Active
	}
	if req.Strategy != nil {
		updates["strategy"] = *req.Strategy
	}
	if req.StartDate != nil {
		if t, err := time.Parse("2006-01-02 15:04:05", *req.StartDate); err == nil {
			updates["start_date"] = t
		}
	}
	if req.EndDate != nil {
		if t, err := time.Parse("2006-01-02 15:04:05", *req.EndDate); err == nil {
			updates["end_date"] = t
		}
	}
	if req.StartTime != nil {
		updates["start_time"] = *req.StartTime
	}
	if req.EndTime != nil {
		updates["end_time"] = *req.EndTime
	}
	if req.Weekdays != nil {
		updates["weekdays"] = *req.Weekdays
	}
	if req.DaysOfMonth != nil {
		updates["days_of_month"] = *req.DaysOfMonth
	}
	if req.IntervalDay != nil {
		updates["interval_day"] = *req.IntervalDay
	}
	if req.Cron != nil {
		updates["cron"] = *req.Cron
	}
	if req.Timezone != nil {
		updates["timezone"] = *req.Timezone
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}

	if err := h.repo.Update(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update maintenance: " + err.Error(),
		})
		return
	}

	m, _ := h.repo.GetByID(uint(id))
	response := h.maintenanceToResponse(m)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Maintenance updated successfully",
		Data:    response,
	})
}

// DeleteMaintenance godoc
// @Summary Delete a maintenance
// @Description Delete a maintenance by its ID
// @Tags maintenances
// @Produce json
// @Security BearerAuth
// @Param id path int true "Maintenance ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /maintenances/{id} [delete]
func (h *MaintenanceHandler) DeleteMaintenance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid maintenance ID",
		})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete maintenance: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Maintenance deleted successfully",
	})
}

func (h *MaintenanceHandler) maintenanceToResponse(m *models.Maintenance) *models.MaintenanceResponse {
	var startDate, endDate, lastStartDate *string
	if m.StartDate != nil {
		s := m.StartDate.Format("2006-01-02 15:04:05")
		startDate = &s
	}
	if m.EndDate != nil {
		s := m.EndDate.Format("2006-01-02 15:04:05")
		endDate = &s
	}
	if m.LastStartDate != nil {
		s := m.LastStartDate.Format("2006-01-02 15:04:05")
		lastStartDate = &s
	}
	return &models.MaintenanceResponse{
		ID:              m.ID,
		Title:           m.Title,
		Description:     m.Description,
		UserID:          m.UserID,
		Active:          m.Active,
		Strategy:        m.Strategy,
		StartDate:       startDate,
		EndDate:         endDate,
		StartTime:       m.StartTime,
		EndTime:         m.EndTime,
		Weekdays:        m.Weekdays,
		DaysOfMonth:     m.DaysOfMonth,
		IntervalDay:     m.IntervalDay,
		Cron:            m.Cron,
		Timezone:        m.Timezone,
		Duration:        m.Duration,
		LastStartDate:   lastStartDate,
	}
}

type MonitorMaintenanceHandler struct {
	repo repository.MonitorMaintenanceRepositoryInterface
}

func NewMonitorMaintenanceHandler(repo repository.MonitorMaintenanceRepositoryInterface) *MonitorMaintenanceHandler {
	return &MonitorMaintenanceHandler{repo: repo}
}

// AddMonitorMaintenance godoc
// @Summary Add maintenance to monitor
// @Description Associate a maintenance window with a monitor
// @Tags monitor-maintenances
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Param maintenance body models.MonitorMaintenanceCreateRequest true "Monitor maintenance data"
// @Success 201 {object} models.APIResponse{data=models.MonitorMaintenanceResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /monitors/{id}/maintenances [post]
func (h *MonitorMaintenanceHandler) AddMonitorMaintenance(c *gin.Context) {
	monitorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	var req models.MonitorMaintenanceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if req.MonitorID != uint(monitorID) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Monitor ID mismatch",
		})
		return
	}

	mm := &models.MonitorMaintenance{
		MonitorID:     req.MonitorID,
		MaintenanceID: req.MaintenanceID,
	}

	if err := h.repo.Create(mm); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to add maintenance to monitor: " + err.Error(),
		})
		return
	}

	response := h.monitorMaintenanceToResponse(mm)
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Maintenance added to monitor successfully",
		Data:    response,
	})
}

// GetMonitorMaintenances godoc
// @Summary Get monitor maintenances
// @Description Get all maintenances for a specific monitor
// @Tags monitor-maintenances
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Success 200 {object} models.APIResponse{data=[]models.MonitorMaintenanceResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /monitors/{id}/maintenances [get]
func (h *MonitorMaintenanceHandler) GetMonitorMaintenances(c *gin.Context) {
	monitorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	mms, err := h.repo.GetByMonitorID(uint(monitorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch monitor maintenances: " + err.Error(),
		})
		return
	}

	var responses []models.MonitorMaintenanceResponse
	for _, mm := range mms {
		responses = append(responses, *h.monitorMaintenanceToResponse(&mm))
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    responses,
	})
}

// DeleteMonitorMaintenance godoc
// @Summary Remove maintenance from monitor
// @Description Remove a maintenance association from a monitor
// @Tags monitor-maintenances
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Param maintenanceId path int true "Maintenance ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /monitors/{id}/maintenances/{maintenanceId} [delete]
func (h *MonitorMaintenanceHandler) DeleteMonitorMaintenance(c *gin.Context) {
	monitorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid monitor ID",
		})
		return
	}

	maintenanceID, err := strconv.ParseUint(c.Param("maintenanceId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid maintenance ID",
		})
		return
	}

	if err := h.repo.Delete(uint(monitorID), uint(maintenanceID)); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to remove maintenance from monitor: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Maintenance removed from monitor successfully",
	})
}

func (h *MonitorMaintenanceHandler) monitorMaintenanceToResponse(mm *models.MonitorMaintenance) *models.MonitorMaintenanceResponse {
	return &models.MonitorMaintenanceResponse{
		ID:            mm.ID,
		MonitorID:     mm.MonitorID,
		MaintenanceID: mm.MaintenanceID,
	}
}