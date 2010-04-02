#!/bin/bash

./makeClean
#-- pkg's --#
#make -C ./SortedList/ install
#make -C ./Server/ install
make -C ./init/ install
make -C ./debug/ install
make -C ./comm/ install
make -C ./node/ install
make -C ./object/ install

#-- bin's --#
make -C ./main/
