package repository

import (
	"encoding/json"
	"errors"
	"time"

	"zimbres/uptime-kuma-api/internal/models"

	"gorm.io/gorm"
)

// MonitorRepositoryInterface defines the interface for monitor repository
type MonitorRepositoryInterface interface {
	Create(monitor *models.Monitor) error
	GetByID(id uint) (*models.Monitor, error)
	GetAll(page, limit int) ([]models.Monitor, int64, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	Exists(id uint) (bool, error)
	HasActiveMaintenance(monitorID uint) (bool, error)
}

// HeartbeatRepositoryInterface defines the interface for heartbeat repository
type HeartbeatRepositoryInterface interface {
	GetLastByMonitorID(monitorID uint) (*models.Heartbeat, error)
	GetByMonitorID(monitorID uint, limit int) ([]models.Heartbeat, error)
}

// TagRepositoryInterface defines the interface for tag repository
type TagRepositoryInterface interface {
	Create(tag *models.Tag) error
	GetByID(id uint) (*models.Tag, error)
	GetAll() ([]models.Tag, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
}

// MonitorTagRepositoryInterface defines the interface for monitor tag repository
type MonitorTagRepositoryInterface interface {
	Create(mt *models.MonitorTag) error
	GetByMonitorID(monitorID uint) ([]models.MonitorTag, error)
	GetByMonitorAndTag(monitorID, tagID uint) (*models.MonitorTag, error)
	Delete(monitorID, tagID uint) error
}

// MaintenanceRepositoryInterface defines the interface for maintenance repository
type MaintenanceRepositoryInterface interface {
	Create(m *models.Maintenance) error
	GetByID(id uint) (*models.Maintenance, error)
	GetAll(page, limit int) ([]models.Maintenance, int64, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
}

// MonitorMaintenanceRepositoryInterface defines the interface for monitor maintenance repository
type MonitorMaintenanceRepositoryInterface interface {
	Create(mm *models.MonitorMaintenance) error
	GetByMonitorID(monitorID uint) ([]models.MonitorMaintenance, error)
	Delete(monitorID, maintenanceID uint) error
}

// StatsRepositoryInterface defines the interface for stats repository
type StatsRepositoryInterface interface {
	GetCurrentPing(monitorID uint) (*float64, error)
	GetAvgPing24h(monitorID uint) (*float64, error)
	GetUptime24h(monitorID uint) (*float64, error)
	GetUptime30d(monitorID uint) (*float64, error)
	GetUptime1y(monitorID uint) (*float64, error)
}

type MonitorRepository struct {
	db *gorm.DB
}

func NewMonitorRepository(db *gorm.DB) *MonitorRepository {
	return &MonitorRepository{db: db}
}

func (r *MonitorRepository) Create(monitor *models.Monitor) error {
	if monitor.AcceptedStatusCodesJSON == "" {
		monitor.AcceptedStatusCodesJSON = `["200-299"]`
	}
	if monitor.Method == "" {
		monitor.Method = "GET"
	}
	if monitor.Conditions == "" {
		monitor.Conditions = "[]"
	}
	if monitor.MQTTCheckType == "" {
		monitor.MQTTCheckType = "keyword"
	}
	if monitor.WSSubprotocol == "" {
		monitor.WSSubprotocol = ""
	}
	if monitor.SSHAuthMethod == "" {
		monitor.SSHAuthMethod = "password"
	}
	if monitor.SNMPVersion == nil || *monitor.SNMPVersion == "" {
		val := "2c"
		monitor.SNMPVersion = &val
	}
	if monitor.HTTPBodyEncoding == nil {
		encoding := "json"
		monitor.HTTPBodyEncoding = &encoding
	}
	return r.db.Create(monitor).Error
}

func (r *MonitorRepository) GetByID(id uint) (*models.Monitor, error) {
	var monitor models.Monitor
	err := r.db.Preload("Tags").First(&monitor, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("monitor not found")
		}
		return nil, err
	}
	return &monitor, nil
}

func (r *MonitorRepository) GetAll(page, limit int) ([]models.Monitor, int64, error) {
	var monitors []models.Monitor
	var total int64

	err := r.db.Model(&models.Monitor{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = r.db.Preload("Tags").Offset(offset).Limit(limit).Find(&monitors).Error
	if err != nil {
		return nil, 0, err
	}

	return monitors, total, nil
}

func (r *MonitorRepository) Update(id uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return errors.New("no updates provided")
	}

	// Convert accepted_statuscodes array to JSON string if provided
	if codes, ok := updates["accepted_statuscodes"].([]string); ok {
		jsonCodes, _ := json.Marshal(codes)
		updates["accepted_statuscodes_json"] = string(jsonCodes)
		delete(updates, "accepted_statuscodes")
	}

	return r.db.Model(&models.Monitor{}).Where("id = ?", id).Updates(updates).Error
}

func (r *MonitorRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Monitor{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("monitor not found")
	}
	return nil
}

func (r *MonitorRepository) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Monitor{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *MonitorRepository) HasActiveMaintenance(monitorID uint) (bool, error) {
	var count int64
	now := time.Now()
	err := r.db.Table("monitor_maintenance mm").
		Joins("JOIN maintenance m ON m.id = mm.maintenance_id").
		Where("mm.monitor_id = ? AND m.active = ?", monitorID, true).
		Where("(m.strategy = 'single' AND m.start_date <= ? AND (m.end_date IS NULL OR m.end_date >= ?)) OR m.strategy != 'single'", now, now).
		Count(&count).Error
	return count > 0, err
}

type HeartbeatRepository struct {
	db *gorm.DB
}

func NewHeartbeatRepository(db *gorm.DB) *HeartbeatRepository {
	return &HeartbeatRepository{db: db}
}

func (r *HeartbeatRepository) GetLastByMonitorID(monitorID uint) (*models.Heartbeat, error) {
	var heartbeat models.Heartbeat
	err := r.db.Where("monitor_id = ?", monitorID).Order("time DESC").First(&heartbeat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &heartbeat, nil
}

func (r *HeartbeatRepository) GetByMonitorID(monitorID uint, limit int) ([]models.Heartbeat, error) {
	var heartbeats []models.Heartbeat
	err := r.db.Where("monitor_id = ?", monitorID).Order("time DESC").Limit(limit).Find(&heartbeats).Error
	return heartbeats, err
}

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *TagRepository) GetByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tag not found")
		}
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) GetAll() ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Find(&tags).Error
	return tags, err
}

