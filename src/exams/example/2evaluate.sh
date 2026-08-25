score() {
    SCORE=0

    if [ -f "/home/$STUDENT/pass" ]; then
        score++
    fi
}
