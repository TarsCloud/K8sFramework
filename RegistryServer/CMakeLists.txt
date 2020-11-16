
include_directories(../thirdparty)
include_directories(../thirdparty/asio/include)

include_directories(/usr/local/tars/cpp/include/servant)

gen_server(tars tarsregistry)

target_link_libraries(tarsregistry ${LIB_MYSQL}) 