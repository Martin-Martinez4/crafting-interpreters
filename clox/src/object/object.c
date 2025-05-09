#include <stdint.h>
#include <stdio.h>
#include <string.h>

#include "memory.h"
#include "object.h"
#include "table.h"
#include "value.h"
#include "vm.h"

#define ALLOCATE_OBJ(type, objectType) \
  (type*)allocateObject(sizeof(type), objectType)

static obj* allocateObject(size_t size, objType type){
  obj* object = (obj*)reallocate(NULL, 0, size);
  object->type = type;

  object->next = vm.objects;
  vm.objects = object;
  return object;
}

static objString* allocateString(char* chars, int length, uint32_t hash){
  objString* string = ALLOCATE_OBJ(objString, OBJ_STRING);
  string->length = length;
  string->chars = chars;
  string->hash = hash;
  push(OBJ_VAL(string));
  tableSet(&vm.strings, string, NIL_VAL);
  pop();
  return string;
}

static uint32_t hashString(const char* key, int length) {
  uint32_t hash = 2166136261u;
  for (int i = 0; i < length; i++) {
    hash ^= (uint8_t)key[i];
    hash *= 16777619;
  }
  return hash;
}

static void printFunction(objFunction* function){
  if(function->name == NULL){
    printf("<script>");
    return;
  }
  printf("<fn %s>", function->name->chars);
}

objClosure* newClosure(objFunction* function){
  objUpValue** upvalues = ALLOCATE(objUpValue*, function->upValueCount);
  for(int i = 0; i < function->upValueCount; i++){
    upvalues[i] = NULL;
  }
  objClosure* closure = ALLOCATE_OBJ(objClosure, OBJ_CLOSURE);
  closure->function = function;
  closure->upValues = upvalues;
  closure->upValueCount = function->upValueCount;
  return closure;
}

objFunction* newFunction(){
  objFunction* function = ALLOCATE_OBJ(objFunction, OBJ_FUNCTION);
  function->arity = 0;
  function->upValueCount = 0;
  function->name = NULL;
  initChunk(&function->chunk);
  return function;
}

objNative* newNative(NativeFn function){
  objNative* native = ALLOCATE_OBJ(objNative, OBJ_NATIVE);
  native->function = function;
  return native;
}

objString* takeString(char* chars, int length){
  uint32_t hash = hashString(chars, length);
  
  objString* interned = tableFindString(&vm.strings, chars, length, hash);
  if(interned != NULL){
    FREE_ARRAY(char, chars, length +1);
    return interned;
  }

  return allocateString(chars, length, hash);
}

objString* copyString(const char* chars, int length){
  uint32_t hash = hashString(chars, length);

  objString* interned = tableFindString(&vm.strings, chars, length, hash);
  if(interned != NULL) return interned;
  
  char* heapChars = ALLOCATE(char, length+1);
  memcpy(heapChars, chars, length);
  heapChars[length] = '\0';
  return allocateString(heapChars, length, hash);
}

objUpValue* newUpValue(Value* slot){
  objUpValue* upValue = ALLOCATE_OBJ(objUpValue, OBJ_UPVALUE);
  upValue->closed = NIL_VAL;
  upValue->location = slot;
  upValue->next = NULL;
  return upValue;
}

void printObject(Value value){
  switch(OBJ_TYPE(value)){
    case OBJ_STRING:
      printf("%s", AS_CSTRING(value));
      break;

    case OBJ_FUNCTION:
      printFunction(AS_FUNCTION(value));
      break;

    case OBJ_NATIVE:
      printf("<native fn>");
      break;

    case OBJ_CLOSURE:
      printFunction(AS_CLOSURE(value)->function);
      break;

    case OBJ_UPVALUE:
      printf("upvalue");
      break;
  }
}



