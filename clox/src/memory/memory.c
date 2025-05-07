#include <stdlib.h>
#include "memory.h"
#include "value.h"
#include "vm.h"

void* reallocate(void* pointer, size_t oldSize, size_t newSize){
  if(newSize == 0){
    free(pointer);
    return NULL;
  }

  void* result = realloc(pointer, newSize);
  if(result == NULL) exit(1);
  return result;
}

void freeObject(obj* object){
  switch(object->type){
    case OBJ_STRING:
      objString* string = (objString*)object;
      FREE_ARRAY(char, string->chars, string->length+1);
      FREE(objString, object);
      break;

    case OBJ_FUNCTION:
      objFunction* function = (objFunction*)object;
      freeChunk(&function->chunk);
      FREE(objFunction, object);
      break;
  }
}

void freeObjects(){
  obj* object = vm.objects;
  while(object != NULL){
    obj* next = object->next;
    freeObject(object);
    object = next;
  }
}

