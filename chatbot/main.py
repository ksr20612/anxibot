import pymysql

if __name__ == '__main__':

    db = None;
    try :
        db = pymysql.connect(
            host='localhost',
            user='root',
            passwd='ehddus17',
            db='chatbot',
            charset='utf8'
        )
        print("DB 연결 성공")

    except Exception as e :
        print(e)

    finally:
        if db is not None :
            db.close()
            print("DB 닫기 성공")

