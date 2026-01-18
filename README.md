# 🔍 Artefacto - Sandbox Detection & XSS Audit Framework

Advanced sandbox detection and fingerprinting tool with XSS vulnerability auditing capabilities for malware analysis platforms.

## 📋 Overview

**Artefacto** is a comprehensive Red Team framework consisting of two main components:

1. **Agent (Go)**: Lightweight executable that collects detailed system information and injects XSS payloads
2. **Visualizer (Django)**: Web-based dashboard for data visualization and XSS hit monitoring

## ✨ Key Features

### 🎯 Agent Capabilities

**System Intelligence:**
- Sandbox and VM detection with multiple indicators
- EDR/AV detection (Defender, CrowdStrike, SentinelOne, Carbon Black, etc.)
- Analysis tools detection (IDA Pro, x64dbg, Wireshark, Process Monitor, etc.)
- Complete system profiling (OS, CPU, RAM, processes, network connections)
- API hooking detection in critical functions
- Sensitive file discovery
- Desktop screenshot capture
- IP-based geolocation

**XSS Audit Module:**
- Multi-vector XSS payload injection (11 vectors, 27 payloads)
- Hostname, filename, process, registry, window title injection
- DNS queries, HTTP requests, file content, PE metadata, environment variables
- APT signature embedding (APT29, APT28, Lazarus, Emotet, Cobalt Strike, Conti, BlackCat, APT41)
- Automatic callback tracking
- Sandbox vulnerability identification

### 📊 Visualizer Features

**Data Visualization:**
- Execution dashboard with detailed views
- Interactive statistics and charts
- Geographic analysis (countries, cities)
- OS and architecture distribution
- VM vs physical system detection
- EDR/AV product statistics

**XSS Audit Dashboard:**
- Real-time XSS hit monitoring
- Sandbox vulnerability tracking
- Vector success rate analysis
- Payload status tracking (injected/triggered)
- Automatic sandbox identification
- Hit timeline and statistics

## 🚀 Quick Start

### Prerequisites

**Agent:**
- Go 1.20+
- Windows target system

**Visualizer:**
- Python 3.8+
- Django 4.2+
- PostgreSQL/MySQL (production)
- Nginx + Gunicorn (production)

### Installation

#### 1. Clone Repository

```bash
git clone https://github.com/yourusername/artefacto.git
cd artefacto
```

#### 2. Configure Agent

```bash
cd artefacto
cp .env.example .env
```

Edit `.env`:
```env
SERVER_URL=http://your-server.com/api/collect
CALLBACK_SERVER=http://your-server.com
XSS_AUDIT=true
DEBUG=0
TIMEOUT=120s
```

#### 3. Compile Agent

```bash
# Standard build
go build -ldflags="-s -w" -trimpath -o agent.exe

# With APT signatures (recommended)
go build -ldflags="-s -w" -trimpath -o agent.exe
```

#### 4. Setup Visualizer

```bash
cd visualizer
pip install -r requirements.txt
python manage.py migrate
python manage.py createsuperuser
python manage.py runserver 0.0.0.0:8000
```

## 📖 Usage

### Running the Agent

```bash
# Execute on target system
.\agent.exe

# Data is automatically exfiltrated to the server
# XSS payloads are injected if XSS_AUDIT=true
```

### Accessing the Visualizer

- **Main Dashboard:** `http://your-server:8000/`
- **Statistics:** `http://your-server:8000/statistics/`
- **XSS Audit:** `http://your-server:8000/dashboard/`
- **Payloads List:** `http://your-server:8000/payloads/`

## 🎯 XSS Audit Module

### How It Works

1. **Injection Phase:** Agent injects 27 unique XSS payloads across 11 vectors
2. **Analysis Phase:** Sandbox analyzes the binary and generates report
3. **Trigger Phase:** Analyst opens report in browser, XSS executes
4. **Callback Phase:** Browser sends callback to your server
5. **Detection Phase:** Server identifies vulnerable sandbox

### Supported Vectors

| Vector | Description | Success Rate |
|--------|-------------|--------------|
| hostname | System hostname | ⭐⭐⭐⭐⭐ 90% |
| filename | File names | ⭐⭐⭐⭐ 70% |
| process | Process names/args | ⭐⭐⭐⭐ 70% |
| registry | Registry keys/values | ⭐⭐⭐ 50% |
| window | Window titles | ⭐⭐⭐ 50% |
| cmdline | Command line args | ⭐⭐ 30% |
| dns | DNS queries | ⭐ 10% |
| http | HTTP requests | ⭐ 10% |
| file-content | File contents | ⭐ 10% |
| pe-metadata | PE metadata | ⭐ 10% |
| environment | Environment vars | ⭐ 10% |

