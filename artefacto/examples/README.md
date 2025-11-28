# Servidor C2 de Ejemplo

Este directorio contiene un servidor C2 simple para recibir datos del agente.

## Instalación

```bash
pip install -r requirements.txt
```

## Uso

```bash
python3 simple_server.py
```

El servidor escuchará en `http://0.0.0.0:8080`

## Endpoints

- `POST /content` - Recibe datos del agente
- `GET /health` - Health check
- `GET /stats` - Estadísticas de reportes recibidos

## Datos Recibidos

Los datos se guardan en el directorio `collected_data/` con el formato:
```
collected_data/
├── HOSTNAME_2025-11-28T10-30-00.json
├── HOSTNAME_2025-11-28T11-15-30.json
└── ...
```

## Configurar el Agente

Asegúrate de que el agente apunte a este servidor:

```bash
export SERVER_URL="http://TU_IP:8080/content"
```

O modifica `config/config.go`:
```go
serverURL = "http://TU_IP:8080/content"
```
