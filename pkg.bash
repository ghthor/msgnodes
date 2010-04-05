#!/bin/bash
#-- pkg's --#
echo && echo --- Making Packages ---
echo
echo --- node ---
make -C ./node/ install

echo
echo --- buffer ---
make -C ./node/buffer/ install

echo
echo --- debug ---
make -C ./debug/ install

echo
echo --- object ---
make -C ./object/ install
