import pymysql
import openpyxl

from config.DatabaseConfig import *  # import all variables

def clear_train_data(db) :
    sql = '''
        delete from Answer
    '''
    with db.cursor() as cursor:
        cursor.execute(sql)
    sql = '''
        alter table Answer auto_increment=1
    '''
    with db.cursor() as cursor :
        cursor.execute(sql)

def insert_data(db, xls_row) :
    ans_intent, ans_entity, ans_query, ans_content = xls_row

    sql = '''
        insert Answer(ans_intent, ans_entity, ans_query, ans_content)
        values(
            '%d', '%d', '%s', '%s'
        )
    ''' % (ans_intent.value, ans_entity.value, ans_query.value, ans_content.value)

    sql = sql.replace("'None'", "null")
    with db.cursor() as cursor :
        cursor.execute(sql)
        print('{} 저장'.format(ans_query.value))
        db.commit()

train_file = './train_data.xlsx'
db = None
try:
    db = pymysql.connect(
        host=DB_HOST,
        user=DB_USER,
        passwd=DB_PASSWORD,
        db=DB_NAME,
        charset='utf8'
    )
    clear_train_data(db)

    wb = openpyxl.load_workbook(train_file)
    sheet = wb['Sheet1']
    for row in sheet.iter_rows(min_row=2):
        insert_data(db, row)

    wb.close()

except Exception as e:
    print(e)

finally:
    if db is not None:
        db.Close()

