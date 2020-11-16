
#pragma once

#include <ostream>
#include <string>

struct SqlStrWrapper {
public:
    explicit SqlStrWrapper(const char *_p) : p(_p) {}

    const char *const p;
};

template<class T>
struct SqlTWrapper {
public:
    explicit SqlTWrapper(const T &_t) : t(_t) {};
    const T &t;
};

template<std::size_t N>
static SqlStrWrapper SqlTr(const char _p[N]) {
    SqlStrWrapper _w(_p);
    return _w;
};

static SqlStrWrapper SqlTr(const char *p) {
    SqlStrWrapper _w(p);
    return _w;
};

static SqlStrWrapper SqlTr(const std::string &p) {
    SqlStrWrapper _w(p.c_str());
    return _w;
};

template<typename T>
typename std::enable_if<std::is_arithmetic<T>::value, SqlTWrapper<T>>::type
static SqlTr(const T &t) {
    SqlTWrapper<T> _w(t);
    return _w;
}

template<typename T>
typename std::enable_if<std::is_enum<T>::value, SqlTWrapper<T>>::type
static SqlTr(const T &t) {
    SqlTWrapper<T> _w(t);
    return _w;
}

std::ostream &operator<<(std::ostream &os, const SqlStrWrapper &t) {
    os << "'" << t.p << "'";
    return os;
}

template<class T>
std::ostream &operator<<(std::ostream &os, const SqlTWrapper<T> &t) {
    os << t.t;
    return os;
}


