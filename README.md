# statusbar

```
  .------------------------------------------------- Network bytes received per second
  |      .------------------------------------------ Network bytes transmitted per second
  |      |       .---------------------------------- Ping time to 1.1.1.1:53
  |      |       |       .-------------------------- Busy CPU per second (100% per thread)
  |      |       |       |    .--------------------- Used RAM
  |      |       |       |    |  .------------------ Free RAM
  |      |       |       |    |  |           .------ Date
  |      |       |       |    |  |     ______|_____
  v      v       v       v    v  v    /            \
2.1 MB  69 kB 1076 ms  125 %  7+10 GB Su 05/07 23:59
\__________________________________________________/
                         |
                 The width is fixed
```

## Usage

```
go get -u github.com/xpetit/statusbar
$(go env GOPATH)/bin/statusbar &
sleep 3
cat /tmp/statusbar
kill %1
```
