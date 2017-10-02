---
title: API Reference

language_tabs: # must be one of https://git.io/vQNgJ
  - shell
  - http
  - python
  - javascript
  - java

toc_footers:
  - <a href='../signup'>Sign Up for an API Key</a>
  - <a href='../'>Back to Landing Page</a>

includes:
  - errors

search: true
---

# Introduction
Welcome to ChatEngine!
We provide a quality API for conversational chatbots.
Simply make an request, and get a response back for your input!
It works just like Cleverbot, but faster.

We have example API code for Shell, Python, JavaScript, and Java.
Check them out by switching programming languages with the tabs on the side.

All POST endpoints accept data in JSON format (MIME type: `application/json`).
Support for form body (MIME type: `application/x-www-form-urlencoded`) data may be added soon.

Libraries/clients used:
 
  - Shell: [curl](https://curl.haxx.se/)
  - HTTP: [HTTP](https://www.w3.org/Protocols/rfc2616/rfc2616.html)...
  - Python: [requests](http://docs.python-requests.org/en/master/)
  - JavaScript: [request](https://github.com/request/request)
  - Java: [OkHttp](https://square.github.io/okhttp/)

<aside class="notice">
Make sure your data is sent encoded in UTF-8!
</aside>

# Authentication
> To authorize, use this code:

```shell
curl "api_endpoint"
  -H "Authorization: super_secret_key"
```

```http
GET api_endpoint
Authorization: super_secret_key
```

```python
import requests

requests.get('api_endpoint', headers={'Authorization': 'super_secret_key'})
```

```javascript
const request = require('request');

request.get({
    'url': 'api_endpoint',
    'headers': {
        'Authorization': 'super_secret_key'
    }
}, (error, response, body) => {

});
```

```java
import okhttp3.OkHttpClient;
import okhttp3.Response;
import okhttp3.Request;

OkHttpClient client = new OkHttpClient();
Response response = client.newCall(new Request.Builder()
        .get()
        .url("api_endpoint")
        .header("Authorization", "super_secret_key")
        .build()).execute();
response.body().close();
```

> Be sure to replace `super_secret_key` with your API key.

ChatEngine uses API keys to authenticate and identify clients.
With every request, you must include a header that looks like this:

`Authorization: super_secret_key`

<aside class="notice">
Make sure you replace <code>super_secret_key</code> with your own API key.
</aside>

# Conversation

## Ask
```shell
curl "https://chatengine.xyz/api/ask"
  -XPOST
  -H "Content-Type: application/json"
  -H "Authorization: super_secret_key"
  -d '{"session": "session_id-101", "query": "Hi there!"}'
```

```http
GET /api/ask
Content-Type: application/json
Authorization: super_secret_key

{
    "session": "session_id-101",
    "query": "Hi there!"
}
```

```python
import json
import requests

data = {
    'session': 'session_id-101',
    'query': 'Hi there!'
}

response = requests.post('https://chatengine.xyz/api/ask',
                        headers={'Authorization': 'super_secret_key'}, data=json.dumps(data))

response_data = json.loads(response.text)
```

```javascript
const request = require('request');

request.post({
    'url': 'https://chatengine.xyz/api/ask',
    'headers': {
        'Authorization': 'super_secret_key'
    }
}, (error, response, body) => {
    let responseData = JSON.parse(response);
});
```

```java
import okhttp3.OkHttpClient;
import okhttp3.Response;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.MediaType;
import org.json.JSONObject;

private static final MediaType JSON_MEDIA_TYPE = MediaType.parse("application/json; charset=utf-8");

OkHttpClient client = new OkHttpClient();
String response = client.newCall(new Request.Builder()
        .post(RequestBody.create(JSON_MEDIA_TYPE, new JSONObject()
                        .put("session", "session_id-101")
                        .put("query", "Hi there!")
                        .toString()))
        .url("https://chatengine.xyz/api/ask")
        .header("Authorization", "super_secret_key")
        .build()).execute().body().string();

JSONObject responseData = new JSONObject(response);
```

> The above code returns JSON structured like this:

```json
{
    "success": true,
    "session": "session_id-101",
    "response": "Heya!",
    "confidence": 0.451
}
```

> `session` will only be present if you didn't specify a session.
In that case, please save the session ID for future requests (in the same conversation).

> If an error occurred, the JSON will be structured like this:

```json
{
    "success": false,
    "error": "You must specify some text!",
    "response": "An error occurred: You must specify some text!",
    "confidence": 0
}
```

This endpoint returns an answer for your query, as a step in a conversation.

<aside class="notice">
Replace <code>session_id-101</code> with your own session ID.
Adding numbers to make it harder to guess is suggested.
Sessions are IP restricted to the IP that started them.
</aside>

### HTTP Request
`POST /api/ask`

<aside class="warning">
The old <code>/ask</code> endpoint has been removed. Please use <code>/api/ask</code>.
</aside>

### Request Data
Parameter | Default | Type   | Description
--------- | ------- | ------ | -----------
session   | random  | string | Used to identify the conversation, and provide context. A random one will be generated and returned if you don't specify one.
query     | N/A     | string | This parameter is required, as it is the whole intent of the endpoint.

<aside class="success">
Remember to send the data in JSON!
</aside>
