#ifndef CLOX_COMPILER_H
#define CLOX_COMPILER_H

#include "object.h"
#include "chunk.h"

objFunction* compile(const char* source);
void markCompilerRoots();

#endif
