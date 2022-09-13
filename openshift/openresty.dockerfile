FROM openresty/openresty:alpine-fat

RUN /usr/local/openresty/luajit/bin/luarocks install lua-resty-jwt && /usr/local/openresty/luajit/bin/luarocks install lua-resty-session && /usr/local/openresty/luajit/bin/luarocks install lua-resty-jwt && /usr/local/openresty/luajit/bin/luarocks install lua-resty-http && /usr/local/openresty/luajit/bin/luarocks install lua-resty-openidc
