def build_pip_url(package):
    """
    This helper ensures the correct url to install from
    the official plugins list
    """
    return f"git+https://github.com/quantum-plugin/{package}.git@main"

def pipfy_name(name):
    return name.replace('-', '_')