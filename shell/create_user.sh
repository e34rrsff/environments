TITLE="Touch Test"
DESCRIPTION="The user must `touch` the file \"test\" in their home directory"

USERNAME="$1"
ADDUSER_OPTS="$2"

adduser "$ADDUSER_OPTS" "$USERNAME"
