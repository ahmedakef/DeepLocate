import sys
import os
import logging
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3' 

from keyword_extraction.text_extract import analyze_text
from image_captioning.generate_folder import analyze_image
import json


def getPathContent(path):
    #Getting needed files
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

    return Image_dict, doc_dict



if __name__ == "__main__":
    if len(sys.argv) > 1:
        path = str(sys.argv[1])

        Image_dict, doc_dict = getPathContent(path)

        print(json.dumps({**doc_dict, **Image_dict}))

    else:
        # user havn't provided path
        print({})
    