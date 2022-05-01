from utils.preprocess import Preprocess
from tensorflow.keras import preprocessing
import pickle

def read_corpus_data(filename):
    with open(filename, 'r') as f:
        data = [line.split('\t') for line in f.read().splitlines()]
        data = data[1:]
    return data

# 1 말뭉치 읽어오기
corpus_data = read_corpus_data('./corpus.txt')

# 2 말뭉치에서 단어 가져오고, 품사 태깅하기
preprocessor = Preprocess()
dict = []
for c in corpus_data:
    pos = preprocessor.pos(c[1])
    for k in pos:
        dict.append(k[0])

# 3 단어 인덱스 데이터로 만들기
tokenizer = preprocessing.text.Tokenizer(oov_token='OOV')
tokenizer.fit_on_texts(dict)
word_index = tokenizer.word_index

# 4 단어 인덱스 데이터 -> 파일
f = open("chatbot_dic.bin", "wb")
try:
    pickle.dump(word_index, f)
except Exception as e:
    print(e)
finally:
    f.close()
