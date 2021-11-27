# verbose-twit-banner

Replace your Twitter profile banner with an image
that contains your profile statistics

![sample.png](sample.png)

# Usage

## Run from Source

```bash
make run
```

## Build a Binary

```bash
make build
```

Executable will be created as verbose-twit-banner

## Ways of Operation

### Run in an Interval

Parameter `-interval=<string>` can be used
to make the program run in a loop.
The banner image will be updated
in the specified time interval (e.g. 5m).
Valid time units are "ns", "us" (or "µs"),
"ms", "s", "m", "h".

### Read-only mode

Parameter `-dry-run` can be used
to make the program output an image.

Could be useful if you want to see
what image it would output
without replacing your banner image

Another advantage of this method is that you don't need
Twitter Access Key & Secret
so you can run it on anyone's Twitter account.

## Twitter API Access

You'll need to create a Project on Twitter to get access to the API

Details: https://developer.twitter.com/en/docs/twitter-api/getting-started/getting-access-to-the-twitter-api

### Obtaining App Consumer Key & Secret

Consumer Key & Secret are for accessing public Twitter profile statistics.
These are used in OAuth 2.0 authentication with Twitter.

API calls made using this set of credentials:

- [GET /2/users/by/username/:username](https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-by-username-username)
- [GET /2/users/:id](https://developer.twitter.com/en/docs/twitter-api/users/lookup/api-reference/get-users-id)

### Obtaining Access Token & Secret

Access Token & Secret are for uploading the generated banner to your Twitter profile.
These are used in OAuth 1.1 authentication with Twitter.

API call made using this set of credentials:

- [POST account/update_profile_banner](https://developer.twitter.com/en/docs/twitter-api/v1/accounts-and-users/manage-account-settings/api-reference/post-account-update_profile_banner)

## Configuration

The binary has a lot of parameters.
Some of them also have equivalent environment variables.

### Command Line Parameters

Use `-h` to find out about the available parameters:
```bash
./verbose-twit-banner -h
```

### Environment Variables

File named [.env.dist](.env.dist) contains
all the available environment variables.

Instead of setting these environment variables
you can also copy the `.env.dist` file into `.env`.

# Tests

```bash
make test
```

# Why?

Why not?
