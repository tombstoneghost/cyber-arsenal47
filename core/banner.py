# Imports
from pyfiglet import Figlet


def print_banner():
    banner_text = "Cyber-Arsenal47"

    f = Figlet(font='slant', width=100)

    print(f.renderText(banner_text))
    print("Developed by @tombstoneghost".center(100))
    print("<- Welcome to Cyber-Arsenal47, The Ultimate Penetration Testing Toolkit ->\n".center(100))
    print("[!] Use 'help' command to list all available commands\n")
