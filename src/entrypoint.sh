#!/usr/bin/env bash

# This project's scripts are meant to be ran on Alpine Linux, which comes with
# busybox instead of the more common GNU coreutils. There are some differences
# between command-line usage to watch out for.
# Reference this documentation: https://linux.die.net/man/1/busybox

cd "${0%/*}"

EXAMS_DIR="./exams"

dialog_cmd="dialog --title 'SSH Exam' --menu 'Available Exams' 0 0 0 "

exam_entries=()
for dir in "$EXAMS_DIR/"*; do
    if [ -d "$dir" ]; then
        exam_entries+=("${dir#$EXAMS_DIR/}")
    fi
done

for n in $(seq 1 ${#exam_entries[@]}); do
    dialog_cmd+="$n \"${exam_entries[($n - 1)]}\" "
done

dialog_cmd+='3>&1 1>&2 2>&3 3>&-'
# Switching STDOUT & STDERR

chosen_exam=$(eval "$dialog_cmd")

clear
for script in "$EXAMS_DIR/${exam_entries[$chosen_exam - 1]}"/*.sh; do
    . "$script"
done

STUDENT="$1"
[ -z "$STUDENT" ] && exit 130

# Scripts under the chosen exam directory should provide these four
# functions:
setup
run
score
cleanup
