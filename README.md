Chrome Web Store cli
=======

This is a small tool to help automate interactions with the chrome web store API.
This can be used in CI or locally as well. It will archive(zip) your extension code,
updating the version and removing any development key in the manifest in the archive.

# Usage
To see full usage run `cws -h` or `cws [command] -h`. To upload and publish an
extension you can run:

```bash
cws deploy --dir=./extension_src
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

# Authentication Setup
To have authentication to the API that we need, we need a already authorized OAuth
client, with offline access so that we can get a single refresh token that will not
expire for a while. The following is a step by step process to do this without the
need to write a script to authenticate.

- Go to your [credentials page for your project](https://console.cloud.google.com/apis/credentials)
- Click *+ CREATE CREDENTIALS* at the top of the page and select *OAuth client ID* and create the app with the following settings.
  - Application Type: Web
  - Name: Your App Name
  - Authorized redirect URIs: https://developers.google.com/oauthplayground
- Save the *Client ID* and *Client Secret*
- Navigate to the [Oauth Playground](https://developers.google.com/oauthplayground)
- Press top right settings (gear) icon (OAuth 2.0 configuration) and fill out the following settings:
  - OAuth flow: Server-side
  - OAuth endpoints: Google
  - Leave the defaults for endpoints
  - Access token location: Authorization header w/ Bearer prefix
  - Access Type: *offline* This is what will allow your refresh token to not expire.
  - Tick *Use your own OAuth credentials*
  - Enter OAuth Client ID and OAuth Client secret
  - Click close to finish.
- At the bottom of the Step 1 *Select & authorize APIs* accordion panel enter required [space-separated scopes](https://developers.google.com/identity/protocols/oauth2/scopes#oauth2):
  - openid (default)
  - https://www.googleapis.com/auth/chromewebstore (access to the chrome webstore API)
- Click Authorize APIs.
  - You may get a warning that the app is not verified, click advanced and then click proceeed.
  - On the next screen choose the account (optional screen) and give the permissions to the app.
- Then you should be redirected back to the playground with *Step 2 Exchange authorization code for tokens* already expanded
- Press Exchange authorization code for tokens. This will get us our refresh token.
- Save the refresh token with your client id and client secret.
