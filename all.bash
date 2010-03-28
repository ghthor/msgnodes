#!/bin/bash

./makeClean
#-- pkg's --#
#make -C ./SortedList/ install
#make -C ./Server/ install

#-- bin's --#
make -C ./Main/
