# statusbar

A statusbar for Linux

```
  .------------------------------------------ Network bytes exchanged per second
  |       .---------------------------------- Highest temperature reading
  |       |      .--------------------------- Busy CPU per second (100% per thread)
  |       |      |       .------------------- Free RAM
  |       |      |       |           .------- Date
  |       |      |       |      _____|______
  v       v      v       v     /            \
2.1 MB │ 62° │ 125 % │ 10 GB │ Su 05/07 23:59
\___________________________________________/
                     |
             The width is fixed
```

### Installation

```
go install github.com/xpetit/statusbar@latest
```

`lm-sensors` is required. Add an entry to `crontab -e`:

```
@reboot sleep 2 && /home/USERNAME/go/bin/statusbar &
```

### Usage

```
nc -U /var/run/user/1000/statusbar.sock
```

With `xfce4-genmon-plugin` I somehow need to run this command instead:

```
sh -c 'nc -U /var/run/user/1000/statusbar.sock </dev/null'
```
