import pandas as pd
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
from tensorflow.keras.preprocessing.text import Tokenizer
from tensorflow.keras.preprocessing.sequence import pad_sequences
from tensorflow.keras.layers import SimpleRNN, Embedding, Dense
from tensorflow.keras.models import Sequential

# 데이터 가져오기
data = pd.read_csv('./yesno.txt', encoding='utf8')
data['v2'] = data['v2'].replace(['yes','no'],[1,0]) 
print(data.info())

# 데이터 확인
print('v1열의 유니크값 : ',data['v1'].nunique())
data.drop_duplicates(subset=['v1'], inplace=True) # 중복 제거

# 데이터 분리 
X_data = data['v1']
Y_data = data['v2']
X_train, X_test, Y_train, Y_test = train_test_split(X_data, Y_data, test_size=0.2, stratify=Y_data) ### stratify = 고르게 카테고리 분류할 레이블

# 토크나이징, 단어 임베딩
tokenizer = Tokenizer()
tokenizer.fit_on_texts(X_train)
X_train_encoded = tokenizer.texts_to_sequences(X_train) # 단어 임베딩 : [[1,2,3], [4,5,6,7], [8,9,10,11]...]
word_to_index = tokenizer.word_index # data 하의 모든 단어에 숫자 매기기
print(word_to_index)
vocab_size = len(word_to_index) + 1
max_len = max(len(l) for l in X_train_encoded) # X_train_encoded 의 가장 긴 문장 길이 뽑아오기
X_train_padded = pad_sequences(X_train_encoded, maxlen = max_len) # maxLen 만큼 배열 길게 만들고, 빈 칸 0으로 채우기
print("padded train data : ", X_train_padded[:5])
print("훈련 데이터 : ", X_train_padded.shape)

# RNN
model = Sequential()
model.add(Embedding(vocab_size, 16))
model.add(SimpleRNN(32))
model.add(Dense(1, activation='sigmoid'))

model.compile(optimizer='rmsprop', loss='binary_crossentropy', metrics=['acc'])
history = model.fit(X_train_padded, Y_train,  epochs=32, batch_size=64, validation_split=0.2) # 훈련데이터 / 검증데이터 분리

X_test_encoded = tokenizer.texts_to_sequences(X_test)
X_test_padded = pad_sequences(X_test_encoded, maxlen=max_len)
print("테스트 정확도 : %.4f" % (model.evaluate(X_test_padded, Y_test)[1]))
print(model.predict(X_test_padded))

# epoch 별 적합도 확인
epochs = range(1, len(history.history['acc']) + 1)
plt.plot(epochs, history.history['loss'])
plt.plot(epochs, history.history['val_loss'])
plt.title('model loss')
plt.ylabel('loss')
plt.xlabel('epoch')
plt.legend(['train', 'val'], loc='upper left')
plt.show()