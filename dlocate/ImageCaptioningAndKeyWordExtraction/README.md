# Overview
Keyword Extraction and image captioning scripts for the project "Advanced File System Search Engine"
# Installation
**IMPORTANT NOTE: You must run python 64 bit version for some libraries to work**
We recommend to use virtualenv for development:
* Create a virtual environment
 *python -m venv env*
* Enable the virtual environment
 *env/activate*
* Install the python dependencies on the virtual environment
*pip install -r requirements.txt*
# Running
* For Keyword Extraction you should run giving text_extract file giving it argument of the path of a document(pdf/docx/odf)  
*text_extract.py PATH/TO/DOC/DOC_NAME.pdf*
* For image captioning you will find to two files you can run one for single images called generate_simage.py or one for folders contains multiple images at once generate_folder.py  
1. For the single image you run script with arugment to the file (jpg or png)  
*generate_simage.py PATH/TO/IMAGE/IMAGE.JPG*
2. For folders you run the script with arugment to images folder
*generate_folder.py PATH/TO/FOLDER*