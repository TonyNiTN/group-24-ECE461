"""module docstring"""
from flask import Flask

app = Flask(__name__)


@app.route("/")
def hello():
    """a dummy docstring"""
    return "Hello World-emile!\nHello Ben!\nHello Mimi! yoyo"


@app.route("/index")
def index():
    """a dummy docstring"""
    return "Index Page"


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080)
