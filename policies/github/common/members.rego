package common.members

isStale(target_last_active, count_months) {
	now := time.now_ns()
	diff := time.diff(now, target_last_active)
	# diff[0] the year index
	diff[0] >= 1
	} else {
	now := time.now_ns()
    diff := time.diff(now, target_last_active)
	# diff[1] the months index
	diff[1] >= count_months
}
