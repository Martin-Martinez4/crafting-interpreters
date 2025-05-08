#include <stdlib.h>
#include "memory.h"
#include "object.h"
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
    case OBJ_STRING:{

      objString* string = (objString*)object;
      FREE_ARRAY(char, string->chars, string->length+1);
      FREE(objString, object);
      break;
    }

    case OBJ_FUNCTION:{

      objFunction* function = (objFunction*)object;
      freeChunk(&function->chunk);
      FREE(objFunction, object);
      break;
    }

    case OBJ_NATIVE:{
      FREE(objNative, object);
      break;
    }

    case OBJ_CLOSURE:{
      objClosure* closure = (objClosure*)object;
      FREE_ARRAY(objUpValue*, closure->upValues,
                 closure->upValueCount);
      FREE(objClosure, object);
      break;
    }

    case OBJ_UPVALUE:
      FREE(objUpValue, object);
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

