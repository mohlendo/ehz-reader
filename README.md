EHZ Reader
==========

EHZ SML parser/reader to be used with IR opto heads (tested with ehz ISKRA MT681)


## How to tunnel USB/Serial Port through network
1. On the server with USB connected

    socat /dev/ttyUSB0,raw,echo=0 tcp-listen:8888,reuseaddr

2. On the client

    socat PTY,raw,echo=0 tcp:<ip>:8888