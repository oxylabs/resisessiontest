# Residential Proxies Session Retention Tester
[![](https://dcbadge.vercel.app/api/server/eWsVUJrnG5)](https://discord.gg/GbxmdGhZjq)

Tool to test % of sessions that manage to stay active for ~10 minutes.

Params:

    Required:
        -p string
            Oxylabs proxy password
        -u string
            Oxylabs proxy username

    Optional:
        -cc string
            country parameter (default "fr")
        -city string
            city parameter (default "nice")
        -iptarget string
            ip check target (default "https://ip.oxylabs.io")
        -sessions int
            sessions to test concurrently (default 100)

Example:
```
./resisessiontest -u myUsername -p myPassword -cc us -city los_angeles
```
