# Simple API package for Facebook Account Kit
## Install 
```
  go get github.com/madding/fb_account_kit
```

## Documentation for Facebook Account Here
  https://developers.facebook.com/docs/accountkit/graphapi
  
## How to use
Create client
```
  client, err := fb_account_kit.CreateClient(<authCode>,
			<FacebookAppId>,
			<FacebookAppSecret>)
```
where <authCode> authentification code, <FacebookAppId> id your facebook application can be taked from development page,
<FacebookAppSecret> application secret from AccountKit tab

After create client if all right you can take profile info:
```
  res, err := client.GetMe()
```
where res its map[string]interface{} like
```
    {  
       id: "1234512345123451",
       phone: {
         number: "+15551234567"
         country_prefix: "1",
         national_number: "5551234567"
       },
       application: {
         id: "5432154321543210"
       }
    }
```

## TODO
- add tests
- add other cases like logout, invalidate_all_tokens, delete
