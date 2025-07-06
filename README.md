# Cyber-Arsenal47: An Automated Network Assessment Toolkit

> 🔒 A hybrid CLI-based toolkit for automating network reconnaissance, service enumeration, and vulnerability assessment—developed using **Python** and **GoLang**.

---

## ✨ About the Project

**Cyber-Arsenal47** is a powerful, modular, and semi-automated penetration testing framework designed to streamline the process of identifying, scanning, and analyzing various network services.

The tool uses Python for the core CLI logic and orchestration, while service-specific modules are written in GoLang for speed and reliability. This CLI-based tool helps simplify routine pentesting tasks and has been presented at **SecTor Arsenal 2024** held at the **Metro Toronto Convention Centre**.

---

## 🧠 Features

- Modular architecture with plug-and-play domain-specific scanners
- CLI interface for loading, configuring, and running modules
- Auto-pentest module for streamlined assessments
- Exploit-DB integration to identify known exploits
- Built-in modules for FTP, SMB, DNS, and more
- Easy to extend and maintain

---

## 🗂️ Directory Structure

```
├── arsenal/               # GoLang modules compiled to .so files
├── app/                   # Python orchestrator logic
├── files/                 # Exploit DB CSV files
├── utils/                 # Helper functions
├── cyber-arsenal.sh       # Main launcher
└── README.md              # Documentation
```

---

## ⚙️ Installation

### Prerequisites:
- Python 3.8+
- Go 1.18+
- Linux (tested on Kali 2024.1)

### Steps:

```bash
git clone https://github.com/tombstoneghost/cyber-arsenal47.git
cd cyber-arsenal47

# Install Python dependencies
pip install -r requirements.txt

# Launch the CLI
chmod +x cyber-arsenal.sh
./cyber-arsenal.sh
```

---

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

---

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

---

## 📦 Sample Outputs

- **Port Scanner Output:**
  ```
  [+] Found open ports: 21 (FTP), 80 (HTTP), 445 (SMB)
  ```

- **SMB Miner:**
  ```
  [+] Guest access enabled
  [+] Shares enumerated: Public, ADMIN$, C$
  ```

- **Exploit_DB Match:**
  ```
  [+] Apache 2.4.29 - Remote Code Execution
  CVE-2017-5638 matched from Exploit-DB
  ```

- **Auto-Pentest:**
  ```
  [*] Initiating automated scan
  [*] Detected: FTP, SMB
  [*] Running respective modules...
  ```

---

## 🏗 Architecture Overview

The tool follows a CLI → Python Core → Go Module pipeline. Python handles the interface and logic, while the compiled GoLang `.so` modules perform the heavy lifting like enumeration and brute-force attempts.

---

## 🛠 Development Status

The tool is under **active development**. While core modules are stable, some auxiliary and exploit modules are still being refined.  
Bug fixes and optimizations are planned in the upcoming weeks.

---

## 🤝 Contributing

We are open for collaboration!  
If you'd like to contribute:

1. Fork this repo
2. Create your branch: `git checkout -b feature/new-module`
3. Commit and push changes
4. Open a pull request

---

## 📚 License

This project is licensed under the GNU General Public License (GPL).

---

## 🙌 Acknowledgements

- Developed by Simardeep Singh (@tombstoneghost)
- Presented at SecTor Arsenal 2024, Toronto

---

## 💬 Contact

Feel free to reach out or connect for collaboration or suggestions.

🔗 [LinkedIn](https://www.linkedin.com/in/simardeepsingh99/)  
🐙 [GitHub](https://github.com/tombstoneghost)