/*
	Digivance MVC Application Framework
	Session Manager Tests
	Dan Mayor (dmayor@digivance.com)

	This file provides some unit tests for the in memory session manager system.

	This package is released under as open source under the LGPL-3.0 which can be found:
	https://opensource.org/licenses/LGPL-3.0
*/

package mvcapp

import (
	"fmt"
	"testing"
	"time"
)

// TestSessionManager tests various aspects of the Session Manager. Also demonstrates
// some of the basic functionality.
func TestSessionManager(t *testing.T) {
	// Needs rewritten
}

func TestTimeoutCalculations(t *testing.T) {
	now := time.Now()
	recentActivity := now.Add(-10 * time.Minute)
	expiredActivity := now.Add(-20 * time.Minute)

	if recentActivity.Add(15 * time.Minute).Before(now) {
		t.Error("Reading in range activity as expired")
	} else {
		fmt.Println("Recent activity valid")
	}

	if expiredActivity.Add(15 * time.Minute).Before(now) {
		fmt.Println("Expired activity rejected correctly")
	} else {
		t.Error("Reading expired activity as valid")
	}
}
