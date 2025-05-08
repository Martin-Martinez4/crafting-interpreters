#ifndef CLOX_OBJECT_H
#define CLOX_OBJECT_H

#include "common.h"
#include "value.h"
#include "chunk.h"

#define OBJ_TYPE(value) (AS_OBJ(value)->type)

#define IS_STRING(value) isObjType(value, OBJ_STRING)
#define IS_NATIVE(value) isObjType(value, OBJ_NATIVE)
#define IS_FUNCTION(value) isObjType(value, OBJ_FUNCTION)

#define AS_STRING(value) ((objString*)AS_OBJ(value))
#define AS_CSTRING(value) (((objString*)AS_OBJ(value))->chars)
#define AS_FUNCTION(value) ((objFunction*)AS_OBJ(value)) 
#define AS_NATIVE(value) (((objNative*)AS_OBJ(value))->function)

typedef enum {
  OBJ_FUNCTION,
  OBJ_NATIVE,
  OBJ_STRING,
} objType;

struct obj {
  objType type;
  struct obj* next;
};

typedef struct {
  obj obj;
  int arity;
  Chunk chunk;
  objString* name;
} objFunction;

typedef Value (*NativeFn)(int argCount, Value* args);

typedef struct {
  obj obj;
  NativeFn function;
} objNative;

// struct composition
struct objString {
  obj obj;
  int length;
  char* chars;
  // cache the hash because it is expensive to calculate
  uint32_t hash;
};

objFunction* newFunction();
objNative* newNative(NativeFn function);

objString* takeString(char* chars, int length);

objString* copyString(const char* chars, int length);

void printObject(Value value);

static inline bool isObjType(Value value, objType type){
  return IS_OBJ(value) && AS_OBJ(value)->type == type;
}

#endif