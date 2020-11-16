cmake_minimum_required(VERSION 2.8)

project(Test-TestServer)

option(TARS_MYSQL "option for mysql" ON)
option(TARS_SSL "option for ssl" ON)
option(TARS_GMSSL "option for ssl" OFF)
option(TARS_HTTP2 "option for http2" OFF)

set(TARS_TOOL_FLAG "--with-tars")

#set(CMAKE_BUILD_TYPE "Debug")
set(CMAKE_BUILD_TYPE "Release")

set(TARS_WEB_HOST "http://tars.test.com")

if(WIN32)
include (c:\\tars\\cpp\\makefile\\tars-tools.cmake)
else()
include (/usr/local/tars/cpp/makefile/tars-tools.cmake)
endif()

include_directories(/usr/local/tars/cpp/thirdparty/include)
link_directories(/usr/local/tars/cpp/thirdparty/lib)

add_subdirectory(ConfigServer)
add_subdirectory(NodeServer)
add_subdirectory(NotifyServer)
add_subdirectory(RegistryServer)

############################################################################

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${CMAKE_SOURCE_DIR}/cmake")

include(cmake/golang.cmake)

add_subdirectory(ImageServer)
add_subdirectory(AdminServer)
add_subdirectory(AgentServer)

install(DIRECTORY deploy DESTINATION ${TARS_PATH}
        PATTERN "*.sh"
        PERMISSIONS OWNER_EXECUTE OWNER_WRITE OWNER_READ GROUP_EXECUTE GROUP_READ)

install(DIRECTORY ${CMAKE_BINARY_DIR}/bin/ DESTINATION ${TARS_PATH}/deploy/k8s-template/Program
        PATTERN "tars*"
        PERMISSIONS OWNER_EXECUTE OWNER_WRITE OWNER_READ GROUP_EXECUTE GROUP_READ WORLD_READ WORLD_EXECUTE)
