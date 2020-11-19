
include_directories(../ThirdParty/asio/include)
include_directories(../ThirdParty/rapidjson/include)
include_directories(/usr/local/tars/cpp/include/servant)

gen_server(tars tarsregistry)

target_link_libraries(tarsregistry ssl crypto) 