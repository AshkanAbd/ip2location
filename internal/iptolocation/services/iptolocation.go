package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"ip_location/internal/iptolocation/models"
	"net/http"

	pkgLog "ip_location/pkg/logger"
	pkgMetrics "ip_location/pkg/metrics"
)

const baseUrl = "http://ip-api.com/json"

type IPToLocation struct {
	client *http.Client
	dbConn *gorm.DB
}

func NewIPToLocation(
	client *http.Client,
	dbConn *gorm.DB,
) *IPToLocation {
	return &IPToLocation{
		client: client,
		dbConn: dbConn,
	}
}

func (i *IPToLocation) GetIPInfo(ip string) (models.IPInfo, error) {
	pkgLog.Debug("Fetching ip info for %s ip from db", ip)
	ipInfo, err := i.getFromDB(ip)
	if err != nil {
		pkgLog.Error(err, "Failed to get ip info for %s ip from db", ip)
		return models.IPInfo{}, err
	}

	if ipInfo.IP != "" {
		pkgMetrics.IncreaseIPSourceCount("database")
		pkgLog.Debug("IP info for %s ip found in db", ip)
		return ipInfo, nil
	}

	pkgLog.Debug("Fetching ip info for %s ip from api", ip)
	ipInfo, err = i.getFromApi(ip)
	if err != nil {
		pkgLog.Error(err, "Failed to get ip info for %s ip from api", ip)
		return models.IPInfo{}, err
	}

	if err := i.saveInDB(ipInfo); err != nil {
		pkgLog.Error(err, "Failed to store ip info in db: %s", err.Error())
	}

	pkgMetrics.IncreaseIPSourceCount("api")
	return ipInfo, nil
}

func (i *IPToLocation) getFromApi(ip string) (models.IPInfo, error) {
	resp, err := i.client.Get(fmt.Sprintf("%s/%s", baseUrl, ip))
	if err != nil {
		return models.IPInfo{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.IPInfo{}, err
	}

	var result struct {
		Status  string `json:"status"`
		Country string `json:"country"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return models.IPInfo{}, err
	}

	if result.Status != "success" {
		return models.IPInfo{}, fmt.Errorf(result.Status)
	}

	return models.IPInfo{
		IP:      ip,
		Country: result.Country,
	}, nil
}

func (i *IPToLocation) getFromDB(ip string) (models.IPInfo, error) {
	ipInfo := models.IPInfo{}

	err := i.dbConn.Model(models.IPInfo{}).Where("ip = ?", ip).First(&ipInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.IPInfo{}, nil
		}

		return models.IPInfo{}, err
	}

	return ipInfo, nil
}

func (i *IPToLocation) saveInDB(ipInfo models.IPInfo) error {
	err := i.dbConn.Create(&ipInfo).Error
	if err != nil {
		return err
	}

	return nil
}
