PATH=$PATH:/home/robking/AGolang/go1.18/go/bin
rm -f mr-*
rm -f wc.so
go build -buildmode=plugin ../mrapps/wc.go