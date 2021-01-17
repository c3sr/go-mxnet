#include "cbits/timer.h"
#include <chrono>

int64_t Getoffset(void) {
	using namespace std::chrono;
	high_resolution_clock::time_point t1 = high_resolution_clock::now();
	system_clock::time_point t2 = system_clock::now();
	system_clock::time_point t3 = system_clock::now();
	high_resolution_clock::time_point t4 = high_resolution_clock::now();
	return (static_cast<int64_t>(duration_cast<microseconds>(t2.time_since_epoch()).count()) 
	      - static_cast<int64_t>(duration_cast<microseconds>(t1.time_since_epoch()).count())
	      + static_cast<int64_t>(duration_cast<microseconds>(t3.time_since_epoch()).count())
	      - static_cast<int64_t>(duration_cast<microseconds>(t4.time_since_epoch()).count())) / 2;
}
