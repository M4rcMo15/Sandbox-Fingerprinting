package collectors

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/m4rcmo15/artefacto/models"
)

// GetPublicIP obtiene la IP pública del sistema
func GetPublicIP() string {
	client := &http.Client{Timeout: 5 * time.Second}
	
	// Intentar varios servicios
	services := []string{
		"https://api.ipify.org",
		"https://icanhazip.com",
		"https://ifconfig.me/ip",
	}
	
	for _, service := range services {
		resp, err := client.Get(service)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		
		return string(body)
	}
	
	return ""
}

// GetGeoLocation obtiene la geolocalización basada en la IP pública
func GetGeoLocation(publicIP string) *models.GeoLocation {
	if publicIP == "" {
		return nil
	}
	
	client := &http.Client{Timeout: 10 * time.Second}
	
	// Usar ip-api.com (gratis, sin API key)
	url := "http://ip-api.com/json/" + publicIP
	
	resp, err := client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	
	var result struct {
		Status      string  `json:"status"`
		Country     string  `json:"country"`
		CountryCode string  `json:"countryCode"`
		Region      string  `json:"regionName"`
		City        string  `json:"city"`
		Lat         float64 `json:"lat"`
		Lon         float64 `json:"lon"`
		ISP         string  `json:"isp"`
		Org         string  `json:"org"`
	}
	
	err = json.Unmarshal(body, &result)
	if err != nil || result.Status != "success" {
		return nil
	}
	
	return &models.GeoLocation{
		Country:      result.Country,
		CountryCode:  result.CountryCode,
		Region:       result.Region,
		City:         result.City,
		Latitude:     result.Lat,
		Longitude:    result.Lon,
		ISP:          result.ISP,
		Organization: result.Org,
	}
}
