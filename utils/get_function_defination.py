import inspect

def extract_params(func):
    signature = inspect.signature(func)

    return [(param.name, param.default) for param in signature.parameters.values()]
