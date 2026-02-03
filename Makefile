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
	install -C ./init/icstz-server.service ~/.config/systemd/user/$(SERVERNAME).service

uninstall-local :
	rm -f ~/.local/bin/$(SERVERNAME)
	rm -f ~/.config/systemd/user/$(SERVERNAME).service

.PHONY : install-local uninstall-local
