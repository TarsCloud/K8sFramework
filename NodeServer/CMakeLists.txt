gen_server(tars tarsnode)

target_link_libraries(tarsnode ${LIB_MYSQL}) 

add_dependencies(tarsnode tarsregistry)