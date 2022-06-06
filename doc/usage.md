# Usage
Hit a url which is protected by distil, get the distil javascript library, example shown;
```
<script type="text/javascript" src="/pvvhnzyazwpzgkhv.js" defer></script>
```
If you need help identifying this, please provide us a URL which we can check this out at. The name of the js file _tends_ to stay the same,
however it is not guaranteed to always be the same.

You must get the data for this javascript library as it may change per ip address / session, base64 this and
dump the information into the following json blob;
```
curl http://www1.ticketmaster.com/pvvhnzyazwpzgkhv.js | base64 | pbcopy
```

request.json;
```
{
  "JsSha1": "af2d0557c23ff2d8f40ccf4bec57e480704634e9",
  "JsUri": "https://www1.ticketmaster.com/pvvhnzyazwpzgkhv.js",
  "JsData": BASE64ED_DATA_FOR_LIB
}
```
If you would like to, you can truncate the "JsUri" to not include the full uri, however please ensure to keep the js filename present.

Hit the endpoint:
```
curl -X POST -H "Auth-Key:{AUTH_KEY_HERE}" https://${API_SERVER_ADDRESS}/api/v1/session -d @request.json | jq .
```

The JSON response contains the following data;
```
{
  "tasks": [
    {
      "uri": "/pvvhnzyazwpzgkhv.js?PID=14CDB9B4-DE01-3FAA-AFF5-65BC2F771745", // This is the URI to hit
      "method": "POST", // This is the method to use when hitting the URI
      "headers": [ // Add any of the headers specified in blob for that URI
        "X-Distil-Ajax:twzvbatvrxzavsfzbzeyurav"
      ],
      "data": "p=%7B%22appName%22%....", // Data (only for posts) to write in the request
      "interval": 0 // Interval of when to perform task (0 is immediately, before other requests)
    },
    {
      "uri": "/pvvhnzyazwpzgkhv.js", // This is the uri to hit
      "method": "HEAD", // Method when hitting URI
      "headers": [],
      "data": "",
      "interval": 270000 // Reoccuring interval in milliseconds to perform this task
    }
  ],
  // Data seperate from tasks is to be applied for all requests following the POST task, but not tasks with intervals
  "Headers": [
    "X-Distil-Ajax:twzvbatvrxzavsfzbzeyurav"
  ]
}
```

## TLDR

 - Hit the session endpoint and provide the distil lib
 - Perform the tasks received from the endpoint
 - One task will be a POST, this must be performed first, this creates a session
 - Another task will be a reoccuring HEAD post, perform this will keep the IP/session alive
 - Use the Headers sent back, along with preserving any cookies that are sent back with the original request
 - A session should be good for upwards of 10k requests, however it might be worth while creating new sessions every 1k