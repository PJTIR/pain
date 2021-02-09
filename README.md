# Networking Project PJTIR

Place all router configs in the `router config` folder in the .cfg format.
Export them by rightclicking and selecting `export config` on a router in GNS3 and saving it in the `router config` folder.

Switches require more work. Perform a `show run` and copy the output into the corresponding .cfg file for that switch.

Make sure to update your changelog with the following whenever you edit a config file for easy backtracking. Do so in the following format:

```
Edit <Config Name> on <Date and time>:
<useful message describing what you did>
```

Simply append this to your changelog file and commit it together with the config change.
