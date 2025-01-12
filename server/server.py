import logging
from flask import Flask, request, abort, json
from werkzeug.exceptions import HTTPException
from markupsafe import escape

from manage_plugins import already_added_plugin, download_plugin, create_plugins_folder, update_lock, get_lock_data, create_plugins_lock

BAD_REQUEST_CODE = 400
SERVER_ERROR_CODE = 500
CREATED_CODE = 201

app = Flask(__name__)
logger = logging.getLogger(__name__)

@app.errorhandler(HTTPException)
def handle_bad_request(error):
    response = error.get_response()
    response.data = json.dumps({
        "name": error.name,
        "description": error.description,
    })
    response.content_type = "application/json"
    return response


@app.route("/addPlugin/")
def add_plugin():
    try:
        plugin_url = escape(request.args.get("url"))
        name = escape(request.args.get("name"))
        
        invalid_parameter = lambda param: param == 'None' or not param

        create_plugins_lock()

        lock_data = get_lock_data()

        if (invalid_parameter(plugin_url) or invalid_parameter(name)):
            abort(BAD_REQUEST_CODE, description="To Add a plugin you must provide a valid name and url")
        
        if(already_added_plugin(lock_data, name)):
            abort(BAD_REQUEST_CODE, description="You already have a plugin with this name. Please, try a differnt one!")

        create_plugins_folder()
        download_plugin(plugin_url, name)
        update_lock(lock_data, plugin_url, name)
    
        return {
            "description": "Plguin Added Successfully"
        }, CREATED_CODE

    except Exception as error:
        logger.error("Add plugin error: " + str(error))
        abort(SERVER_ERROR_CODE, description="It wasn't possible to add your plugin. Please, try again later!")

