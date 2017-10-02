# Errors
The ChatEngine API uses the following error codes:

Error Code | Name | Meaning
---------- | ---- | -------
400 | Bad Request | Your request is ill-formatted or missing data.
401 | Unauthorized | Your API key is invalid.
403 | Forbidden | That endpoint isn't available to your key's tier.
404 | Not Found | That endpoint couldn't be found.
405 | Method Not Allowed | You tried to use an endpoint with the wrong method.
406 | Not Acceptable | You sent a format that's not JSON.
410 | Gone | That endpoint has been removed.
418 | I'm A Teapot | Tip me over and pour me out.
429 | Too Many Requests | You're requesting too many things, too fast.
500 | Internal Server Error | The server exploded processing your request.
503 | Service Unavailable | We're temporarily offline for maintenance. Try again later.
