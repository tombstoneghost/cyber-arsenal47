from prompt_toolkit.completion import Completer, Completion, WordCompleter
from prompt_toolkit.document import Document

class ModuleCompleter(Completer):
    def __init__(self, modules):
        self.modules = modules

    def get_completions(self, document, complete_event):
        text = document.text_before_cursor
        parts = text.split()
        
        # If no input or first word, suggest modules
        if len(parts) == 0 or len(parts) == 1:
            for module in self.modules.keys():
                if module.startswith(text):
                    yield Completion(module, start_position=-len(text))
        # If second word, suggest functions of the module
        elif len(parts) > 1:
            module = parts[0]
            if module in self.modules:
                func_text = parts[1] if len(parts) > 1 else ""
                for func in self.modules[module]:
                    if func.startswith(func_text):
                        yield Completion(func, start_position=-len(func_text))


class CommandCompleter(Completer):
    def __init__(self, commands, module_completer):
        self.command_completer = WordCompleter(list(commands.keys()), ignore_case=True)
        self.module_completer = module_completer
        self.in_use_mode = False

    def get_completions(self, document, complete_event):
        text = document.text_before_cursor
        parts = text.split()

        # Toggle between command and module completion based on mode
        if self.in_use_mode:
            # Provide completions for modules and functions after 'use'
            sub_document = Document(text=' '.join(parts), cursor_position=document.cursor_position)
            for completion in self.module_completer.get_completions(sub_document, complete_event):
                yield completion
        else:
            # Provide completions for commands
            if len(parts) == 1:
                for completion in self.command_completer.get_completions(document, complete_event):
                    yield completion
            # Provide completions for modules after the 'use' command
            elif len(parts) > 1:
                command = parts[0]
                if command == 'use' and len(parts) > 1:
                    module_text = parts[1]
                    for module in self.module_completer.modules.keys():
                        if module.startswith(module_text):
                            yield Completion(module, start_position=-len(module_text))
