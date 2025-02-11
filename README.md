# Residential Proxies Session Retention Tester

[![Oxylabs promo code](https://raw.githubusercontent.com/oxylabs/product-integrations/refs/heads/master/Affiliate-Universal-1090x275.png)](https://oxylabs.go2cloud.org/aff_c?offer_id=7&aff_id=877&url_id=112)

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
