package common.secrets

is_stale(date) {
    diff := time.diff(time.now_ns(),date)
    diff[0] >= 1
}
