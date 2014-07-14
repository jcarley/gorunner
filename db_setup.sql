CREATE DATABASE gorunner;
USE gorunner;

CREATE USER 'gorunner-admin'@'localhost' IDENTIFIED BY 'letmein123';

GRANT USAGE ON *.* TO 'gorunner-admin'@'localhost' IDENTIFIED BY 'letmein123'
WITH MAX_QUERIES_PER_HOUR 0
MAX_UPDATES_PER_HOUR 0
MAX_CONNECTIONS_PER_HOUR 0
MAX_USER_CONNECTIONS 0;

GRANT Create Routine, Insert, Lock Tables, References, Select, Drop, Delete, Index, Alter Routine, Create View, Create Temporary Tables, Show View, Trigger, Event, Create, Update, Execute, Grant Option, Alter ON `gorunner`.* TO `gorunner-admin`@`localhost`;
