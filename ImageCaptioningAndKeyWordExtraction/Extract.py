import sys
import os
import logging
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3' 
logging.disable(logging.WARNING)
import tensorflow
logging.getLogger('tensorflow').disabled = True


from keyword_extraction.text_extract import analyze_text
from image_captioning.generate_folder import analyze_image
import json

#Getting needed files
path = str(sys.argv[1])
doc_files = []
img_files = []
for r, d, f in os.walk(path):
    for file in f:
        if((file.split(".")[-1]) in ['txt' , 'pdf' , 'docx' ,  'py' , "cpp" , "c" , "odt" , 'md']):
            doc_files.append([os.path.join(r, file) , file])
        elif((file.split(".")[-1]) in ['jpg' , 'png']):
            img_files.append([os.path.join(r, file) , file])
#Keyword extraction
doc_dict = analyze_text(doc_files)
#Image Caption
Image_dict = analyze_image(img_files)

print(json.dumps({**doc_dict, **Image_dict}))