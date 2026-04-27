bash
psql -U postgres -d mydb
Use code with caution.

1. The "Clean" Way (Recommended)
This removes all rows and resets the auto-incrementing ID counter back to 1.
sql
TRUNCATE TABLE users RESTART IDENTITY;
Use code with caution.

2. The "Simple" Way
This removes all rows but keeps the ID counter where it left off (e.g., if the last ID was 5, the next inserted row will be 6).
sql
DELETE FROM users;
Use code with caution.

3. The "Nuclear" Way
If you want to remove the entire table structure (the table itself will disappear), use:
sql
DROP TABLE users;
Use code with caution.

Important Tips:
Don't forget the semicolon (;) at the end of the command, or PostgreSQL won't execute it.
If you have other tables linked to users (Foreign Keys), you might need to add CASCADE:
sql
TRUNCATE TABLE users RESTART IDENTITY CASCADE;
Use code with caution.