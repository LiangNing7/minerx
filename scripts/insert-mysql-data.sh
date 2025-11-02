#!/bin/bash

function casbinrules() {
  mysql -h127.0.0.1 -uminerx -p"miner(#)666" -D onex << EOF
INSERT INTO casbin_rule VALUES (id,'p', 'alice', 'data1', 'read', 'allow', '', '');
INSERT INTO casbin_rule VALUES (id,'p', 'bob', 'data2', 'write', 'deny', '', '');
INSERT INTO casbin_rule VALUES (id,'p', 'data2_admin', 'data2', 'read', 'allow', '', '');
INSERT INTO casbin_rule VALUES (id,'p', 'data2_admin', 'data2', 'write', 'allow', '', '');
INSERT INTO casbin_rule VALUES (id,'g', 'alice', 'data2_admin', 'deny', '', '', '');
EOF
}

$*
