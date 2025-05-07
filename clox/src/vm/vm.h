#ifndef CLOX_VM_H
#define CLOX_VM_H

#include "chunk.h"
#include "value.h"
#include "compiler.h"
#include "table.h"
#include "object.h"


#define FRAMES_MAX 64
#define STACK_MAX (FRAMES_MAX * UINT8_COUNT)

typedef struct {
    objFunction* function;
    uint8_t* ip;
    Value* slots;
} CallFrame;

typedef struct {
    Chunk* chunk;
    uint8_t* ip;
    CallFrame frames[FRAMES_MAX];
    int frameCount;
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
