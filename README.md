EHZ Reader
==========

EHZ SML parser/reader to be used with IR opto heads (tested with ehz ISKRA MT681)


## How to tunnel usb port
1. Host with USB connected

    socat /dev/ttyUSB0,raw,echo=0 tcp-listen:8888,reuseaddr

2. Client

    socat PTY,raw,echo=0,link=<link> tcp:<ip>:8888