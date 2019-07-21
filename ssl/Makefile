NAME = mutateme

.PHONY:cert
cert:
	openssl req -x509 -nodes -days 3650 -sha256 -subj /C=aa/ST=aa/L=aa/CN=$(NAME).default.svc.internal -newkey rsa:2048 -keyout $(NAME).key -out $(NAME).pem

.PHONY:clean
clean:
	rm -vf *.key *.pem *.cert *.crt *.csr

