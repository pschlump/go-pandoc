
all:
	go vet
	go build


run: all
	./go-pandoc run --config app.conf 

