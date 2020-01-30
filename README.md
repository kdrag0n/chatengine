# ChatEngine

ChatEngine was a conversational natural language chatbot written in Go. It was in development from July 9, 2017 to October 1, 2017.

**THIS CODE IS AVAILABLE SOLELY FOR REFERENCE PURPOSES.** It has not been modified from its original state. Use at your own risk; the quality of this code leaves a lot to be desired.

The original README can be found below for reference.

<details>
<summary>Original README</summary>

> MySQL connection string example:
> ```
> username:password@protocol(address)/dbname?param=value
> ```
> 
> `protocol` = the protocol used, like `unix` or `tcp`.
> 
> Redis settings:
> ```
>     "redis_protocol": "[tcp or unix]",
>     "redis_server": "[host:port or socket path]",
>     "redis_password": "s00per-sekr3t_p455w0rd",
>     "redis_database": 2
> ```
> 
> Example using TCP on localhost:
> ```
>     "redis_protocol": "tcp",
>     "redis_server": "127.0.0.1:6379",
>     "redis_database": 1,
>     "redis_password": "n0h-b0dyG3t51n"
> ```
> 
> Example using a Unix domain socket:
> ```
>     "redis_protocol": "unix",
>     "redis_server": "/run/redis/redis.sock"
> ```
> 
> Database is a number, 0 - max. Server default max is 16.
> 
> # Notes
> 
> `keys.json` is now `config.json`.

</details>
