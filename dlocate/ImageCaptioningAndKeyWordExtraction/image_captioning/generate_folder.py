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

from desc_cleaning import clean_text

import sys
import os

# extract features from each photo in the directory
def extract_features(directory):
	# load the model
	#model = VGG16()
	model = InceptionV3()

	# re-structure the model
	model.layers.pop()
	model = Model(inputs=model.inputs, outputs=model.layers[-1].output)
	# extract features from each photo
	features = dict()
	for name in os.listdir(directory):
		# load an image from file
		filename = directory + '/' + name
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

def generate_all(model, tokenizer, names , photos, max_length):
    predicted = {}
    for name in names:
        yhat = generate_desc(model, tokenizer, photos[name], max_length)
        predicted[name] = yhat
    return predicted
 
# load the tokenizer
tokenizer = load(open('tokenizer.pkl', 'rb'))
# pre-define the max sequence length (from training)
max_length = 34
# load the model
model = load_model('model_inceptionV3_18.h5')
# load and prepare the photograph
directory = str(sys.argv[1])
features = extract_features(directory)
names=[]

for path, subdirs, files in os.walk('Images'):
   for filename in files:
       f = str(filename.split('.')[0])
       names.append(f)

# generate description
pridictions = generate_all(model, tokenizer,names , features , max_length)
for key,value in pridictions.items():
    print(str(key) + ": " + str(clean_text(value[9:-7])))