cmake_minimum_required(VERSION 2.8)

project(Test-TestServer)

option(TAF_MYSQL "option for mysql" ON)
option(TAF_SSL "option for ssl" ON)
option(TAF_GMSSL "option for ssl" OFF)
option(TAF_HTTP2 "option for http2" OFF)

set(JCE2CPP_FLAG "--with-taf")

#set(CMAKE_BUILD_TYPE "Debug")
set(CMAKE_BUILD_TYPE "Release")

set(TAF_WEB_HOST "http://taf.test.com")

if(WIN32)
include (c:\\taf\\cpp\\makefile\\taf-tools.cmake)
else()
include (/usr/local/taf/cpp/makefile/taf-tools.cmake)
endif()

include_directories(/usr/local/taf/cpp/thirdparty/include)
link_directories(/usr/local/taf/cpp/thirdparty/lib)

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
add_subdirectory(ControlServer)

install(DIRECTORY ../builder/deploy DESTINATION ${TAF_PATH}
        PATTERN "*.sh"
        PERMISSIONS OWNER_EXECUTE OWNER_WRITE OWNER_READ GROUP_EXECUTE GROUP_READ)

install(DIRECTORY ${CMAKE_BINARY_DIR}/bin/ DESTINATION ${TAF_PATH}/deploy/k8s-template/Program
        PATTERN "taf*"
        PERMISSIONS OWNER_EXECUTE OWNER_WRITE OWNER_READ GROUP_EXECUTE GROUP_READ WORLD_READ WORLD_EXECUTE)
