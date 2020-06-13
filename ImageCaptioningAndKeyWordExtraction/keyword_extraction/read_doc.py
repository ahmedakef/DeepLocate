import os
import PyPDF2 
#import slate3k as slate
import docx
from odf import text, teletype
from odf.opendocument import load


def read_files(p):
    filename, file_extension = os.path.splitext(p)
    if (file_extension in [".txt" , ".py" , ".cpp" , ".c"] ):
        with open(p, 'r') as file:
            data = file.read()
        return data
    
    
    elif(file_extension == ".pdf"): ## TODO: needs tuning sometimes does not work
        #print("handling pdf:")
        data = []
        t = ""
        with open(p, 'rb') as file:
            pdfReader = PyPDF2.PdfFileReader(file) 
            for i in range (pdfReader.numPages):
                pageObj = pdfReader.getPage(i) 
                data.append(pageObj.extractText())
            t = '\n'.join(data)
        return(t)
    
    
        ## TODO: TEXTRACT STILL NOT FIXED
        ''' 
            if (not len(text)):
                print("HERE")
                text = textract.process(p, method='tesseract', language='eng')
        '''
    
    
        '''
        ### Testing slate3k
            with open(p, 'rb') as file:
                extracted_text = slate.PDF(file)
            return(extracted_text)
        '''
        
    elif(file_extension == ".docx"):
        #print("handling docx:")
        doc = docx.Document(p)
        data = []
        for para in doc.paragraphs:
            data.append(para.text)
        return '\n'.join(data)
        
        
        
    elif(file_extension == ".odt"):
        #print("handling odt:")
        textdoc = load(p)
        data = []
        allparas = textdoc.getElementsByType(text.P)
        for i in range (len(allparas)):
            data.append(teletype.extractText(allparas[i]))
        return '\n'.join(data)
         
        
                
        

          

