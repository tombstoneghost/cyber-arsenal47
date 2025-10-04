"""
Main CLI handler module.

This file serves as the entry point for the command-line interface (CLI) of the application.
It is responsible for parsing user input, handling commands, and coordinating the application's core functionality.
"""

from prompt_toolkit import PromptSession
from prompt_toolkit.completion import WordCompleter
from time import sleep

from core.utils.generate_module_data import generate_module_data, get_modules
from core.utils.completer import CommandCompleter, ModuleCompleter

from .banner import print_banner
from .commands import COMMANDS

from core.logger import get_logger

class CLI:
    def __init__(self) -> None:
        self.commands = COMMANDS

        self.commands['help'] = self.show_help
        self.commands['exit'] = self.quit
        self.commands['back'] = self.back
        self.commands['use'] = self.use_module
        self.commands['options'] = self.show_options
        self.commands['set'] = self.set_option
        self.commands['unset'] = self.unset_option
        self.commands['run'] = self.run_module

        self.running = True
        self.session = PromptSession()

        self.current_context = ""
        self.loaded_module = ""

        self.modules = generate_module_data()
        self.module_dict = get_modules.get_modules()
        
        self.module_completer = ModuleCompleter(self.module_dict)
        self.command_completer = CommandCompleter({}, self.module_completer)

        self.session = PromptSession(completer=self.command_completer)

        self.logger = get_logger(__name__)


    # Start the CLI
    def start(self):
        print_banner()

        while self.running:
            try:
                if self.current_context:
                    prompt_message = f"cyber-arsenal47 ({self.current_context}) > "
                else:
                    prompt_message = "cyber-arsenal47 > "

                user_input = self.session.prompt(prompt_message).strip().split()
                
                if not user_input:
                    continue

                command = user_input[0]
    
                if command in self.commands:
                    self.command_completer.in_use_mode = ("use" in command)

                    if len(user_input) > 1:
                        try:
                            self.commands[command](user_input[1:])
                        except Exception as e:
                            self.logger.error(f"Error executing command '{command}': {e}")
                    else:
                        try:
                            self.commands[command]()
                        except Exception as e:
                            self.logger.error(f"Incomplete command. Type 'help' for a list of available commands")
                else:
                    self.logger.error(f"Unknown command: {command}. Type 'help' for a list of commands.")

            except KeyboardInterrupt:
                print("\nExiting...")
                self.running = False
            except Exception as e:
                self.logger.error(f"An error occurred: {e}")
    

    # Quit the tool
    def quit(self):
        """Exit the tool"""
        self.running = False
        print("[!] Exiting....")
        sleep(1)
        exit(0)


    # Show Help
    def show_help(self):
        """Displays all available commands"""
        for command, func in self.commands.items():
            print(f"\t{command}: {func.__doc__}")

    
    # Go back to main context
    def back(self):
        """Go back to main context"""
        self.current_context = ""
        self.loaded_module = ""
        self.command_completer.in_use_mode = False

    
    # Use a module
    def use_module(self, args):
        """Load a module for use\tUsage: use <domain/module>"""
        module = args[0]
        self.loaded_module = module
        self.current_context = module
        self.command_completer.in_use_mode = True
        if len(self.current_context) == 0:
            print("[X] Error: Invalid module name.")
            self.current_context = ""
            self.command_completer.in_use_mode = False
        print(f"[!] Module selected: {self.current_context}")


    # Show module options
    def show_options(self):
        """Show options for the currently loaded module"""
        if not self.loaded_module:
            print("[X] Error: No module loaded. Use the 'use' command to load a module.")
            return

        try:
            domain, module = self.loaded_module.split("/")
            option_data = self.modules[domain][module]['params']

            print("\nUse 'set' command to set value to the below parameter")
            for k, v in option_data.items():
                print(f"{str(k).upper():<{20}}\t:\t{v:<{20}}")
            print("")
        except Exception as e:
            print(f"[X] Error processing command: {e}")


    # Set Module Option
    def set_option(self, args):
        """Set value for an option. Usage: set [option] [value]"""
        if not args or len(args) < 2:
            print("[X] Usage: set [option] [value]")
            return
        
        domain, module = self.loaded_module.split("/")
        option = args[0].strip().lower()
        value = " ".join(args[1:]).strip()

        if option not in self.modules[domain][module]['params']:
            print(f"[X] Error: Invalid option '{option}' for module '{self.loaded_module}'")
            return

        self.modules[domain][module]['params'][option] = value
        print(f"[!] Set option '{option}' to '{value}'")


    # UnSet Module Option
    def unset_option(self, args):
        """Unset value for an option. Usage: unset [option]"""
        if not args or len(args) < 1:
            print("[X] Usage: unset [option]")
            return

        domain, module = self.loaded_module.split("/")
        option = args[0].strip()

        if option not in self.modules[domain][module]['params']:
            print(f"[X] Error: Invalid option '{option}' for module '{self.loaded_module}'")
            return

        self.modules[domain][module]['params'][option] = ""
        print(f"[!] Unset option '{option}'")    


    # Run the loaded module
    def run_module(self):
        """Run the currently loaded module"""
        if not self.loaded_module:
            print("[X] Error: No module loaded. Use the 'use' command to load a module.")
            return

        try:
            domain, module = self.loaded_module.split("/")
            module_info = self.modules[domain][module]
            func = module_info['func']
            params = module_info['params']

            print(f"[!] Running module: {self.loaded_module} with parameters: {params}")
            
            params = params.values()
            params = list(params)

            func(*params)

            print(f"[!] Module {self.loaded_module} executed successfully.")
        except Exception as e:
            print(f"[X] Error running module: {e}")

