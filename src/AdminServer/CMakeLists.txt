

cmake_minimum_required(VERSION 2.8)

project(tarsadmin Go) # select GO compile

add_go_executable(tarsadmin)

message(${CMAKE_CURRENT_SOURCE_DIR})

add_custom_target(swagger
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    COMMAND docker pull quay.io/goswagger/swagger
    COMMENT "docker install swagger")

add_custom_target(openapi
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    COMMAND sh create-openapi.sh
    COMMENT "create openapi")