#!/usr/bin/env python3
import pymysql
import bcrypt

# Database connection
connection = pymysql.connect(
    host='mysql',
    user='reports',
    password='reports',
    database='reports',
    charset='utf8mb4'
)

try:
    with connection.cursor() as cursor:
        # Generate bcrypt hash
        password = b'password123'
        hashed = bcrypt.hashpw(password, bcrypt.gensalt())
        password_hash = hashed.decode('utf-8')
        
        print(f"Generated hash: {password_hash}")
        
        # Update password
        sql = "UPDATE users SET password = %s WHERE email = %s"
        cursor.execute(sql, (password_hash, 'admin@test.ru'))
        connection.commit()
        
        print(f"Rows affected: {cursor.rowcount}")
        
        # Verify
        cursor.execute("SELECT password FROM users WHERE email = %s", ('admin@test.ru',))
        result = cursor.fetchone()
        
        if result:
            stored_hash = result[0]
            print(f"Stored hash: {stored_hash}")
            
            # Test verification
            if bcrypt.checkpw(password, stored_hash.encode('utf-8')):
                print("Password verification: SUCCESS")
            else:
                print("Password verification: FAILED")
        else:
            print("User not found")

finally:
    connection.close()

