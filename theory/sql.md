# SQL

## SQL Delete vs Truncate vs DROP Table
```sql
TRUNCATE TABLE users;
DELETE FROM users;
```
- If we want to delete all rows from a table, it is more efficient to use the TRUNCATE statement instead of the DELETE statement.
- it is more efficient than the DELETE statement because it does not log the deleted rows and deletes all the rows in one step, rather than deleting one row at a time.
- drop table deletes the table itself

## SQL primary key and Composit Primary key

### Primary key
- they are UNIQUE and can not be NULL.
- you can have ONLY ONE primary key in a table
```SQL
CREATE TABLE departments (
    id INTEGER PRIMARY KEY,
    name TEXT
);
```
### Composit Primary key
- it is a primary key made of multiple columns
- It is used when a single column cannot uniquely identify each row in a table.
```SQL
CREATE TABLE employees (
    department_id INTEGER,
    employee_id INTEGER,
    name TEXT,
    PRIMARY KEY (department_id, employee_id)
);
```

## SQL serial vs identity vs sequence vs UUID 

- serial: This is a shorthand for creating an auto-incrementing integer column.
- identity: introduced in PostgreSQL 10, this is the SQL standard way to create auto-incrementing columns.
- sequence: You can define your own sequence, and set the start value and increment value.
- UUID: It is a 16 byte value that is randomly generated and is guaranteed to be unique across all instances of a database
```SQL
CREATE SEQUENCE example_id_seq START WITH 10 INCREMENT BY 10;

CREATE TABLE example (
    id_0 INTEGER PRIMARY KEY,
    id_1 SERIAL,
    id_2 INTEGER GENERATED ALWAYS AS IDENTITY,
    id_3 INTEGER DEFAULT nextval('example_id_seq')
    id_4 INTEGER DEFAULT uuid_generate_v4()
);
```

## SQL database INDEXES what are they:

imagine this 
```SQL
SELECT * FROM users WHERE email = 'bob@email.com';
```
- Without an index the database would need to check EVERY ROW one by one to check if it matches the WHERE condition
- With an index you keep in a separate place the index value and a pointer to the actual row, uses more memory but makes lookup faster


## SQL inner JOIN, LEFT JOIN, Right JOIN, FULL JOIN:
- INNER JOIN (the same as a normal JOIN) returns only matching rows from both tables.
- LEFT JOIN returns all rows from the left table, plus matched rows from the right.
- RIGHT JOIN returns all rows from the right table, plus matched rows from the left.
- Unmatched rows in LEFT/RIGHT JOIN get NULL values for missing columns.
- FULL JOIN: grabs all (LEFT and RIGHT join combination)

## What are SQL Common Table Expressions:
With PostgreSQL you can temporarily save the results of a query to a Common Table Expression (CTE). This is useful when you want to reuse the results of a query in multiple places in your query.

```SQL 
WITH high_value_orders AS (
    SELECT customer_id, SUM(price) as total_spend
    FROM orders
    GROUP BY customer_id
    HAVING SUM(price) > 1000
)
```
A CTE is not persisted in the database. It is only used for the duration of the query. A CTE should not have a semicolon at the end.
## What is soft delete:
Add a deleted_at, and on delete do not remove, just register the deleted_at
## what are SQL views
A SQL View is like a saved query or a virtual table.
It doesn't store actual data — it stores a SELECT statement, and when you query the view, it runs that SELECT.
- Simplify complex queries
- Improve readability and reuse
- Add security (hide certain columns or tables)
- Help with reporting