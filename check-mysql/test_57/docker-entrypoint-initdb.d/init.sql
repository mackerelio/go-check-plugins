CREATE USER 'repl'@`%` IDENTIFIED BY 'repl';
GRANT REPLICATION SLAVE ON *.* TO 'repl'@`%`;
GRANT SELECT ON *.* TO 'repl'@`%`;
