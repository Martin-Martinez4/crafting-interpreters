#ifndef CLOX_OBJECT_H
#define CLOX_OBJECT_H

#include "common.h"
#include "value.h"

#define OBJ_TYPE(value) (AS_OBJ(value)->type)
#define IS_STRING(value) isObjType(value, OBJ_STRING)

#define AS_STRING(value) ((objString*)AS_OBJ(value))
#define AS_CSTRING(value) (((objString*)AS_OBJ(value))->chars)

typedef enum {
  OBJ_STRING,
} objType;

struct obj {
  objType type;
  struct obj* next;
};

// struct composition
struct objString {
  obj obj;
  int length;
  char* chars;
  // cache the hash because it is expensive to calculate
  uint32_t hash;
};

objString* takeString(char* chars, int length);

objString* copyString(const char* chars, int length);

void printObject(Value value);

static inline bool isObjType(Value value, objType type){
  return IS_OBJ(value) && AS_OBJ(value)->type == type;
}

#endif