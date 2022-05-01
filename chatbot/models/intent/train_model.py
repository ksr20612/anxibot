import pandas as pd
import tensorflow as tf
from tensorflow.keras import preprocessing
from tensorflow.keras.models import Model
from tensorflow.keras.layers import Input, Embedding, Dense, Dropout, Conv1D, GlobalMaxPool1D, concatenate

# 1 데이터 읽기
data = pd.read_csv("./total_train_data.csv", delimiter=',')
queries = data['query'].tolist()
intents = data['intent'].tolist()

# 2 전처리 프로세서 생성
from utils.preprocess import Preprocess
preprocessor = Preprocess(word2index_dic='../../train_tools/dict/chatbot_dict.bin',
                       userdic="../../utils/user_dic.tsv")

# 3 주요 개체 추출하여 시퀀스 생성
sequences = []
for sentence in queries:
    pos = preprocessor.pos(sentence)
    keywords = preprocessor.get_keywords(pos, without_tag=True)
    seq = preprocessor.get_wordidx_sequence(keywords)
    sequences.append(seq)

from config.GlobalParams import MAX_SEQ_LEN
padded_seqs = preprocessing.sequence.pad_sequences(sequences, maxlen=MAX_SEQ_LEN, padding='post')

ds = tf.data.Dataset.from_tensor_slices((padded_seqs, intents))
ds = ds.shuffle(len(queries), reshuffle_each_iterator=True)

train_size = int(len(padded_seqs)*0.7)
val_size = int(len(padded_seqs)*0.2)
test_size = int(len(padded_seqs)*0.1)

train_ds = ds.take(train_size).batch(20)  # 최적화 필요
val_ds = ds.skip(train_size).take(val_size).batch(20)
test_ds = ds.skip(train_size+val_size).take(test_size).batch(20)

# 4 하이퍼 파라미터 설정
dropout_prob = 0.5
EMB_SIZE = 128
EPOCH = 5
VOCAB_SIZE = len(preprocessor.word_index) + 1

# 5 CNN 모델 구현
# 5.1 임베딩(희소벡터->분산벡터)
input_layer = Input(shape=(MAX_SEQ_LEN,))
embedding_layer = Embedding(VOCAB_SIZE, EMB_SIZE, input_length=MAX_SEQ_LEN)(input_layer)
dropout_emb = Dropout(rate=dropout_prob)(embedding_layer)

# 5.2 합성곱 연산 + 풀링 연산 구현
conv1 = Conv1D(
    filters=128,
    kernel_size=3,
    padding='valid',
    activation=tf.nn.relu
)(dropout_emb)
conv2 = Conv1D(
    filters=128,
    kernel_size=4,
    padding='valid',
    activation=tf.nn.relu
)(dropout_emb)
conv3 = Conv1D(
    filters=128,
    kernel_size=5,
    padding='valid',
    activation=tf.nn.relu
)(dropout_emb)
pool1 = GlobalMaxPool1D()(conv1)
pool2 = GlobalMaxPool1D()(conv2)
pool3 = GlobalMaxPool1D()(conv3)
concat = concatenate([pool1, pool2, pool3])

# 5.3 완전연결계층 구현
hidden = Dense(128, activation=tf.nn.relu)(concat)
dropout_hidden = Dropout(rate=dropout_prob)(hidden)
logits = Dense(5, name='logits')(dropout_hidden)
predictions = Dense(5, activation=tf.nn.softmax)(logits)

# 5.4 모델 생성
model = Model(inputs=input_layer, outputs=predictions)
model.compile(optimizer='adam',
              loss='sparse_categorical_crossentropy',
              metrics=['accuracy'])

# 6 모델 학습
model.fit(train_ds, validation_data=val_ds, epochs=EPOCH, verbose=1)

# 7 모델 평가
loss, accuracy = model.evaluate(test_ds, verbose=1)
print('Accuracy : %f' % (accuracy*100))
print('Loss : %f' % (loss))

# 8 모델 저장
model.save('intent_model.h5')
