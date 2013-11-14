# Contrail

A small wrapper around glog to make category and trace logging dead simple.

See the godoc documentation for details.

CHANGELOG:

* 2013-10-17: Initial implementation of contexts and trace IDs.
* 2013-10-21: Added LogWriter to adapt any non-formattd log function to io.Writer.
* 2013-10-22: Added SetStderrThreshold to Logger interface.
* 2013-11-14: Return exported interface, not internal struct.
