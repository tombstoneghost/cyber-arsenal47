# Cyber-Arsenal47: Modular Automated Penetration Testing Toolkit

> **Cyber-Arsenal47** is a modular, automation-focused penetration testing toolkit built with Go and Python. It leverages high-performance Go modules compiled into a shared library (`arsenal.so`), which are dynamically loaded and orchestrated by a Python CLI driver.


## 🚩 Project Overview

Cyber-Arsenal47 streamlines network reconnaissance, enumeration, and exploitation by combining the speed of Go with the flexibility of Python. The toolkit is designed for extensibility, automation, and ease of use, making it suitable for both red and blue teams.

## ✨ Core Features

- **Modular Design:** Plug-and-play Go modules for scanners, miners, exploits, and automation tasks.
- **Python CLI Interface:** Interactive command-line interface with dynamic module discovery and autocompletion.
- **AutoPentest Module:** Orchestrates multi-stage scanning and exploitation workflows with a single command.
- **Logging & Results:** Each run generates detailed logs and summary files in the `log/` directory.
- **Automatic Go Module Building:** Build all Go modules into `arsenal.so` using `scripts/build_modules.sh`.
- **Configurable Modules:** YAML files in `configs/` allow easy customization of service mappings and module options.

## 🏗 Architecture Overview

- **Go Modules:** Each scanner, miner, or exploit is implemented as a Go module.
- **Shared Library (`arsenal.so`):** All Go modules are compiled into a single shared object using cgo.
- **Python Driver:** The CLI (in Python) loads `arsenal.so` via `ctypes`, exposing all modules as callable functions.
- **Workflow:**  
  `Python CLI` → `arsenal.so` (Go modules) → Results/logs


## ⚙️ Setup & Usage

### Prerequisites

- **Go** (>= 1.18)
- **Python 3** (>= 3.8)
- **gcc** (for cgo compilation)
- **Linux** (tested on Kali 2024.1)

### Installation & Build

```bash
git clone https://github.com/tombstoneghost/cyber-arsenal47.git
cd cyber-arsenal47

# Install Python dependencies
pip install -r requirements.txt

# Build Go modules as a shared library
chmod +x build_modules.sh
./build_modules.sh

# Run the CLI
./cyber-arsenal47.sh
```

## 🖥 Sample CLI Output

```
$ ./cyber-arsenal.sh
[!] Building Arsenal Modules
[+] Modules build successfully
[sudo] password for user: 
   ______      __                    ___                               ____ _______
  / ____/_  __/ /_  ___  _____      /   |  _____________  ____  ____ _/ / // /__  /
 / /   / / / / __ \/ _ \/ ___/_____/ /| | / ___/ ___/ _ \/ __ \/ __ `/ / // /_ / / 
/ /___/ /_/ / /_/ /  __/ /  /_____/ ___ |/ /  (__  )  __/ / / / /_/ / /__  __// /  
\____/\__, /_.___/\___/_/        /_/  |_/_/  /____/\___/_/ /_/\__,_/_/  /_/  /_/   
     /____/                                                                         

          <- Welcome to Cyber-Arsenal47, The Ultimate Penetration Testing Toolkit ->
```
```
# Use a scanner module
use scanners/port_scanner
run

# Use the automated pentest module
use automate/auto_pentest
set target <IP>
run
```

## 📁 Folder Structure
```
cyber-arsenal47/
├── arsenal/      # Go modules (scanners, miners, exploits, automation)
│   └── arsenal.go
├── core/         # Python CLI logic and utilities
├── cmd/          # CLI entry point (cli.py)
├── configs/      # YAML configuration files
├── log/          # Output logs and summary files
├── scripts/      # Build scripts (build_modules.sh)
├── [requirements.txt](http://_vscodecontentref_/4)
└── [README.md](http://_vscodecontentref_/5)
```

## 🚀 Modules Overview

### ✅ Implemented Modules
- `scanners/port_scanner`
- `scanners/ftp_login`
- `scanners/smb_login`
- `auxiliary/ftp_miner`
- `auxiliary/smb_miner`
- `auxiliary/dns_snooper`
- `exploit/exploit_db`
- `automate/auto_pentest`

### 🧪 In Progress / Planned
- `auxiliary/rdp_miner`
- `scanners/ldap_login`
- `scanners/mssql_login`
- `auxiliary/snmp_miner`
- `scanners/websocket_scanner`
- `exploit/nfs_enum`
- `auxiliary/smtp_miner`


## 📦 Logging & Results

- **Log Files**: Each module and AutoPentest run generates a timestamped log in `log/` (e.g., log/`auto_pentest_target_YYYYMMDD_HHMMSS.log`).
- **Summary Files**: Human-readable summaries are saved alongside logs (e.g., `log/auto_pentest_target_YYYYMMDD_HHMMSS.txt`).
- **Configurable Output**: Log and summary file locations are managed automatically per run.

## 🤝 Contributing

**We welcome contributions!**
To add a new module or feature:

- **Go Modules**:
  - Add your scanner/miner/exploit in arsenal/.
  - Register it in arsenal.go for export.
  - Follow the existing module structure and documentation.

- **Python CLI**:
  - Integrate new modules in core/ and update cmd/cli.py as needed.
  - Ensure new commands and options are documented.

- **Configs**:
  - Add or update YAML files in configs/ for service-port mappings or module options.

**Contribution Steps:**

1. Fork the repo and create a feature branch.
2. Add or update modules/configs.
3. Test your changes.
4. Open a pull request with a clear description.

## 🛣️ Roadmap
- New modules: RDP miner, LDAP login, MSSQL login, SNMP miner, WebSocket scanner, NFS enum, SMTP miner.
- Continuous Integration (CI) for automated testing and builds.
- Enhanced reporting and sample output files.
- Improved CLI UX and error handling.


## 📚 License

This project is licensed under the GNU General Public License (GPL).

## 🙌 Acknowledgements

- Developed by **Simardeep Singh* (@tombstoneghost)
- Presented at SecTor Arsenal 2024, Toronto

## 💬 Contact

Feel free to reach out or connect for collaboration or suggestions.

🔗 [LinkedIn](https://www.linkedin.com/in/simardeepsingh99/)  
🐙 [GitHub](https://github.com/tombstoneghost)
