#include <stdlib.h>
#include <stdio.h>
#include "chunk.h"
#include "memory.h"
#include "vm.h"

void initChunk(Chunk* chunk) {
  chunk->count = 0; 
  chunk->capacity = 0;
  chunk->code = NULL;
  chunk->lines = NULL;
  chunk->lineCount = 0;
  chunk->lineCapacity = 0;

  initValueArray(&chunk->constants);
}

void freeChunk(Chunk* chunk){
  FREE_ARRAY(uint8_t, chunk->code, chunk->capacity);
  FREE_ARRAY(int, chunk->lines, chunk->lineCapacity);

  freeValueArray(&chunk->constants);
  
  initChunk(chunk);
}

static void setLine(Chunk* chunk, int line){
  if(chunk->lineCapacity < chunk->lineCount + 1){
    int oldCap = chunk->lineCapacity;
    chunk->lineCapacity = GROW_CAPACITY(oldCap);
    chunk->lines = GROW_ARRAY(int, chunk->lines, oldCap, chunk->lineCapacity);
  }
  if(chunk->lines[chunk->lineCount-2] == line){
    chunk->lines[chunk->lineCount-1]++;
  }
  else{
    
    chunk->lineCount += 2;
    chunk->lines[chunk->lineCount-2] = line;
    chunk->lines[chunk->lineCount-1] = chunk->lines[chunk->lineCount-3] > 0 ? chunk->lines[chunk->lineCount-3] + 1 : 1;
  }
}

int getLine(Chunk* chunk, int offset){
  // printf("offset: %d\n", offset);
  int index = 1;
  int l = chunk->lines[index];

  while(offset >= l){
    index += 2;
    l = chunk->lines[index];
  }

  return chunk->lines[index-1];

}

void writeChunk(Chunk* chunk, uint8_t byte, int line){
    if(chunk->capacity < chunk->count + 1){
        int oldCap = chunk->capacity;
        chunk->capacity = GROW_CAPACITY(oldCap);
        chunk->code = GROW_ARRAY(uint8_t, chunk->code, oldCap, chunk->capacity);
        // chunk->lines = GROW_ARRAY(int, chunk->lines, oldCap, chunk->capacity);
    }

    chunk->code[chunk->count] = byte;
    // chunk->lines[chunk->count] = line;
    chunk->count++;

    setLine(chunk, line);

}

int addConstant(Chunk* chunk, Value value){
  // temp push onto stack so GC can see it
  push(value);
  writeValueArray(&chunk->constants, value);
  pop();
  return chunk->constants.count - 1;
}



