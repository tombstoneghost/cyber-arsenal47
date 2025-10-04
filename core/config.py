# Imports
import os
import yaml


# Config Class
class Config:
    def __init__(self, config_file: str = "configs/core.yml"):
        self.config_file = config_file
        self.config_data = self.load_config()

    def load_config(self):
        """Load configuration from a YAML file"""
        if not os.path.exists(self.config_file):
            raise FileNotFoundError(f"Configuration file {self.config_file} not found.")

        with open(self.config_file, 'r') as file:
            try:
                config = yaml.safe_load(file)
                return config
            except yaml.YAMLError as e:
                raise Exception(f"Error parsing YAML file: {e}")

    def get(self, key: str):
        """Get a configuration value"""
        return self.config_data["core"][0].get(key)