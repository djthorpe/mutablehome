
This is the beginning of an interface between mutablehome and google smart home.
Here is what I did to get this far:

* Created a certificate using "Let's Encrypt" for my web server:

```
sudo apt-get install certbot
sudo certbot certonly --standalone
```

This creates two files, and then I am able to use them on my server (which needs to be
accessible though the internet):

```
CERT_PATH=/opt/gaffer/etc/mutablehome.myhome.com
CLIENT_ID=1234
go run ./cmd/mutablehome-googleactions/... \
  -port 9001
  -client_id ${CLIENT_ID}
  -sslkey ${CERTPATH}/privkey1.pem -sslcert ${CERTPATH}/cert1.pem \
  -debug
```

Here, the Client ID is used when creating your project in the Google Console
[here](https://console.actions.google.com/u/0/). You need to create a project,
choose "Smart Home" and then run through the "Quick Setup", etc. If for example
your web server is accessible at https://mutablehome.myhome.com:9001/ as per the
example above, you need to use the following values:

  * Setup Account Linking, linking type OAuth implicit
  * Client Information, Client ID 1234 (as per above)
  * Client Information, Authorization URL https://mutablehome.myhome.com:9001/oauth2
  * Create Smart Home Actions, Fulfillment URL https://mutablehome.myhome.com:9001/mutablehome

For whatever reason, you can't use the simulator here, you need to use your Google Home app,
and there you'll need to "Add" on the home screen and then "Set up device". On the next page,
choose "Works with Google" and then choose your test action, which should have the name
you chose with "[test]" written before it.

The "Sync" intent is done but clearly needs some work. The query, execute etc. are not started,
so you won't actually be able to use this code yet and it needs heavily re-structured too.

I would definately check out the google documentation [here](https://developers.google.com/actions/smarthome/develop/create).
