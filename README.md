# GoSquash

:warning: **This is not production-ready code!** Although I've deployed it, I'm an irresponsible scoundrel.

This is a tiny app designed to compress images for web requests. It's far from production ready and is mainly for my own learning!

The brief for this app is to sit alongside a Django webapp (to which these pesky users keep uploading massive images) and serve user-uploaded files which are compressed on the fly. Ideally this would be integrated with some sort of Django plugin, but that's beyond the scope of the project for now.

# Deployment Steps

I've deployed this behind Nginx on a Linux server running Ubuntu with the following config.

### Bring on systemd

To keep this Gopher running in the background we'll use systemd with the following config:

```
[Unit]
Description=gosquash

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=/home/user/go/gosquash/main

[Install]
WantedBy=multi-user.target
```

At the following path: `/lib/systemd/system/gosquash.service`

And kick it off using: `sudo service gosquash start`

You can check the status using `sudo service gosquash status`

### Nginx Config

Under `/etc/nginx/sites-available` I modified the config to include:

```
# Gosquash product images
location /product-images {
    proxy_pass http://localhost:9990;
}
```

Then reload Nginx: `sudo nginx -s reload`

---

# Next Up

I want to do the following:

- [x] Write a fancy todo list
- [ ] Move config to command line args
- [ ] Figure out why routing doesn't work with the config above (right now it only works when routing to root `/`)
- [ ] Add a batch job for compressing a set of images all in one go
- [ ] Benchmark the mother
- [ ] Add proper logging to an output file