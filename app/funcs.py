# Imports
import ctypes

from utils.get_modules import get_modules

# Ctypes Libary Configuration
arsenal = ctypes.cdll.LoadLibrary('./arsenal/arsenal.so')


# List all available domains
def list_domains():
    """List all available domains"""
    modules = get_modules()

    domains = modules.keys() 

    for domain in domains:
        print(f"[->] {domain}")


# List all available modules in a domain
def list_modules(args):
    """List all available modules in a domain\tUsage: module [domain] """

    all_modules = get_modules()

    try:
        modules = all_modules[args[0]]

        for m in modules:
            print(f"[>] {args[0]}/{m}")
    except KeyError:
        print("[X] Domain not found. Type 'domain' command to list all available domains")


# List all available modules
def list_all_modules():
    """List all available modules and their domains"""
    all_modules = get_modules()

    print("\n[!] Listing all available modules in the arsenal\n")

    for d in all_modules.keys():
        print(f"[$] Domain - {d}")
        print("*" * 25)
        for m in all_modules[d]:
            print(f"[>] {d}/{m}")
        print()



# Port Scanner
def port_scanner(target: str = "", port_start: str = "1", port_end: str = "10000", nmap: str = "True", host_discovery: str = "False", protocol:str = "tcp"):
    """Scan for open ports on the provided target(s)"""
    
    tool = arsenal.PortScanner
    tool.argtypes = [ctypes.c_int, ctypes.POINTER(ctypes.POINTER(ctypes.c_char)), ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]

    targets = target.split(',')

    target_count = len(targets)
    target_strings = (ctypes.POINTER(ctypes.c_char) * target_count)(
        *[ctypes.create_string_buffer(s.encode('utf-8')) for s in targets]
    )

    tool(target_count, target_strings, port_start.encode('utf-8'), port_end.encode('utf-8'), nmap.encode('utf-8'), host_discovery.encode('utf-8'), protocol.encode('utf-8'))


# Directory Scanner
def directory_scanner(target: str = "", wordlist: str = "/usr/share/wordlists/dirbuster/directory-list-2.3-medium.txt", error_codes: str = "404", extensions: str = "", file_only: str="False", threads: str = "10"):
    """Checks for open directories"""
    print("[!] Running Directory Scanner")

    tool = arsenal.DirectoryScanner
    tool.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int, ctypes.POINTER(ctypes.POINTER(ctypes.c_char)), ctypes.c_int, ctypes.POINTER(ctypes.POINTER(ctypes.c_char)), ctypes.c_char_p, ctypes.c_int]

    error_code = error_codes.split(',')
    error_code_count = len(error_code)
    error_code_strings = (ctypes.POINTER(ctypes.c_char) * error_code_count)(
        *[ctypes.create_string_buffer(s.encode('utf-8')) for s in error_code]
    )

    extension = extensions.split(',')
    extension_count = len(extension)
    extension_strings = (ctypes.POINTER(ctypes.c_char) * extension_count)(
        *[ctypes.create_string_buffer(s.encode('utf-8')) for s in extension]
    )

    thread_count = int(threads)

    tool(target.encode('utf-8'), wordlist.encode('utf-8'), error_code_count, error_code_strings, extension_count, extension_strings, file_only.encode('utf-8'), thread_count)


# VHost Scanner
def vhost_scanner(target: str = "", wordlist: str = "/usr/share/seclists/Discovery/DNS/subdomains-top1million-20000.txt", ssl: str="False", ignore_body_length: str="", threads: str = "10"):
    """Checks for available Vhost"""
    print("[!] Running Vhost Scanner")
    
    tool = arsenal.VhostScanner
    tool.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]

    threads_count = int(threads)

    tool(target.encode('utf-8'), wordlist.encode('utf-8'), ssl.encode('utf-8'), ignore_body_length.encode('utf-8'), threads_count)


# DB Miner
def db_miner(db_Type: str = "", host : str = "", port : str = "", username : str = "", password : str = "", db_Name : str = ""):
    """"DB Miner"""
    print("[!] Running DB Miner")

    tool = arsenal.DBMiner
    tool.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]

    tool(db_Type.encode('utf-8'), host.encode('utf-8'), port.encode('utf-8'), username.encode('utf-8'), password.encode('utf-8'), db_Name.encode('utf-8'))


# Hash Crack
def hash_crack(hash: str = "", wordlist: str = "/usr/share/wordlists/rockyou.txt"):
    """Cracks the provided hash"""
    print("[!] Running Hash Cracker")

    tool = arsenal.HashCracker
    tool.argtypes = [ctypes.c_char_p, ctypes.c_char_p]

    tool(hash.encode('utf-8'), wordlist.encode('utf-8'))


