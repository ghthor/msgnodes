#!/bin/bash

./makeClean
#-- pkg's --#
#make -C ./SortedList/ install
#make -C ./Server/ install
make -C ./Init/ install
make -C ./Debug/ install
make -C ./Comm/ install
make -C ./Node/ install
make -C ./Object/ install

#-- bin's --#
make -C ./Main/
