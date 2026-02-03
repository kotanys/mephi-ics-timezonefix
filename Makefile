SERVERDIR:=cmd/icstz-server
SERVERNAME:=icstz-server
SERVER:=$(SERVERDIR)/$(SERVERNAME)

build : icstz-server

$(SERVERNAME) : $(SERVER)
.PHONY : $(SERVERNAME)

$(SERVER) : $(SERVERDIR)/*.go
	go build -C $(SERVERDIR)

install-local : $(SERVERNAME)
	install -C $(SERVER) ~/.local/bin/$(SERVERNAME)

uninstall-local :
	rm -f ~/.local/bin/$(SERVERNAME)
