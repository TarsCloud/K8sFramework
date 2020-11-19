include_directories(../RegistryServer)
include_directories(/usr/local/tars/cpp/include/servant)

gen_server(tars tarsnode)

add_dependencies(tarsnode tarsregistry)