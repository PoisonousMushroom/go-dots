# acct_v5
# ------------------------------------------------------------

DROP TABLE IF EXISTS `acct_v5`;

CREATE TABLE acct_v5 (
  `agent_id` INT(4) UNSIGNED NOT NULL,
  `class_id` CHAR(16) NOT NULL,
  `mac_src` CHAR(17) NOT NULL,
  `mac_dst` CHAR(17) NOT NULL,
  `vlan` INT(2) UNSIGNED NOT NULL,
  `ip_src` CHAR(15) NOT NULL,
  `ip_dst` CHAR(15) NOT NULL,
  `src_port` INT(2) UNSIGNED NOT NULL,
  `dst_port` INT(2) UNSIGNED NOT NULL,
  `ip_proto` CHAR(6) NOT NULL,
  `tos` INT(4) UNSIGNED NOT NULL,
  `packets` INT UNSIGNED NOT NULL,
  `bytes` BIGINT UNSIGNED NOT NULL,
  `flows` INT UNSIGNED NOT NULL,
  `stamp_inserted` DATETIME NOT NULL,
  `stamp_updated` DATETIME,
  PRIMARY KEY (`agent_id`, `class_id`, `mac_src`, `mac_dst`, `vlan`, `ip_src`, `ip_dst`, `src_port`, `dst_port`, `ip_proto`, `tos`, `stamp_inserted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
