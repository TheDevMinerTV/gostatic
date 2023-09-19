# GoStatic

Probably the fastest (and simplest) way to serve static files.

```bash
docker run --rm -p 8080:80 -v $PWD:/static:ro ghcr.io/thedevminertv/gostatic
```

Additionally, you can add the following flags:

```bash
-cache <duration> # Enable caching
-compress-level <level> # Enable compression (-1=remove module, 0=none, 2=best)
-log-requests # Log requests (this will slow down the server by ~40%)
```

You can just add these to the end of the `docker run` command.
