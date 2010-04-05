include $(GOROOT)/src/Make.$(GOARCH)

TARG=ghthor/node/atomic[_Type]
GOFILES=\
	atomic[_Type].go\
	type.go\

include $(GOROOT)/src/Make.pkg
