#ifndef CLOX_OBJECT_H
#define CLOX_OBJECT_H

#include "common.h"
#include "value.h"
#include "chunk.h"

#define OBJ_TYPE(value) (AS_OBJ(value)->type)

#define IS_STRING(value) isObjType(value, OBJ_STRING)
#define IS_NATIVE(value) isObjType(value, OBJ_NATIVE)
#define IS_CLOSURE(value) isObjType(value, OBJ_CLOSURE)
#define IS_FUNCTION(value) isObjType(value, OBJ_FUNCTION)

#define AS_STRING(value) ((objString*)AS_OBJ(value))
#define AS_CSTRING(value) (((objString*)AS_OBJ(value))->chars)
#define AS_CLOSURE(value) ((objClosure*)AS_OBJ(value))
#define AS_FUNCTION(value) ((objFunction*)AS_OBJ(value)) 
#define AS_NATIVE(value) (((objNative*)AS_OBJ(value))->function)

typedef enum {
  OBJ_CLOSURE,
  OBJ_FUNCTION,
  OBJ_NATIVE,
  OBJ_STRING,
  OBJ_UPVALUE,
} objType;

struct obj {
  objType type;
  struct obj* next;
};


typedef struct {
  obj obj;
  int arity;
  int upValueCount;
  Chunk chunk;
  objString* name;
} objFunction;


typedef struct{
  obj obj;
  Value* location;
  Value closed;
  struct objUpValue* next;
} objUpValue;
typedef struct {
  obj obj;
  objFunction* function;
  objUpValue** upValues;
  int upValueCount;
} objClosure;

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


objClosure* newClosure(objFunction* function);
objFunction* newFunction();
objNative* newNative(NativeFn function);
objUpValue* newUpValue(Value* slot);

objString* takeString(char* chars, int length);

objString* copyString(const char* chars, int length);

void printObject(Value value);

static inline bool isObjType(Value value, objType type){
  return IS_OBJ(value) && AS_OBJ(value)->type == type;
}

#endif