setup() {
    adduser -D -s /usr/bin/nologin "$STUDENT"

    local wd="/home/$STUDENT"
    local home_files="${BASH_SOURCE[0]%/*}/skel"

    #mount -t tmpfs tmpfs "$wd"
    cp -r "$home_files"/* "$home_files"/.* "$wd"

    bins_dir="$wd/.local/bin"
    mkdir -p "$bins_dir"

    local allowed_programs=(
        "/usr/bin/bash"
        "/usr/bin/ls"
        "/usr/bin/whoami"
    )

    for program in "${allowed_programs[@]}"; do
        [ -f "$program" ] && ln -s "$program" "$bins_dir/${program##*/}"
    done

    chown "$STUDENT":"$STUDENT" -R "$wd"
}
