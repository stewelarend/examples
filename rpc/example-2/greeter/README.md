# Example 2

Same as Example1 but we use a request structure in the operation instead of a handler function.
The request structure must implement rpc.IHandler with method Exec()
The request structure may also implement rpc.IValidator with method Validate() to validate a request before calling Exec()

Note that the body is parsed first and URL params are individually parsed afterwards to overwrite the body values as needed.
Body is only parsed on method POST.
URL is parse for POST and GET.

The server config is in ./conf/config.json
The server runs on HTTP port specified in that file and can be tested with:

$ curl -s -XPOST "http://localhost:8000/hello?name=Jan&ageu8=100" -d'{"name":"koos","age16":16}'
"Hi {Name:Jan Age:0 Age8:0 Age16:16 Age32:0 Age64:0 Ageu8:100 Ageu16:0 Ageu32:0 Ageu64:0}!"
