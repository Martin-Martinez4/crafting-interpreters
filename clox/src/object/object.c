#include <stdint.h>
#include <stdio.h>
#include <string.h>

#include "memory.h"
#include "object.h"
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
  return string;
}

static uint32_t hashString(const char* key, int length){
  
}

objString* takeString(char* chars, int length){
  uint32_t hash = hashString(chars, length);
  return allocateString(chars, length, hash);
}

objString* copyString(const char* chars, int length){
  char* heapChars = ALLOCATE(char, length+1);
  memcpy(heapChars, chars, length);
  heapChars[length] = '\0';
  return allocateString(heapChars, length,);
}

void printObject(Value value){
  switch(OBJ_TYPE(value)){
    case OBJ_STRING:
      printf("%s", AS_CSTRING(value));
      break;
  }
}



