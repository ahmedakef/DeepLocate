# Overview
Keyword Extraction and image captioning scripts for the project "Advanced File System Search Engine"

# Installation
**IMPORTANT NOTE: You must run python 64 bit version for some libraries to work**
We recommend to use virtualenv for development:

Create a virtual environment
```
python -m venv env
```

Enable the virtual environment
```
env/activate
```

Install the python dependencies on the virtual environment
```
pip install -r requirements.txt
python -m spacy download en_core_web_sm
```

# Running
run `Extract.py` and give it directory containing files/images 
```
python -W ignore Extract.py  PATH/TO/DOC/DIRECTORY
```
* `-W ignore` to ignore the warnings
