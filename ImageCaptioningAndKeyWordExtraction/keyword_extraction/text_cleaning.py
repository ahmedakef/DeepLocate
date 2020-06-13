import en_core_web_sm
import sys
from nltk.stem.snowball import SnowballStemmer

stemmer = SnowballStemmer(language='english')
nlp = en_core_web_sm.load()
text = ' '.join(sys.argv[1:])
w = nlp(text)

def clean_doc( doc, candidate_pos = ['NOUN' , 'PROPN' , 'VERB' , 'ADJ'] ):
        """Store those words only in cadidate_pos"""
        string = ""
        for sent in doc.sents:
            for token in sent:
                # Store words only with cadidate POS tag
                if token.pos_ in candidate_pos and token.is_stop is False:
                    string += stemmer.stem(token.text.lower()) + " "               
        return string.strip()
print(clean_doc(w))
