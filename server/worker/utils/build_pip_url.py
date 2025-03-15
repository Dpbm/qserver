def build_pip_url(package):
    """
    This helper ensures the correct url to install from
    the official plugins list
    """
    return f"git+https://github.com/quantum-plugins/{package}.git@main"


def pipfy_name(name):
    """
    Fix name to match with the requirements that Pip has
    """
    return name.replace("-", "_")
