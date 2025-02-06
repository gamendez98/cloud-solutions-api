CREATE USER cloud_solutions WITH PASSWORD 'cloud_solutions';
CREATE DATABASE cloud_solutions;
GRANT ALL PRIVILEGES ON DATABASE cloud_solutions TO cloud_solutions;

\c cloud_solutions

GRANT ALL ON SCHEMA public TO cloud_solutions;

CREATE EXTENSION vector;