# DNS Cache Snooping
def dns_cache_snooping(domain: str = "", dns_server: str = ""):
    """DNS Cache Snooping"""
    print("[!] Running DNS Cache Snooping")
    
    tool = arsenal.DNSCacheSnooping
    tool.argtypes = [ctypes.c_char_p, ctypes.c_char_p]

    tool(domain.encode('utf-8'), dns_server.encode('utf-8'))


# FTP Login
def ftp_login(target: str = "", user_file: str = "", pass_file: str = "", worker_count: str = "5"):
    """FTP Login Brute Force"""
    print("[!] Running FTP Brute Force")

    tool = arsenal.FTPBruteForcer
    tool.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]

    workerCount = int(worker_count)

    tool(target.encode('utf-8'), user_file.encode('utf-8'), pass_file.encode('utf-8'), workerCount)


# FTP Miner
def ftp_miner(target: str = "", username: str = "", password: str = ""):
    """FTP Miner to mine all the data from FTP Server"""
    print("[!] Running FTP Miner")

    tool = arsenal.FTPMiner
    tool.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]

    tool(target.encode('utf-8'), username.encode('utf-8'), password.encode('utf-8'))

# SMB Login
def smb_login(target: str = "", user_file: str = "", pass_file: str = "", worker_count: str = "5"):
    """SMB Login Brute Force"""
    print("[!] Running SMB Brute Force")

    tool = arsenal.SMBBruteForce
    tool.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]

    workerCount = int(worker_count)

    tool(target.encode('utf-8'), user_file.encode('utf-8'), pass_file.encode('utf-8'), workerCount)

# SMB Miner
def smb_miner(target: str = "", username: str = "", password: str = ""):
    """SMB Miner to mine all the data from SMB Share"""
    print("[!] Running SMB Miner")

    tool = arsenal.SMBMiner
    tool.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]

    tool(target.encode('utf-8'), username.encode('utf-8'), password.encode('utf-8'))


# Exploit DB Search
def exploit_db_search(query: str = ""):
    """Search Exploit DB for the provided query"""
    print("[!] Running ExploitDB Search")

    tool = arsenal.ExploitDBSearch
    tool.argtypes = [ctypes.c_char_p]

    tool(query.encode('utf-8'))


# Auto Pentest
def auto_pentest(target: str = "", port_start: str = "1", port_end: str = "10000", nmap: str = "True", host_discovery: str = "False", protocol:str = "tcp",
	dirScanWordlist: str = "/usr/share/wordlists/dirb/common.txt", error_codes: str = "404", extensions: str = "", file_only: str="False",
	vHostScanWordlist: str = "/usr/share/seclists/Discovery/DNS/subdomains-top1million-20000.txt", ssl: str="False", ignore_body_length: str="", 
    user_file: str = "", pass_file: str = "", worker_count: str = "5"):
    """Auto Pentest module to scan and find all possible vulnerabilities"""
    print("[!] Running Auto Pentest")

    tool = arsenal.AutoPentest
    tool.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p,
                     ctypes.c_char_p, ctypes.c_int, ctypes.POINTER(ctypes.POINTER(ctypes.c_char)), ctypes.c_int, ctypes.POINTER(ctypes.POINTER(ctypes.c_char)), ctypes.c_char_p,
                     ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
    
    error_code = error_codes.split(',')
    error_code_count = len(error_code)
    error_code_strings = (ctypes.POINTER(ctypes.c_char) * error_code_count)(
        *[ctypes.create_string_buffer(s.encode('utf-8')) for s in error_code]
    )

    extension = extensions.split(',')
    extension_count = len(extension)
    extension_strings = (ctypes.POINTER(ctypes.c_char) * extension_count)(
        *[ctypes.create_string_buffer(s.encode('utf-8')) for s in extension]
    )

    workerCount = int(worker_count)

    tool(target.encode('utf-8'), port_start.encode('utf-8'), port_end.encode('utf-8'), nmap.encode('utf-8'), host_discovery.encode('utf-8'), protocol.encode('utf-8'),
		dirScanWordlist.encode('utf-8'), error_code_count, error_code_strings, extension_count, extension_strings, file_only.encode('utf-8'),
		vHostScanWordlist.encode('utf-8'), ssl.encode('utf-8'), ignore_body_length.encode('utf-8'), user_file.encode('utf-8'), pass_file.encode('utf-8'), workerCount)
    
