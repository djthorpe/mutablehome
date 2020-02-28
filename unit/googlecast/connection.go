/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package googlecast

import "sync"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type connection struct {
	sync.Mutex
}
