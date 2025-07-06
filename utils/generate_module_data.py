from utils import get_modules, get_function_defination
from app import funcs

def generate_module_data():
    module_data = {}

    raw_data = get_modules.get_modules()

    # Initialize
    for domain, modules in raw_data.items():
        module_data[domain] = {}

        for module in modules:
            module_data[domain][module] = {}
            module_data[domain][module]['params'] = {}
            
            func_attr = getattr(funcs, module)
            func_params = get_function_defination.extract_params(func_attr)
            
            for params in func_params:
                module_data[domain][module]['params'][params[0]] = params[1]
            
            module_data[domain][module]['func'] = func_attr

    return module_data

    
