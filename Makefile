# Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGEN=$(GOCMD) generate

# App parameters
GOPI=github.com/djthorpe/gopi
GOLDFLAGS += -X $(GOPI).GitTag=$(shell git describe --tags)
GOLDFLAGS += -X $(GOPI).GitBranch=$(shell git name-rev HEAD --name-only --always)
GOLDFLAGS += -X $(GOPI).GitHash=$(shell git rev-parse HEAD)
GOLDFLAGS += -X $(GOPI).GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOFLAGS = -ldflags "-s -w $(GOLDFLAGS)" 

PKG_CONFIG_PATH_DARWIN=/usr/local/Cellar/libcoap/4.2.1/lib/pkgconfig

darwin: test-darwin

ecovacs:
	$(GOINSTALL) $(GOFLAGS) ./cmd/ecovacs

test-darwin:
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH_DARWIN) $(GOTEST) -v ./libcoap2/...

clean: 
	$(GOCLEAN)
