#!/bin/bash

mysql -h127.0.0.1 -P3306 -uroot -p'minerx(#)666' << EOF
grant all on *.* TO 'minerx'@'%' identified by "minerx(#)666";
flush privileges;
EOF
