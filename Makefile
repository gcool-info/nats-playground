.PHONY: certificates
certificates: 
	mkcert -install
	# Client certificate
	mkcert -client -cert-file config/client-cert.pem -key-file config/client-key.pem localhost ::1 client@localhost
	# Server certificate
	mkcert -client -cert-file config/server-cert.pem -key-file config/server-key.pem n1 n2 n3 ::1 server@localhost
	# Certificate authority
	cp "$$(mkcert -CAROOT)"/rootCA.pem ./config/rootCA.pem

.PHONY: cleanup
cleanup: 
	rm -rf ./*.pem ./config/*.pem ./persistent-data/*
	mkcert -uninstall
	