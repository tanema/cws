Chrome Web Store cli
=======

This is a small tool to help automate interactions with the chrome web store API.
This can be used in CI or locally as well. It will archive(zip) your extension code,
updating the version and removing any development key in the manifest in the archive.

# Usage
To see full usage run `cws -h` or `cws [command] -h`. To upload and publish an
extension you can run:

```bash
cws deploy ./extension_src
```

# Config
`cws` uses a json config so that you can keep it in your repo (but not committed use `.gitignore`)
and constantly run cws commands in short form like `cws status` It also supports
environment variables as well which take precedence over the json config. This
is great for CI usage.

To see how to obtain these values, please look at the [Authentication Setup](#AuthenticationSetup)

### ENV Vars

| Environment Variable | Value
|----------------------|---------
|`CWS_EXTENSION_ID`    | Chrome Web Store Extension ID
|`CWS_CLIENT_ID`       | Google OAuth Client ID
|`CWS_CLIENT_SECRET`   | Google OAuth Client Secret
|`CWS_REFRESH_TOKEN`   | Google OAuth Refresh Token

### JSON config example

```json
{
  "extension_id": "your-extension-id",
  "client_id": "your-client-id",
  "client_secret": "your-client-secret",
  "refresh_token": "your-refresh-token"
}
```

# Screen Shot
`cws` has helpful error output and actions to help complete a process.

![Screen Shot 2022-10-11 at 10 54 21 AM](https://user-images.githubusercontent.com/463193/195125995-5975f5b6-3572-43f8-aa8d-4c1f3810b534.png)

# Authentication Setup
To have authentication to the API that we need, we need a already authorized OAuth
client, with offline access so that we can get a single refresh token that will not
expire for a while. The following is a step by step process to do this without the
need to write a script to authenticate.

- Go to your [credentials page for your project](https://console.cloud.google.com/apis/credentials)
- Click *+ CREATE CREDENTIALS* at the top of the page and select *OAuth client ID* and create the app with the following settings.
  - Application Type: Web
  - Name: Your App Name
  - Authorized redirect URIs: http://localhost:3333 (cws will start a local server to wait for the response)
- Use the client id and secret to run `cws init [client-id] [client-secret]`
- Click on the outputted link and click Authorize APIs.
  - On the next screen choose the account (optional screen) and give the permissions to the app.
  - You may get a warning that the app is not verified, do not worry, it is referring to your oauth client, click advanced and then click proceeed.
- Once you close the tab, you should now have a `chrome_webstore.json` file. Fill in the
  extension_id of your extension.