func (r *TagRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.Tag{}).Where("id = ?", id).Updates(updates).Error
}

func (r *TagRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Tag{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("tag not found")
	}
	return nil
}

type MonitorTagRepository struct {
	db *gorm.DB
}

func NewMonitorTagRepository(db *gorm.DB) *MonitorTagRepository {
	return &MonitorTagRepository{db: db}
}

func (r *MonitorTagRepository) Create(mt *models.MonitorTag) error {
	return r.db.Create(mt).Error
}

func (r *MonitorTagRepository) GetByMonitorID(monitorID uint) ([]models.MonitorTag, error) {
	var monitorTags []models.MonitorTag
	err := r.db.Where("monitor_id = ?", monitorID).Find(&monitorTags).Error
	return monitorTags, err
}

func (r *MonitorTagRepository) GetByMonitorAndTag(monitorID, tagID uint) (*models.MonitorTag, error) {
	var mt models.MonitorTag
	err := r.db.Where("monitor_id = ? AND tag_id = ?", monitorID, tagID).First(&mt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &mt, nil
}

func (r *MonitorTagRepository) Delete(monitorID, tagID uint) error {
	result := r.db.Where("monitor_id = ? AND tag_id = ?", monitorID, tagID).Delete(&models.MonitorTag{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("monitor tag not found")
	}
	return nil
}

type MaintenanceRepository struct {
	db *gorm.DB
}

func NewMaintenanceRepository(db *gorm.DB) *MaintenanceRepository {
	return &MaintenanceRepository{db: db}
}

func (r *MaintenanceRepository) Create(m *models.Maintenance) error {
	if m.Weekdays == "" {
		m.Weekdays = "[]"
	}
	if m.DaysOfMonth == "" {
		m.DaysOfMonth = "[]"
	}
	return r.db.Create(m).Error
}

func (r *MaintenanceRepository) GetByID(id uint) (*models.Maintenance, error) {
	var m models.Maintenance
	err := r.db.First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("maintenance not found")
		}
		return nil, err
	}
	return &m, nil
}

