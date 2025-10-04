# Imports
import sys
from pathlib import Path

# Add the repo root to sys.path
sys.path.append(str(Path(__file__).resolve().parent.parent))

from core.main import CLI

ca47 = CLI()
ca47.start()
