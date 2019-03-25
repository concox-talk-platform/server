/**
 * Copyrights (c) 2019. All rights reserved.
 * Group data type module
 * Author: tesion
 * Date: March 25th 2019
 */
package group

type RoleType uint8

const (
    GROUP_NORMAL_USER RoleType = iota
    GROUP_MANAGER
)