func (r *MaintenanceRepository) GetAll(page, limit int) ([]models.Maintenance, int64, error) {
	var maintenances []models.Maintenance
	var total int64

	err := r.db.Model(&models.Maintenance{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = r.db.Offset(offset).Limit(limit).Find(&maintenances).Error
	if err != nil {
		return nil, 0, err
	}

	return maintenances, total, nil
}

func (r *MaintenanceRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.Maintenance{}).Where("id = ?", id).Updates(updates).Error
}

func (r *MaintenanceRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Maintenance{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("maintenance not found")
	}
	return nil
}

type MonitorMaintenanceRepository struct {
	db *gorm.DB
}

func NewMonitorMaintenanceRepository(db *gorm.DB) *MonitorMaintenanceRepository {
	return &MonitorMaintenanceRepository{db: db}
}

func (r *MonitorMaintenanceRepository) Create(mm *models.MonitorMaintenance) error {
	return r.db.Create(mm).Error
}

func (r *MonitorMaintenanceRepository) GetByMonitorID(monitorID uint) ([]models.MonitorMaintenance, error) {
	var mms []models.MonitorMaintenance
	err := r.db.Where("monitor_id = ?", monitorID).Find(&mms).Error
	return mms, err
}

func (r *MonitorMaintenanceRepository) Delete(monitorID, maintenanceID uint) error {
	result := r.db.Where("monitor_id = ? AND maintenance_id = ?", monitorID, maintenanceID).Delete(&models.MonitorMaintenance{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("monitor maintenance not found")
	}
	return nil
}

type StatsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// GetCurrentPing gets the current ping from the latest minutely stat
func (r *StatsRepository) GetCurrentPing(monitorID uint) (*float64, error) {
	var stat models.StatMinutely
	err := r.db.Where("monitor_id = ?", monitorID).Order("timestamp DESC").First(&stat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &stat.Ping, nil
}

// GetAvgPing24h gets average ping from last 24 hours of hourly stats
func (r *StatsRepository) GetAvgPing24h(monitorID uint) (*float64, error) {
	now := time.Now().Unix()
	cutoff := now - 24*60*60 // 24 hours ago

	var avgPing float64
	err := r.db.Model(&models.StatHourly{}).
		Where("monitor_id = ? AND timestamp >= ?", monitorID, cutoff).
		Select("AVG(ping)").Scan(&avgPing).Error
	if err != nil {
		return nil, err
	}
	if avgPing == 0 {
		return nil, nil
	}
	return &avgPing, nil
}

// GetUptime24h gets uptime percentage from last 24 hours of hourly stats
func (r *StatsRepository) GetUptime24h(monitorID uint) (*float64, error) {
	now := time.Now().Unix()
	cutoff := now - 24*60*60 // 24 hours ago

	var totalUp, totalDown int64
	err := r.db.Model(&models.StatHourly{}).
		Where("monitor_id = ? AND timestamp >= ?", monitorID, cutoff).
		Select("SUM(up), SUM(down)").Row().Scan(&totalUp, &totalDown)
	if err != nil {
		return nil, err
	}

	if totalUp+totalDown == 0 {
		return nil, nil
	}

	uptime := float64(totalUp) / float64(totalUp+totalDown) * 100
	return &uptime, nil
}

// GetUptime30d gets uptime percentage from last 30 days of daily stats
func (r *StatsRepository) GetUptime30d(monitorID uint) (*float64, error) {
	now := time.Now().Unix()
	cutoff := now - 30*24*60*60 // 30 days ago

	var totalUp, totalDown int64
	err := r.db.Model(&models.StatDaily{}).
		Where("monitor_id = ? AND timestamp >= ?", monitorID, cutoff).
		Select("SUM(up), SUM(down)").Row().Scan(&totalUp, &totalDown)
	if err != nil {
		return nil, err
	}

	if totalUp+totalDown == 0 {
		return nil, nil
	}

	uptime := float64(totalUp) / float64(totalUp+totalDown) * 100
	return &uptime, nil
}

// GetUptime1y gets uptime percentage from last 365 days of daily stats
func (r *StatsRepository) GetUptime1y(monitorID uint) (*float64, error) {
	now := time.Now().Unix()
	cutoff := now - 365*24*60*60 // 365 days ago

	var totalUp, totalDown int64
	err := r.db.Model(&models.StatDaily{}).
		Where("monitor_id = ? AND timestamp >= ?", monitorID, cutoff).
		Select("SUM(up), SUM(down)").Row().Scan(&totalUp, &totalDown)
	if err != nil {
		return nil, err
	}

	if totalUp+totalDown == 0 {
		return nil, nil
	}

	uptime := float64(totalUp) / float64(totalUp+totalDown) * 100
	return &uptime, nil
}