#include <stdlib.h>
#include "memory.h"
#include "table.h"
#include "vm.h"
#include "compiler.h"

#ifdef DEBUG_LOG_GC
#include <stdio.h>
#include "debug.h"
#endif

#define GC_HEAP_GROW_FACTOR 2

void* reallocate(void* pointer, size_t oldSize, size_t newSize){

  if(newSize > oldSize){
    #ifdef DEBUG_STRESS_GC
      collectGarbage();
    #endif

    if(vm.bytesAllocated > vm.nextGC){
      collectGarbage();
    }
  }

  if(newSize == 0){

    free(pointer);
    return NULL;
  }

  void* result = realloc(pointer, newSize);
  if(result == NULL) exit(1);
  return result;
}

void freeObject(obj* object){

  #ifdef DEBUG_LOG_GC
    printf("%p free type %d\n", (void*)object, object->type);
  #endif

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

    case OBJ_CLASS:{
      objClass* klass = (objClass*)object;
      freeTable(&klass->methods);
      FREE(objClass, object);
      break;
    }

    case OBJ_INSTANCE:{
      objInstance* instance = (objInstance*)object;
      freeTable(&instance->fields);
      FREE(objInstance, object);
      break;
    }

    case OBJ_CLOSURE:{
      objClosure* closure = (objClosure*)object;
      FREE_ARRAY(objUpValue*, closure->upValues, closure->upValueCount);
      FREE(objClosure, object);
      break;
    }

    case OBJ_UPVALUE:
      FREE(objUpValue, object);
      break;

    case OBJ_BOUND_METHOD:
      FREE(objBoundMethod, object);
      break;
  }
}

void markObject(obj* object){
  if(object == NULL) return;
  if(object->isMarked) return;

  #ifdef DEBUG_LOG_GC
    printf("%p mark", (void*)object);
    printValue(OBJ_VAL(object));
    printf("\n");
  #endif

  object->isMarked = true;

  if(vm.grayCapacity < vm.grayCount + 1){
    vm.grayCapacity = GROW_CAPACITY(vm.grayCapacity);
    vm.grayStack = (obj**)realloc(vm.grayStack, sizeof(obj*) + vm.grayCapacity);

    if(vm.grayStack == NULL) exit(1);
  }

  vm.grayStack[vm.grayCount++] = object;


}

void markValue(Value value){
  if(IS_OBJ(value)) markObject(AS_OBJ(value));
}

static void markArray(ValueArray* array){
  for(int i = 0; i < array->count; i++){
    markValue(array->values[i]);
  }
}

static void blackenObject(obj* object){
  #ifdef DEBUG_LOG_GC
    printf("%p blacken", (void*)object);
    printValue(OBJ_VAL(object));
    printf("\n");
  #endif
  switch(object->type){
    case OBJ_CLASS:{
      objClass* klass = (objClass*)object;
      markObject((obj*)klass->name);
      markTable(&klass->methods);
      break;
    }

    case OBJ_INSTANCE:{
      objInstance* instance = (objInstance*)object;
      markObject((obj*)instance->klass);
      markTable(&instance->fields);
      break;
    }

    case OBJ_CLOSURE:{
      objClosure* closure = (objClosure*)object;
      markObject((obj*)closure->function);
      for(int i = 0; i < closure->upValueCount; i++){
        markObject((obj*)closure->upValues[i]);
      }
      break;
    }
    case OBJ_FUNCTION:{
      objFunction* function = (objFunction*)object;
      markObject((obj*)function->name);
      markArray(&function->chunk.constants);
    }
    case OBJ_UPVALUE:
      markValue(((objUpValue*)object)->closed);
      break;
    case OBJ_NATIVE:
    case OBJ_STRING:
    break;

    case OBJ_BOUND_METHOD:{
      objBoundMethod* bound = (objBoundMethod*)object;
      markValue(bound->receiver);
      markObject((obj*)bound->method);
    }
  }
}

static void markRoots(){
  for(Value* slot = vm.stack; slot < vm.stackTop; slot++){
    markValue(*slot);
  }

  for(int i = 0; i < vm.frameCount; i++){
    markObject((obj*)vm.frames[i].closure);
  }

  for(objUpValue* upValue = vm.openUpValues; upValue != NULL; upValue = upValue->next){
    markObject((obj*)upValue);
  }

  markTable(&vm.globals);
  markCompilerRoots();
  markObject((obj*)vm.initString);
}

static void traceReference(){
  while(vm.grayCount > 0){
    obj* object = vm.grayStack[--vm.grayCount];
    blackenObject(object);
  }
}

static void sweep(){
  obj* previous = NULL;
  obj* object = vm.objects;

  while(object != NULL){
    if(object->isMarked){
      object->isMarked = false;
      previous = object;
      object = object->next;
    }else{
      obj* unreached = object;
      object = object->next;
      if(previous != NULL){
        previous->next = object;
      }else{
        vm.objects = object;
      }

      freeObject(unreached);
    }
  }
}

static void traceReferences(){
  while(vm.grayCount > 0){
    obj* object = vm.grayStack[--vm.grayCount];
    blackenObject(object);
  }
}

void collectGarbage(){
  #ifdef DEBUG_LOG_GC
    printf("-- gc begin\n");
    size_t before = vm.bytesAllocated;
  #endif

  markRoots();
  traceReferences();
  tableRemoveWhite(&vm.strings);
  sweep();

  vm.nextGC = vm.bytesAllocated * GC_HEAP_GROW_FACTOR;

  #ifdef DEBUG_LOG_GC
    printf("-- gc end\n");
    printf("  collected %zu bytes (from %zu to %zu) next at %zu\n", before-vm.bytesAllocated, before, vm.bytesAllocated, vm.nextGC);
  #endif
}

void freeObjects(){
  obj* object = vm.objects;
  while(object != NULL){
    obj* next = object->next;
    freeObject(object);
    object = next;
  }

  free(vm.grayStack);
}