### APT Signatures

The agent embeds signatures from major APT groups to increase sandbox detection priority:

- **APT29** (Cozy Bear): SUNBURST, TEARDROP
- **APT28** (Fancy Bear): X-Agent, Sofacy
- **Lazarus Group**: WannaCry, HOPLIGHT
- **Emotet**: Banking trojan signatures
- **Cobalt Strike**: Beacon patterns
- **Conti**: Ransomware indicators
- **BlackCat** (ALPHV): Ransomware signatures
- **APT41**: Dual espionage/crime group

## 📁 Project Structure

```
artefacto/
├── artefacto/                  # Go agent
│   ├── collectors/             # Data collection modules
│   ├── config/                 # Configuration management
│   ├── exfil/                  # Data exfiltration
│   ├── models/                 # Data models
│   ├── signatures/             # APT signatures
│   ├── xss/                    # XSS audit module
│   ├── utils/                  # Utilities
│   └── main.go                 # Entry point
│
├── visualizer/                 # Django visualizer
│   ├── collector/              # Main app
│   ├── xss_audit/              # XSS audit app
│   ├── deploy/                 # Deployment scripts
│   └── visualizer/             # Django config
│
└── docs/                       # Documentation
    ├── APT_SIGNATURES_GUIDE.md
    ├── COMPILAR_CON_APT.md
    └── DESPLIEGUE_TEMPLATES_XSS.md
```

## 🔒 Security Considerations

### Authentication

Protect your visualizer with HTTP Basic Auth:

```bash
sudo apt install apache2-utils
sudo htpasswd -c /etc/nginx/auth/.htpasswd username
```

### HTTPS

Enable HTTPS with Let's Encrypt:

```bash
sudo certbot --nginx -d your-domain.com
```

### Firewall

Restrict access to your server:

```bash
sudo ufw allow from YOUR_IP to any port 8000
sudo ufw enable
```

## 📊 Data Collected

The agent collects comprehensive system information:

- **Sandbox Detection:** VM indicators, sandbox artifacts
- **System Info:** OS, architecture, language, timezone
- **Processes:** Running processes, command lines
- **Network:** Active connections, listening ports
- **Security:** EDR/AV products, security drivers
- **Analysis Tools:** Debuggers, disassemblers, monitors
- **Files:** Recent files, sensitive documents
- **Hooks:** Hooked API functions
- **Geolocation:** Country, city, ISP
- **Screenshot:** Desktop capture

## 🛠️ Development

### Agent Development

```bash
cd artefacto
go mod download
go build -o agent.exe
go test ./...
```

### Visualizer Development

```bash
cd visualizer
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python manage.py runserver
```

## 📚 Documentation

- **Agent README:** `artefacto/README.md`
- **Visualizer README:** `visualizer/README.md`
- **APT Signatures:** `APT_SIGNATURES_GUIDE.md`
- **Compilation Guide:** `COMPILAR_CON_APT.md`
- **Deployment Guide:** `DESPLIEGUE_TEMPLATES_XSS.md`

## 🎯 Use Cases

### Red Team Operations
- Identify sandbox environments before payload execution
- Detect EDR/AV products for evasion
- Profile target systems

### Security Research
- Audit sandbox platforms for XSS vulnerabilities
- Analyze malware analysis infrastructure
- Responsible disclosure to vendors

### Penetration Testing
- System reconnaissance
- Security product detection
- Environment fingerprinting

## ⚠️ Responsible Use

This tool is designed for:
- ✅ Authorized Red Team operations
- ✅ Security research with permission
- ✅ Penetration testing engagements
- ✅ Responsible vulnerability disclosure

**NOT for:**
- ❌ Unauthorized access
- ❌ Malicious activities
- ❌ Privacy violations

## 🐛 Known Issues

- Filename injection fails on Windows (invalid characters)
- DNS queries with special characters are rejected
- HTTP requests may be blocked by sandboxes
- File content XSS not analyzed by most sandboxes

## 🔄 Roadmap

- [ ] Linux agent support
- [ ] macOS agent support
- [ ] More XSS payload variants
- [ ] Improved sandbox detection
- [ ] API authentication
- [ ] Multi-user support
- [ ] Export reports (PDF, JSON)

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👥 Authors

- **Marc Monfort** - Initial work

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📧 Contact

For questions, issues, or responsible disclosure:
- Open an issue on GitHub
- Contact: [your-email@example.com]

## ⚠️ Disclaimer

This tool is provided for educational and authorized security testing purposes only. The authors are not responsible for any misuse or damage caused by this tool. Always obtain proper authorization before testing systems you do not own.

---

**Made with ❤️ for the security community**
