# Imports
import os

ARSENAL_PATH = "./arsenal"


def get_modules():
    dir_walk = os.walk(ARSENAL_PATH)

    directories = [x[0][10:] for x in dir_walk]


    directories = list(filter(None, directories))

    modules = {}

    for d in directories:
        module, *func = d.split('/')

        if module not in modules:
            modules[module] = []

        if func:
            modules[module].append(func[0])

    return modules

