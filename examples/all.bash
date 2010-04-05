#!/bin/bash

#-- pkg's --#
cd .. && ./pkg.bash && cd examples

#-- Examples --#

# BufferNode Test
echo && echo --- Making Examples ---
echo
echo --- BufferNode --
make -C ./buffer/ clean
make -C ./buffer/
