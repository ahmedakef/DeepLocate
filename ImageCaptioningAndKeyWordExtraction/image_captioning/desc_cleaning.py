import en_core_web_sm
import sys
from nltk.stem.snowball import SnowballStemmer
def clean_text( t, candidate_pos = ['NOUN' , 'PROPN' , 'VERB' , 'ADJ'] ):
		stemmer = SnowballStemmer(language='english')
		nlp = en_core_web_sm.load()
		doc = nlp(t)
		string = ""
		for sent in doc.sents:
			for token in sent:
				if token.pos_ in candidate_pos and token.is_stop is False:
					string += stemmer.stem(token.text.lower()) + " "
		return string.strip()