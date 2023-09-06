# Residential Proxies Session Retention Tester

Tool to test % of session that manage to stay active for 10 minutes.

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
            ip check target (default "https://ipinfo.io/ip")
        -sessions int
            sessions to test concurrently (default 100)

Example:
```
./sessiontest -u myUsername -p myPassword -cc us -city los_angeles
```