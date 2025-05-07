#ifndef CLOX_VM_H
#define CLOX_VM_H

#include "chunk.h"
#include "value.h"
#include "compiler.h"
#include "table.h"


#define STACK_MAX 256

typedef struct {
    Chunk* chunk;
    uint8_t* ip;
    Value stack[STACK_MAX];
    Table globals;
    Table strings;
    Value* stackTop;
    obj* objects;
} VM;

typedef enum {
    INTERPRET_OK,
    INTERPRET_COMPILE_ERROR,
    INTERPRET_RUNTIME_ERROR,
} InterpreterResult;

extern VM vm;

void initVM();
void freeVM();
InterpreterResult interpret(const char* source);

void push(Value value);
Value pop();


#endif
