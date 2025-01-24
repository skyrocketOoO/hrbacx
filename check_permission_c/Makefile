MODULES = check_permission_c
EXTENSION = check_permission_c
DATA = sql/check_permission_c--1.0.sql

PG_CONFIG = pg_config
CFLAGS = $(shell $(PG_CONFIG) --cflags)
LDFLAGS = $(shell $(PG_CONFIG) --ldflags) -lm
SOURCES = check_permission_c.c

REGRESS = check_permission_c_test

include $(shell $(PG_CONFIG) --pgxs)
