# Imports
from prompt_toolkit import PromptSession
from prompt_toolkit.completion import WordCompleter
from time import sleep

from utils.generate_module_data import generate_module_data, get_modules
from utils.completer import CommandCompleter, ModuleCompleter

from .banner import print_banner
from .commands import COMMANDS

import sys

class App:
    def __init__(self) -> None:
        self.commands = COMMANDS
        self.commands['help'] = self.show_help
        self.commands['exit'] = self.quit
        self.commands['use'] = self.use_module
        self.commands['back'] = self.back
        self.commands['options'] = self.show_options
        self.commands['set'] = self.set_option
        self.commands['unset'] = self.unset_option
        self.commands['run'] = self.run_module

        self.running = True

        self.current_context = ""
        self.loaded_module = None

        self.modules = generate_module_data()
        self.module_dict = get_modules.get_modules()

        self.module_completer = ModuleCompleter(self.module_dict)

        self.command_completer = CommandCompleter(self.commands, self.module_completer)

        self.session = PromptSession(completer=self.command_completer)

    
    # Show Help
    def show_help(self):
        """Displays all available commands"""
        for command, func in self.commands.items():
            print(f"\t{command}: {func.__doc__}")
    

    # Exit the tool
    def quit(self):
        """Exit"""
        self.running = False
        print("Exiting....")
        
        sys.exit(0)


    # Back Option
    def back(self):
        """Get back to the main menu"""
        self.loaded_module = None
        self.current_context = ""


    # Use a module
    def use_module(self, args):
        """Use a module to run.\tUsage: use [domain]/[module]"""
        self.loaded_module = args[0].strip()
        self.current_context = self.loaded_module

        if len(self.current_context) == 0:
            print("[X] Error: Invalid module name.")
            self.current_context = ""

        print(f"[!] Module selected: {self.current_context}")

    # Show module options
    def show_options(self):
        """List all options for the selected module"""
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
        option, value = str(args[0]).lower(), " ".join(args[1:])

        if domain not in self.modules or module not in self.modules[domain]:
            print(f"[X] Error: Module {self.loaded_module} not found.")
            return

        self.modules[domain][module]['params'][option] = value

        print(f"[!] {option} => {value}")

    # UnSet Module Option
    def unset_option(self, args):
        """UnSet value for an option. Usage: unset [option]"""
        domain, module = self.loaded_module.split("/")
        option = str(args[0]).lower()

        self.modules[domain][module]['params'][option] = ""

        print(f"[!] Unsetted {option}")


    # Run the loaded module
    def run_module(self):
        """Run the loaded module"""
        domain, module = self.loaded_module.split("/")
        func = self.modules[domain][module]['func']
        
        params = self.modules[domain][module]['params'].values()
        params = list(params)

        func(*params)


    # Process Command
    def process_commands(self, user_input: str):
        args = user_input.strip().split()

        if not args:
            return
        
        command = args[0]

        if command in self.commands:
            self.command_completer.in_use_mode = ("use" in command)
            
            if len(args) > 1:
                try:
                    self.commands[command](args[1:])
                except Exception as e:
                    print(f"[X] Error processing command: {e}")
            else:
                try:
                    self.commands[command]()
                except TypeError as e:
                    print(f"[X] Incomplete command. Type 'help' for a list of available commands")
        else:
            print(f"[X] Unknown Command: {command}. Type 'help' for a list of available commands")


    def run(self):
        # Display Banner
        print_banner()

        try:
            # Initialize 
            while self.running:
                prompt_text = f"(ca47 {self.current_context})> "

                user_input = self.session.prompt(prompt_text, completer=self.command_completer)
                self.process_commands(user_input=user_input)
        except KeyboardInterrupt as e:
            self.quit()
