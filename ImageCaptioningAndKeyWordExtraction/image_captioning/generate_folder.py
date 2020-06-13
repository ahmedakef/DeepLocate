import os
import sys

# ignore tensorflow wanrnings
import logging
logging.basicConfig(level=logging.INFO)

stderr = sys.stderr
sys.stderr = open(os.devnull, 'w')
import keras
sys.stderr = stderr

from pickle import load
from numpy import argmax
from keras.preprocessing.sequence import pad_sequences
#from keras.applications.vgg16 import VGG16
#from keras.applications.vgg16 import preprocess_input

from keras.applications.inception_v3 import InceptionV3
from keras.applications.inception_v3 import preprocess_input

from keras.preprocessing.image import load_img
from keras.preprocessing.image import img_to_array

from keras.models import Model
from keras.models import load_model

from image_captioning.desc_cleaning import clean_text




# extract features from each photo in the directory
def extract_features(img_files):
	# load the model
	#model = VGG16()
	model = InceptionV3()

	# re-structure the model
	model.layers.pop()
	model = Model(inputs=model.inputs, outputs=model.layers[-1].output)
	# extract features from each photo
	features = dict()
	for filename,name in img_files:
		# load an image from file
		image = load_img(filename, target_size=(299, 299))
		# convert the image pixels to a numpy array
		image = img_to_array(image)
		# reshape data for the model
		image = image.reshape((1, image.shape[0], image.shape[1], image.shape[2]))
		# prepare the image for the VGG model
		image = preprocess_input(image)
		# get features
		feature = model.predict(image, verbose=0)
		# get image id
		image_id = name.split('.')[0]
		# store feature
		features[image_id] = feature
	return features
 
# map an integer to a word
def word_for_id(integer, tokenizer):
	for word, index in tokenizer.word_index.items():
		if index == integer:
			return word
	return None
 
# generate a description for an image
def generate_desc(model, tokenizer, photo, max_length):
	# seed the generation process
	in_text = 'startseq'
	# iterate over the whole length of the sequence
	for i in range(max_length):
		# integer encode input sequence
		sequence = tokenizer.texts_to_sequences([in_text])[0]
		# pad input
		sequence = pad_sequences([sequence], maxlen=max_length)
		# predict next word
		yhat = model.predict([photo,sequence], verbose=0)
		# convert probability to integer
		yhat = argmax(yhat)
		# map integer to word
		word = word_for_id(yhat, tokenizer)
		# stop if we cannot map the word
		if word is None:
			break
		# append as input for generating the next word
		in_text += ' ' + word
		# stop if we predict the end of the sequence
		if word == 'endseq':
			break
	return in_text

def generate_all(model, tokenizer, img_files , photos, max_length):
    predicted = {}
    for fpath,name in img_files:
    	yhat = generate_desc(model, tokenizer, photos[name.split('.')[0]], max_length)
    	predicted[fpath] = yhat
    return predicted

def analyze_image(img_files):
	script_loc = os.path.join(os.path.dirname(os.path.realpath(sys.argv[0])) , "image_captioning")
	# load the tokenizer
	tokenizer_path = os.path.join(script_loc, 'tokenizer.pkl')
	tokenizer = load(open( tokenizer_path ,'rb'))
	# pre-define the max sequence length (from training)
	max_length = 34
	# load the model
	model_path = os.path.join(script_loc, 'model_inceptionV3_18.h5')
	model = load_model(model_path)
	# load and prepare the photograph
	features = extract_features(img_files)
	# names=[]
	# for f,n in img_files:
	# 	names.append(n.split('.')[0])
	# generate description
	pridictions = generate_all(model, tokenizer,img_files , features , max_length)
	dic = {}
	for path,des in pridictions.items():
		mini_dic = {}
		clean_des = clean_text(des[9:-7])
		for w in clean_des.split():
			mini_dic[w] = 1
		dic[path] = mini_dic

	#for key,value in pridictions.items():
	#    print(str(key) + ": " + str(clean_text(value[9:-7])))
	return dic