#ifndef CLOX_TABLE_H
#define CLOX_TABLE_H

#include "common.h"
#include "value.h"


typedef struct {
    objString* key;
    Value value;
} Entry;

typedef struct {
    int count;
    int capacity;
    Entry* entries;
} Table;

void initTable(Table* table);
void freeTable(Table* table);
bool tableSet(Table* table, objString* key, Value value);
bool tableDelete(Table* table, objString* key);
bool tableGet(Table* table, objString* key, Value* value);
void tableAddAll(Table* from, Table* to);
objString* tableFindString(Table* table, const char* chars, int lenght, uint32_t hash);

#endif