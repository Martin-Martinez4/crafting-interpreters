cmake_minimum_required(VERSION 3.10)
project(clox C)

# Generate compile_commands.json
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

set(CMAKE_BUILD_TYPE Debug)



set(sources ${sources}
  ${CMAKE_CURRENT_SOURCE_DIR}/src/memory/memory.c
  ${CMAKE_CURRENT_SOURCE_DIR}/src/value/value.c
  ${CMAKE_CURRENT_SOURCE_DIR}/src/object/object.c
  ${CMAKE_CURRENT_SOURCE_DIR}/src/table/table.c
  ${CMAKE_CURRENT_SOURCE_DIR}/src/chunk/chunk.c
  ${CMAKE_CURRENT_SOURCE_DIR}/src/scanner/scanner.c
  ${CMAKE_CURRENT_SOURCE_DIR}/src/compiler/compiler.c
  ${CMAKE_CURRENT_SOURCE_DIR}/src/debug/debug.c
  ${CMAKE_CURRENT_SOURCE_DIR}/src/vm/vm.c
)

# breaks linking in C?
# if(WIN32)
#     # Create your game executable target as usual
#     add_executable(
#   ${CMAKE_PROJECT_NAME}
#   WIN32
#   ${CMAKE_SOURCE_DIR}/src/main.c
#   ${sources}

# )
# else()
add_executable(
    ${CMAKE_PROJECT_NAME}
    ${CMAKE_SOURCE_DIR}/src/main.c
    ${sources} 
  )
# endif()

include_directories(
  ${CMAKE_PROJECT_NAME}
  PUBLIC
  ${CMAKE_CURRENT_SOURCE_DIR}/src/
  ${CMAKE_CURRENT_SOURCE_DIR}/src/common/
  ${CMAKE_CURRENT_SOURCE_DIR}/src/memory/
  ${CMAKE_CURRENT_SOURCE_DIR}/src/value/
  ${CMAKE_CURRENT_SOURCE_DIR}/src/object/
  ${CMAKE_CURRENT_SOURCE_DIR}/src/table/
  ${CMAKE_CURRENT_SOURCE_DIR}/src/chunk/
  ${CMAKE_CURRENT_SOURCE_DIR}/src/scanner/
  ${CMAKE_CURRENT_SOURCE_DIR}/src/compiler/
  ${CMAKE_CURRENT_SOURCE_DIR}/src/debug/
  ${CMAKE_CURRENT_SOURCE_DIR}/src/vm/
)
  


