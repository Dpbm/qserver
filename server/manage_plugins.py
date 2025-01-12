from typing import Dict
import os
import json

LockFile = Dict[str, str]

plugins_lock = os.path.join(os.getcwd(), 'plugins.lock')
plugins_folder_path = os.path.join(os.getcwd(), 'plugins')

SUCCESS_RETURNCODE = 0

def get_lock_data() -> LockFile:
    with open(plugins_lock, 'r') as lock_file:
        data = json.load(lock_file)
        return data

def already_added_plugin(lock:LockFile, name:str) -> bool:
    return name in list(lock.keys())

def create_plugins_folder():
    if(not os.path.exists(plugins_folder_path)):
        os.makedirs(plugins_folder_path)

def create_plugins_lock():
    if(not os.path.exists(plugins_lock)):
        with open(plugins_lock, 'w') as lock_file:
            json.dump({}, lock_file)

def download_plugin(url:str, name:str):
    import subprocess

    target_plugin_directory = os.path.join(plugins_folder_path, name)

    result = subprocess.run(['git', 'clone', url, target_plugin_directory])
    if(result.returncode != SUCCESS_RETURNCODE):
        raise Exception("Failed on Download plugin")
    
    # clear unnecessary files
    readme_path = os.path.join(target_plugin_directory, 'README.md')
    dot_git_path = os.path.join(target_plugin_directory, '.git')
    dot_github_path = os.path.join(target_plugin_directory, '.github')
    dot_gitignore_path = os.path.join(target_plugin_directory, '.gitignore')
    license_path = os.path.join(target_plugin_directory, 'LICENSE')

    subprocess.run(["rm","-rf", readme_path, dot_git_path, dot_gitignore_path, dot_github_path, license_path])

def update_lock(lock:LockFile, url:str, name:str):
    with open(plugins_lock, 'w') as lock_file:
        lock[name] = url
        json.dump(lock,lock_file)
