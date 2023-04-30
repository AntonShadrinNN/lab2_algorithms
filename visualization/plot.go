package visualization

import (
	v1 "lab2_algorithms/tests/contest/v1"
	v2 "lab2_algorithms/tests/contest/v2"
	v3 "lab2_algorithms/tests/contest/v3"
)

func Plot() {
	_, _, execMap := v1.Bench()
	totalBrute := v2.Bench()
	_, _, execTree := v3.Bench()

	//Draw(0, 18, "Total time", "visualization/res/exec_time",
	//	totalMap, totalBrute, totalTree)
	//Draw(0, 18, "Preparation time", "visualization/res/prep_time",
	//	prepMap, nil, prepTree)
	Draw(0, 18, "Execution time", "visualization/res/exec_time_log_new",
		execMap, totalBrute, execTree)
}
