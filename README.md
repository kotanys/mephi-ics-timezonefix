## Proxy server for fixing timezone information in calendar files

The RFC 5545 aka the standard that defines the iCalendar (.ics) format describes 3 ways
to declare the time of the event:
1. LOCAL TIME: time with no timezone info attached, said to be _"floating"_
2. UTC TIME: time in the UTC+0 timezone
3. LOCAL TIME AND TIME ZONE REFERENCE: time in the specified timezone

The service I happened to export calendar files from used LOCAL TIME despite
the timezone being declared in the file. Floating time is unreliable
for obvious reasons (I could be in a different timezone than the one of the event),
so I coded this little server that adds timezone information on the fly by converting 
type 1 time to type 3 time.

I have contacted the people behind the service in question and they have fixed the calendar export 
functionality on their end, but this proxy server proved to be a neat temporary solution.

This thing uses no 3rd-party libraries, parsing the ICS with regexes. So fingers crossed
that the calendar file is valid enough I guess :)

### Build and run
```shell
$ make build install-local
$ ~/.local/bin/icstz-server 24501
```
Here, 24501 is a TCP port that the server listens to. Can be any other port though.

For systemd systems, `make install-local` installs a user service named `icstz-server.service`.
To run the server as a systemd service, execute
```shell
$ systemctl --user enable --now icstz-server.service
```

Alternatively, run as a docker container:
```shell
$ docker compose up -d
```

Then instead of downloading the calendar from the remote server, download it from the proxy
by connecting to `http://localhost:24501/fix?id=ARG`, where `ARG` is the `%v` of the `s.FmtString` defined
in [main.go](./cmd/icstz-server/main.go), in the `main` function.
This is very much tied to the particular website I was using, shouldn't be difficult to "untie"
if I ever need this thing again.

> [!NOTE]
> The docker container and the systemd service
> listen to TCP port 24501 by default.

> [!NOTE]
> The URL that the server downloads the calendar from is 
> hardcoded in the [main.go](./cmd/icstz-server/main.go), sorry for that.
