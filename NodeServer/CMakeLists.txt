include_directories(../RegistryServer)
include_directories(./)
include_directories(/usr/local/tars/cpp/include/servant)

gen_server(tars tarsnode)


target_link_libraries(tarsnode ${LIB_MYSQL}) 

add_dependencies(tarsnode tarsregistry)