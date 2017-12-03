// Copyright 2017 Joshua J Baker. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// +build darwin netbsd freebsd openbsd dragonfly

package base

import (
    "syscall"
    "time"
)

const (
    evListen = syscall.EVFILT_READ
    evRead   = syscall.EVFILT_READ
    evWrite  = syscall.EVFILT_READ
    evClient = syscall.EVFILT_READ | syscall.EVFILT_WRITE
)

type PollBase struct {
    events []syscall.Kevent_t
}

func addPollEvent(p, fd int, masks uint32) (err error) {
    if masks&EventRead == EventRead {
        _, err = syscall.Kevent(p, []syscall.Kevent_t{{
            Ident:  uint64(fd),
            Flags:  syscall.EV_ADD,
            Filter: syscall.EVFILT_READ}}, nil, nil)
    }
    if masks&EventWrite == EventWrite {
        _, err = syscall.Kevent(p, []syscall.Kevent_t{{
            Ident:  uint64(fd),
            Flags:  syscall.EV_ADD,
            Filter: syscall.EVFILT_WRITE}}, nil, nil)
    }

    return
}

func deletePollEvent(p, fd int, masks uint32) (err error) {
    if masks&EventRead == EventRead {
        _, err = syscall.Kevent(p, []syscall.Kevent_t{{
            Ident:  uint64(fd),
            Flags:  syscall.EV_DELETE,
            Filter: syscall.EVFILT_READ}}, nil, nil)
    }
    if masks&EventWrite == EventWrite {
        _, err = syscall.Kevent(p, []syscall.Kevent_t{{
            Ident:  uint64(fd),
            Flags:  syscall.EV_DELETE,
            Filter: syscall.EVFILT_WRITE}}, nil, nil)
    }

    return
}

func AddRead(p, fd int, readon, writeon *bool) error {
    if readon != nil {
        if *readon {
            return nil
        }
        *readon = true
    }
    _, err := syscall.Kevent(p,
        []syscall.Kevent_t{{Ident: uint64(fd),
            Flags: syscall.EV_ADD, Filter: syscall.EVFILT_READ}},
        nil, nil)
    return err
}
func DelRead(p, fd int, readon, writeon *bool) error {
    if readon != nil {
        if !*readon {
            return nil
        }
        *readon = false
    }
    _, err := syscall.Kevent(p,
        []syscall.Kevent_t{{Ident: uint64(fd),
            Flags: syscall.EV_DELETE, Filter: syscall.EVFILT_READ}},
        nil, nil)
    return err
}

func AddWrite(p, fd int, readon, writeon *bool) error {
    if writeon != nil {
        if *writeon {
            return nil
        }
        *writeon = true
    }
    _, err := syscall.Kevent(p,
        []syscall.Kevent_t{{Ident: uint64(fd),
            Flags: syscall.EV_ADD, Filter: syscall.EVFILT_WRITE}},
        nil, nil)
    return err
}

func DelWrite(p, fd int, readon, writeon *bool) error {
    if writeon != nil {
        if !*writeon {
            return nil
        }
        *writeon = false
    }
    _, err := syscall.Kevent(p,
        []syscall.Kevent_t{{Ident: uint64(fd),
            Flags: syscall.EV_DELETE, Filter: syscall.EVFILT_WRITE}},
        nil, nil)
    return err
}

func MakePoll() (p int, err error) {
    return syscall.Kqueue()
}

func CreatePoll() (p int, err error) {
    return syscall.Kqueue()
}

func MakeEvents(n int) interface{} {
    return make([]syscall.Kevent_t, n)
}

func Wait(p int, evs interface{}, timeout time.Duration) (n int, err error) {
    if timeout < 0 {
        timeout = 0
    }
    ts := syscall.NsecToTimespec(int64(timeout))
    return syscall.Kevent(p, nil, evs.([]syscall.Kevent_t), &ts)
}

func GetFD(evs interface{}, i int) int {
    return int(evs.([]syscall.Kevent_t)[i].Ident)
}

func closePoll(p int) (err error) {
    return syscall.Close(p)
}

func (r *PollReactor) wait(timeout int64) (n int, err error) {
    if timeout < 0 {
        return syscall.Kevent(r.pollerFD, nil, r.events, nil)
    }

    ts := syscall.NsecToTimespec(int64(timeout))
    return syscall.Kevent(r.pollerFD, nil, r.events, &ts)
}

func (r *PollReactor) createEvents(n int) (err error) {
    r.events = make([]syscall.Kevent_t, n)
    return
}
