#!/usr/bin/env bash

set -o errexit

function install_go() {
	local go_version="1.20.2"
	local download="https://storage.googleapis.com/golang/go${go_version}.linux-amd64.tar.gz"

		if go version 2>&1 | grep -q "${go_version}"; then
		return
	fi

	# remove previous older version
	rm -rf /usr/local/go

	# retry downloading on spurious failure
	curl -sSL --fail -o /tmp/go.tar.gz \
		--retry 5 --retry-connrefused \
		"${download}"

	tar -C /tmp -xf /tmp/go.tar.gz
	sudo mv /tmp/go /usr/local
	sudo chown -R root:root /usr/local/go
}

install_go

# Ensure that the GOPATH tree is owned by vagrant:vagrant
mkdir -p /opt/gopath
chown -R vagrant:vagrant /opt/gopath

# Ensure Go is on PATH
if [ ! -e /usr/bin/go ] ; then
	ln -s /usr/local/go/bin/go /usr/bin/go
fi
if [ ! -e /usr/bin/gofmt ] ; then
	ln -s /usr/local/go/bin/gofmt /usr/bin/gofmt
fi


# Ensure new sessions know about GOPATH
if [ ! -f /etc/profile.d/gopath.sh ] ; then
	cat <<EOT > /etc/profile.d/gopath.sh
export GOPATH="/opt/gopath"
export PATH="/opt/gopath/bin:\$PATH"
EOT
	chmod 755 /etc/profile.d/gopath.sh
fi
