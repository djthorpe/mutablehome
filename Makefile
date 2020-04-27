# Go parameters
GO=go

# App parameters
GOPI=github.com/djthorpe/gopi/v2/config
GOLDFLAGS += -X $(GOPI).GitTag=$(shell git describe --tags)
GOLDFLAGS += -X $(GOPI).GitBranch=$(shell git name-rev HEAD --name-only --always)
GOLDFLAGS += -X $(GOPI).GitHash=$(shell git rev-parse HEAD)
GOLDFLAGS += -X $(GOPI).GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOFLAGS = -ldflags "-s -w $(GOLDFLAGS)" 

all:
	@echo "Synax: make protogen|ecovacs|tradfri|clean"

protogen:
	@echo Compiling protocol buffers
	@$(GO) generate ./protobuf/...

ecovacs:
	@echo Installing ecovacs to /opt/gaffer
	@install -d /opt/gaffer/bin
	@$(GO) build -o /opt/gaffer/bin/ecovacs $(GOFLAGS) ./cmd/ecovacs

tradfri: protogen
	@echo Installing tradfri to /opt/gaffer
	@install -d /opt/gaffer/bin
	@install -d /opt/gaffer/sbin
	@$(GO) build -o /opt/gaffer/bin/tradfri $(GOFLAGS) ./cmd/tradfri
	@$(GO) build -o /opt/gaffer/sbin/tradfri-service $(GOFLAGS) ./cmd/tradfri-service

clean: 
	$(GOCLEAN)
